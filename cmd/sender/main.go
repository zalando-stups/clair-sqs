// The sending notifications of Clair to SNS
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/zalando/clair-sqs/clair"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/zalando/clair-sqs/queue"
)

type notificationEnvelope struct {
	Notification struct {
		Name string
	}
}

func main() {
	snsTopicArn := os.Getenv("SENDER_TOPIC_ARN")
	snsTopicRegion := os.Getenv("SENDER_TOPIC_REGION")
	clairUrl := os.Getenv("CLAIR_URL")

	if snsTopicArn == "" || snsTopicRegion == "" || clairUrl == "" {
		log.Fatal("SENDER_TOPIC_ARN, SENDER_TOPIC_REGION or CLAIR_URL not set")
	}

	log.Printf("Sender Topic ARN: %v", snsTopicArn)
	log.Printf("Sender Topic Region: %v", snsTopicRegion)

	svc := sns.New(session.New(&aws.Config{Region: &snsTopicRegion}))

	http.HandleFunc("/trigger", func(w http.ResponseWriter, r *http.Request) {
		// Clair triggered us, read the notification name
		decoder := json.NewDecoder(r.Body)
		var notification notificationEnvelope
		err := decoder.Decode(&notification)
		if err != nil {
			log.Println("Got invalid notification JSON.")
			return
		}

		err = clair.ProcessNotification(clairUrl, notification.Notification.Name, func(newLayers []string, oldLayers []string) error {
			layers := append(newLayers, oldLayers...)

			for _, layer := range layers {
				// send the raw notification details to the SNS notifcation
				details, err := clair.GetLayer(clairUrl, layer)
				if err != nil {
					return err
				}

				if err := queue.SendNotification(svc, snsTopicArn, details); err != nil {
					return err
				}

				log.Printf("New notification %v page forwarded to %v", notification.Notification.Name, snsTopicArn)
			}

			return nil
		})

		if err != nil {
			log.Printf("Couldn't process notification %v: %v", notification.Notification.Name, err)
		}

		// mark notification as processed
		if err = clair.DeleteNotification(clairUrl, notification.Notification.Name); err != nil {
			log.Printf("Couldn't delete notification %v after it was processed: %v", notification.Notification.Name, err)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":7070", nil))
}
