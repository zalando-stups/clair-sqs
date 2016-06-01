package queue

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"

	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/base64"
	"log"
	"time"
)

const (
	SNS_CONTENT_TYPE_KEY              = "CLAIR.CONTENTTYPE"
	SNS_CONTENT_TYPE_VALUE_JSON       = "application/json"
	SNS_CONTENT_TYPE_VALUE_BASE64GZIP = "application/base64gzip"
	MAX_MESSAGE_SIZE                  = 250 * 1024 // some metadata overhead to 256k
)

func ProcessMessages(sqsService *sqs.SQS, sqsQueueUrl string, messageProcessor func(msgid, msg string) error) {
	for {
		// wait and receive messages
		req := &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(sqsQueueUrl),
			WaitTimeSeconds:     aws.Int64(20),
			MaxNumberOfMessages: aws.Int64(1),
		}
		resp, err := sqsService.ReceiveMessage(req)

		if err != nil {
			log.Printf("Cannot fetch messages from SQS queue %v: %v", sqsQueueUrl, err)
			time.Sleep(10 * time.Second)
			continue
		}

		for _, msg := range resp.Messages {
			err = messageProcessor(*msg.MessageId, *msg.Body)

			if err != nil {
				log.Printf("Couldn't process message %v: %v (message left in queue)", msg.MessageId, err)
				continue
			}

			// delete message from queue
			_, err = sqsService.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(sqsQueueUrl),
				ReceiptHandle: msg.ReceiptHandle,
			})

			if err != nil {
				log.Printf("Could't delete message %v from queue: %v", msg.MessageId, err)
				continue
			}
		}
	}
}

func SendNotification(snsService *sns.SNS, snsTopicArn string, json []byte) error {
	attributes := make(map[string]*sns.MessageAttributeValue)

	var message []byte
	if len(json) <= MAX_MESSAGE_SIZE {
		// send raw json
		attributes[SNS_CONTENT_TYPE_KEY] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(SNS_CONTENT_TYPE_VALUE_JSON),
		}

		message = json
	} else {
		// we need to gzip the content
		attributes[SNS_CONTENT_TYPE_KEY] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(SNS_CONTENT_TYPE_VALUE_BASE64GZIP),
		}

		var b bytes.Buffer
		b64 := base64.NewEncoder(base64.StdEncoding, &b)
		gz, err := gzip.NewWriterLevel(b64, flate.BestCompression)
		if err != nil {
			return err
		}

		if _, err = gz.Write(json); err != nil {
			return err
		}

		if err = gz.Flush(); err != nil {
			return err
		}
		if err = b64.Flush(); err != nil {
			return err
		}
		if err = gz.Close(); err != nil {
			return err
		}
		if err = b64.Close(); err != nil {
			return err
		}
		message = b.Bytes()
	}

	if len(message) > MAX_MESSAGE_SIZE {
		log.Printf("SNS message probably too big: %v > %v trying anyway (%v: %v).", len(message), MAX_MESSAGE_SIZE,
			SNS_CONTENT_TYPE_KEY, attributes[SNS_CONTENT_TYPE_KEY].StringValue)
	}

	snsmsg := &sns.PublishInput{
		TopicArn:          aws.String(snsTopicArn),
		MessageAttributes: attributes,
		Message:           aws.String(string(message)),
	}
	_, err := snsService.Publish(snsmsg)
	return err
}
