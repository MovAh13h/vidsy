package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"

	awsclient "github.com/MovAh13h/vidsy/internal/aws"
	"github.com/MovAh13h/vidsy/internal/common"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Consumer struct {
	Awsclient *awsclient.AwsClient
	QueueUrl *string
}

func NewConsumer(name, profile *string) (*Consumer, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile(*profile))

	if err != nil {
		return nil, err
	}

	awsclient, err := awsclient.NewAwsClient(&cfg)

	if err != nil {
		return nil, err
	}

	o, err := awsclient.Sts.GetCallerIdentity(context.Background())

	if err != nil {
		return nil, err
	}

	q, err := awsclient.Sqs.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
		QueueName: name,
		QueueOwnerAWSAccountId: o.Account,
	})

	if err != nil {
		return nil, err
	}

	return &Consumer{Awsclient: awsclient, QueueUrl: q.QueueUrl}, nil
}

func (c *Consumer) Start() {
	log.Println("Starting Consumer...")

	var wg sync.WaitGroup

	for {
		m, err := c.Awsclient.Sqs.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
			QueueUrl: c.QueueUrl,
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     20,
		})

		if err != nil {
			log.Printf("%v", err)
		}

		if len(m.Messages) > 0 {
			var qj common.QueueJob

			err = json.Unmarshal([]byte(*m.Messages[0].Body), &qj)

			if err != nil {
				log.Printf("%v", err)
			}

			wg.Add(1)
			go c.process(&qj, &wg)	
		}
	}
}

func (c *Consumer) process(qj *common.QueueJob, wg *sync.WaitGroup) {
	// dlock, err := c.Awsclient.Dynamo.AcquireLock(context.Background(), &qj.Src)

	// if err != nil {
	// 	log.Printf("%v", err)
	// 	return
	// }

	baseSrcPath := path.Base(qj.Src)
	bucketName := "vidsy-store"
	fileName := baseSrcPath[:len(baseSrcPath)-len(path.Ext(baseSrcPath))]

	err := os.MkdirAll(fmt.Sprintf("./%s/%s/out/%s", bucketName, fileName, qj.Resolution.String()), 0755)

	if err != nil {
		log.Printf("%v", err)
		return
	}

	// _, err := os.Stat(fmt.Sprintf("./%s/%s/%s.mp4", bucketName, fileName, fileName))

	// if os.IsNotExist(err) {
	// 	o, err := c.Awsclient.S3.GetObject(context.Background(), &s3.GetObjectInput{
	// 		Bucket: &bucketName,
	// 		Key: &baseSrcPath,
	// 	})

	// 	if err != nil {
	// 		log.Printf("%v", err)
	// 		return
	// 	}

	// 	err = 

	// 	if err != nil {
	// 		log.Printf("%v", err)
	// 		return
	// 	}

	// 	file_desc, err := os.Create(fmt.Sprintf("./%s/%s/%s.mp4", bucketName, fileName, fileName))

	// 	if err != nil {
	// 		log.Printf("%v", err)
	// 		return
	// 	}

	// 	reader := bufio.NewReader(o.Body)
	// 	writer := bufio.NewWriter(file_desc)

	// 	_, err = io.Copy(writer, reader)

	// 	if err != nil {
	// 		log.Printf("%v", err)
	// 		return
	// 	}

	// 	writer.Flush()
	// 	file_desc.Close()
	// } else if err != nil {
	// 	log.Printf("%v", err)
	// 	return
	// }

	// c.Awsclient.Dynamo.ReleaseLock(context.Background(), dlock)
	
	if qj.OutputFormat.String() == "MPEG_DASH" {
		c.dashProcessor(qj, bucketName, fileName)
	} else if qj.OutputFormat.String() == "HLS" {
		c.hlsProcessor(qj, bucketName, fileName)
	}

	wg.Done()
}

func (c *Consumer) hlsProcessor(qj *common.QueueJob, bucket, filename string) {
	fpath := fmt.Sprintf("./%s/%s/%s.mp4", "vidsy-store", filename, filename)
	opath := fmt.Sprintf("./%s/%s/out/%s/playlist.m3u8", bucket, filename, qj.Resolution.String())

	fmt.Println("ffmpeg", "-i", fpath, "-vf", getScale(qj), "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", opath)

	cmd := exec.Command("ffmpeg", "-i", fpath, "-vf", getScale(qj), "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", opath)

	_, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("%v", err)
		return
	}

	err = filepath.Walk(fmt.Sprintf("./%s/%s/out/%s", bucket, filename, qj.Resolution.String()), func (filePath string, info os.FileInfo, err error) error {
		if err != nil {
            return err
        }

        if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
            return nil
        }

        fileBytes, err := ioutil.ReadFile(filePath)

        if err != nil {
            return err
        }

		key := fmt.Sprintf("%s/%s/%s", filename, qj.Resolution.String(), path.Base(filePath))

		_, err = c.Awsclient.S3.PutObject(context.Background(), &s3.PutObjectInput{
            Bucket: &bucket,
            Key:    &key,
            Body:   bytes.NewReader(fileBytes),
            ACL:    types.ObjectCannedACLPrivate, // Set ACL to private or public-read as needed
        })

        if err != nil {
            return err
        }

        return nil
	})

	if err != nil {
		log.Printf("%v", err)
		return
	}
}

func (c *Consumer) dashProcessor(qj *common.QueueJob, bucket, key string) {
	fpath := fmt.Sprintf("./%s/%s/%s.mp4", "vidsy-store", "key", "key")
	opath := fmt.Sprintf("./%s/%s/out/%s/playlist.m3u8", bucket, key, qj.Resolution.String())

	cmd := exec.Command("ffmpeg", "-i", fpath, "-codec:", "copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", opath)

	out, err := cmd.Output()

	if err != nil {
		fmt.Println(out)
		log.Printf("%v", err)
		return
	}

	file, err := os.Open(opath)

	if err != nil {
		log.Printf("%v", err)
		return
	}

	defer file.Close()	

	_, err = c.Awsclient.S3.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key: &key,
		Body: bufio.NewReader(file),
	})

	if err != nil {
		log.Printf("%v", err)
	}
}

func getScale(qj *common.QueueJob) string {
	switch qj.Resolution.String() {
		case "P144":
			return "scale=256:144"
		case "P240":
			return "scale=426:240"
		case "P360":
			return "scale=640:360"
		case "P480":
			return "scale=854:480"
		case "P720":
			return "scale=1280:720"
		case "P1080":
			return "scale=1920:1080"
		case "P1440":
			return "scale=2560:1440"
		case "P2160":
			return "scale=3840:2160"
		default:
			return "scale=256:144"
	}
} 







