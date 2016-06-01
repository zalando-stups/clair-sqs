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
	"github.com/zalando/go-tokens/tokens"
)

type pushMessage struct {
	Layer struct {
		Name       string
		Path       string
		Headers    map[string]string
		ParentName string
		Format     string
	}
}

func main() {
	sqsQueueUrl := os.Getenv("RECEIVER_QUEUE_URL")
	sqsQueueRegion := os.Getenv("RECEIVER_QUEUE_REGION")
	snsTopicArn := os.Getenv("RECEIVER_TOPIC_ARN")
	snsTopicRegion := os.Getenv("RECEIVER_TOPIC_REGION")
	clairUrl := os.Getenv("CLAIR_URL")
	accessTokenUrl := os.Getenv("ACCESS_TOKEN_URL")

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

	log.Printf("Receiver Access Token URL: %v", accessTokenUrl)
	reqs := []tokens.ManagementRequest{
		tokens.NewPasswordRequest("fetch-layer", "uid"),
	}
	tokenManager, err := tokens.Manage(accessTokenUrl, reqs)
	if err != nil && err != tokens.ErrMissingURL {
		// if tokens.ErrMissingUrl, then tokenManage is nil and can be tested (feature off)
		panic(err)
	}

	queue.ProcessMessages(sqsService, sqsQueueUrl, func(msgid, msg string) error {
		var jsonMessages []pushMessage

		if msg[0:1] == "[" {
			// batch message
			if err := json.Unmarshal([]byte(msg), &jsonMessages); err != nil {
				return err
			}
		} else {
			// single message
			var jsonMessage pushMessage
			if err := json.Unmarshal([]byte(msg), &jsonMessage); err != nil {
				return err
			}
			jsonMessages = append(jsonMessages, jsonMessage)
		}

		for _, jsonMessage := range jsonMessages {
			if jsonMessage.Layer.Headers == nil {
				jsonMessage.Layer.Headers = make(map[string]string)
			}

			// optional: enrich mesage with authorization token
			var jsonMessageBytes []byte
			if tokenManager != nil {
				token, err := tokenManager.Get("fetch-layer")
				if err != nil {
					return err
				}
				jsonMessage.Layer.Headers["Authorization"] = "Bearer " + token.Token
			}

			jsonMessageBytes, err = json.Marshal(jsonMessage)
			if err != nil {
				return err
			}

			// forward to Clair
			if err := clair.PushLayer(clairUrl, jsonMessageBytes); err != nil {
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

			log.Printf("Indexed layer %v (%v) from message %v and sent result notification.", jsonMessage.Layer.Name, jsonMessage.Layer.Path, msgid)
		}

		return nil
	})
}
