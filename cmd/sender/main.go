// The sending notifications of Clair to SNS
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type notificationEnvelope struct {
	Notification struct {
		Name string
	}
}

type notificationDetailEnvelope struct {
	Notification struct {
		NextPage string
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Clair triggered us, read the notification name
		decoder := json.NewDecoder(r.Body)
		var notification notificationEnvelope
		err := decoder.Decode(&notification)
		if err != nil {
			log.Println("Got invalid notification JSON.")
			return
		}

		notificationUrl := fmt.Sprintf("%v/v1/notifications/%v", clairUrl, notification.Notification.Name)

		// fetch detailed information from clair
		var page = ""
		for {
			var pageUrl string
			if page == "" {
				pageUrl = fmt.Sprintf("%v?limit=%v", notificationUrl, 10)
			} else {
				pageUrl = fmt.Sprintf("%v?limit=%v&page=%v", notificationUrl, 10, page)
			}

			details, err := http.Get(pageUrl)
			if err != nil {
				log.Printf("Couldn't get details about notification %v", notification.Notification.Name)
				return
			}

			detailsBytes, err := ioutil.ReadAll(details.Body)
			if err != nil {
				log.Printf("Couldn't read response from network socket for notification %v", notification.Notification.Name)
			}
			detailsString := string(detailsBytes)

			var detailsJson notificationDetailEnvelope
			err = json.Unmarshal(detailsBytes, &detailsJson)
			if err != nil {
				log.Println("Got invalid notification detail JSON.")
				return
			}

			// send the raw notification details to the SNS notifcation
			message := &sns.PublishInput{
				Message:  aws.String(detailsString),
				TopicArn: aws.String(snsTopicArn),
			}
			_, err = svc.Publish(message)

			if err != nil {
				log.Printf("Couldn't send notification to SNS topic %v: %v", snsTopicArn, err.Error())
				return
			}

			log.Printf("New notification %v page %v forwarded to %v", notification.Notification.Name, page, snsTopicArn)

			// another page to process?
			if detailsJson.Notification.NextPage != "" {
				page = detailsJson.Notification.NextPage
			} else {
				break
			}
		}

		// mark notification as processed
		_, err = http.NewRequest("DELETE", notificationUrl, nil)
		if err != nil {
			log.Printf("Couldn't delete notification %v after it was processed: %v", notification.Notification.Name, err)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":7070", nil))
}
