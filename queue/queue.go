package queue

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"time"
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
	message := &sns.PublishInput{
		Message:  aws.String(string(json)),
		TopicArn: aws.String(snsTopicArn),
	}
	_, err := snsService.Publish(message)
	return err
}
