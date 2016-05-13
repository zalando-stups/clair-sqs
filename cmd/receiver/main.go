// The receiver is receiving SQS messages and forwards them to Clair
package main

import (
	"log"
	"os"

	"github.com/zalando/clair-sqs/clair"

	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/zalando/clair-sqs/queue"
)

type pushMessage struct {
	Layer struct {
		Name          string
		Path          string
		Authorization string
		ParentName    string
		Format        string
	}
}

func main() {
	sqsQueueUrl := os.Getenv("RECEIVER_QUEUE_URL")
	sqsQueueRegion := os.Getenv("RECEIVER_QUEUE_REGION")
	snsTopicArn := os.Getenv("RECEIVER_TOPIC_ARN")
	snsTopicRegion := os.Getenv("RECEIVER_TOPIC_REGION")
	clairUrl := os.Getenv("CLAIR_URL")

	if sqsQueueUrl == "" || sqsQueueRegion == "" {
		log.Fatal("RECEIVER_QUEUE_URL or RECEIVER_QUEUE_REGION not set")
	}

	if snsTopicArn == "" || snsTopicRegion == "" {
		log.Fatal("RECEIVER_TOPIC_ARN or RECEIVER_TOPIC_REGION not set")
	}

	if clairUrl == "" {
		log.Fatal("CLAIR_URL not set")
	}

	log.Printf("Receiver Queue URL: %v", sqsQueueUrl)
	log.Printf("Receiver Queue Region: %v", sqsQueueRegion)

	sqsService := sqs.New(session.New(&aws.Config{Region: &sqsQueueRegion}))

	log.Printf("Receiver Topic ARN: %v", snsTopicArn)
	log.Printf("Receiver Topic Region: %v", snsTopicRegion)

	snsService := sns.New(session.New(&aws.Config{Region: &snsTopicRegion}))

	queue.ProcessMessages(sqsService, sqsQueueUrl, func(msgid, msg string) error {
		var jsonMessage pushMessage
		if err := json.Unmarshal([]byte(msg), &jsonMessage); err != nil {
			return err
		}

		// forward to Clair
		if err := clair.PushLayer(clairUrl, []byte(msg)); err != nil {
			return err
		}

		// send details to SNS
		details, err := clair.GetLayer(clairUrl, jsonMessage.Layer.Name)
		if err != nil {
			return err
		}
		if err := queue.SendNotification(snsService, snsTopicArn, details); err != nil {
			return err
		}

		log.Printf("Forwarded message %v to Clair", msgid)

		return nil
	})
}
