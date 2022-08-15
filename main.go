package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"lambda-index-nir/service/application/service"
	"lambda-index-nir/service/infraestructure/dto"
	"lambda-index-nir/service/infraestructure/dydb"
	"log"
)

var TopicArn string
var TableName string
var AwsRegion string

func makeBody(body string) (dto.Document, error) {
	var doc dto.Document
	err := json.Unmarshal([]byte(body), &doc)

	if err != nil {
		return doc, err
	}

	return doc, nil

}

func handler(ctx context.Context, event events.SQSEvent) error {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	repository := dydb.NewIndexRepository(sess, "NIR_Index")
	service := service.NewIndexService(repository)

	for _, message := range event.Records {

		doc := dto.Document{}
		err := json.Unmarshal([]byte(message.Body), &doc)

		if err != nil {
			log.Fatalln("error...: ", err)
			return err
		}

		err = service.CreateIndex(doc.Id, doc.Title, doc.Body)

		if err != nil {
			log.Fatalln("error...: ", err)
			return err
		}

		println("****** Sucesso *******")
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
	}

	return nil
}

func main() {
	//AwsRegion = "us-east-1"
	//TopicArn = "arn:aws:sns:us-east-1:149501088887:mestrado-document-created" //os.Getenv("BAR")
	TableName = "NIR_Index"
	lambda.Start(handler)
}
