// The receiver is receiving SQS messages and forwards them to Clair
package main

import (
    "fmt"
    "time"
    "os"
    "log"
    "strings"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sqs"
    "net/http"
)

func main() {
    sqsQueueUrl := os.Getenv("RECEIVER_QUEUE_URL")
    sqsQueueRegion := os.Getenv("RECEIVER_QUEUE_REGION")
    clairUrl := os.Getenv("CLAIR_URL")

    if sqsQueueUrl == "" || sqsQueueRegion == "" || clairUrl == "" {
        log.Fatal("RECEIVER_QUEUE_URL, RECEIVER_QUEUE_REGION or CLAIR_URL not set")
    }

    log.Printf("Receiver Queue URL: %v", sqsQueueUrl)
    log.Printf("Receiver Queue Region: %v", sqsQueueRegion)

    svc := sqs.New(session.New(&aws.Config{Region: &sqsQueueRegion}))

    for {
        // wait and receive messages
        req := &sqs.ReceiveMessageInput{
            QueueUrl: aws.String(sqsQueueUrl),
            WaitTimeSeconds: aws.Int64(20),
            MaxNumberOfMessages: aws.Int64(1),
        }
        resp, err := svc.ReceiveMessage(req)

        if err != nil {
            log.Printf("Cannot fetch messages from SQS queue %v: %v", sqsQueueUrl, err)
            time.Sleep(10 * time.Second)
            continue
        }

        for _, msg := range resp.Messages {
            // forward to Clair
            _, err = http.Post(fmt.Sprintf("%v/v1/layers", clairUrl), "application/json", strings.NewReader(*msg.Body))

            if err != nil {
                log.Printf("Couldn't forward message %v to Clair: %v", msg.MessageId, err)
                continue
            }

            log.Printf("Forwarded message %v to Clair", msg.MessageId)

            // delete message from queue
            _, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
                QueueUrl: aws.String(sqsQueueUrl),
                ReceiptHandle: msg.ReceiptHandle,
            })

            if err != nil {
                log.Printf("Could't delete message %v from queue: %v", msg.MessageId, err)
                continue
            }
        }
    }
}
