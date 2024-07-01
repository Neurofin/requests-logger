package serverConfigs

import (
	"application-manager/src/store/types"
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var SqsClient *sqs.Client

func SetupSqsClient() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	SqsClient = sqs.NewFromConfig(cfg)
	if SqsClient == nil {
		panic("Error setting up SQS client")
	}
}

func ListenToDocumentUploadEvents(queName string, eventHandler func(eventBody types.S3EventBody)) {
	if SqsClient == nil {
		panic("Sqs client not initialised")
	}

	result, err := SqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(queName),
	})
	if err != nil {
		panic(err)
	}

	queUrl := result.QueueUrl

	// Listner loop
	for {
		messagesPipe, err := SqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            queUrl,
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     0,
		})
		if err != nil {
			panic(err)
		}

		// Messages loop
		for _, message := range messagesPipe.Messages {

			body := *message.Body
			data := []byte(body)

			// If the data received is not in json format, ignore
			if !json.Valid(data) {
				continue
			}

			// Assumption here is that the event received will be of S3ObjectCreated type
			eventBody := types.S3EventBody{}
			if err := json.Unmarshal(data, &eventBody); err != nil {
				//TODO: Log Error
				continue
			}

			go eventHandler(eventBody)

			deleteInput := sqs.DeleteMessageInput{
				QueueUrl:      queUrl,
				ReceiptHandle: message.ReceiptHandle,
			}
			SqsClient.DeleteMessage(context.Background(), &deleteInput)
		}
	}
}
