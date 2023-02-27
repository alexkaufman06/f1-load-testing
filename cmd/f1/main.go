package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/form3tech-oss/f1/v2/pkg/f1"
	"github.com/form3tech-oss/f1/v2/pkg/f1/testing"
	"net/http"
	"time"
)

func main() {
	f := f1.New()
	f.Add("testScenario", testScenario)
	f.Execute()
}

func testScenario(t *testing.T) testing.RunFn {
	// Configure an SQS client
	goawsEndpoint := "http://localhost:4100"
	region := "us-west-2"
	s, err := session.NewSession(&aws.Config{
		Endpoint:    &goawsEndpoint,
		Region:      &region,
		Credentials: credentials.NewStaticCredentials("foo", "bar", ""),
	})
	if err != nil {
		t.Require().NoError(err)
	}
	sqsClient := sqs.New(s)
	queueName := "test-queue"
	queueUrlResponse, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		t.Require().NoError(err)
	}
	queueUrl := queueUrlResponse.QueueUrl

	// Consume SQS messages from the queue and put them in a channel
	messagesChan := make(chan string, 100)
	stopChan := make(chan bool)
	go func() {
		for {
			select {
			case <-stopChan:
				return
			default:
				messages, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
					QueueUrl: queueUrl,
				})
				t.Require().NoError(err)
				for _, message := range messages.Messages {
					if message.Body != nil {
						messagesChan <- *message.Body
					}
				}
			}
		}
	}()
	t.Cleanup(func() {
		stopChan <- true
	})

	runFn := func(t *testing.T) {
		// Test iteration code
		res, err := http.Post("http://localhost:8080/payments", "application/json", nil)
		t.Require().NoError(err)
		t.Require().Equal(http.StatusAccepted, res.StatusCode)
		timer := time.NewTimer(10 * time.Second)
		for {
			select {
			case <-timer.C:
				t.Require().Fail("no message received after timeout")
				return
			case <-messagesChan:
				t.Logger().Info("message received, iteration success")
				return
			}
		}
	}
	return runFn
}
