package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"lambda-index-nir/service/application/service"
	"lambda-index-nir/service/infraestructure/dto"
	zapplog "lambda-index-nir/service/infraestructure/log"
	"lambda-index-nir/service/infraestructure/memory"
)

var TableName string
var AwsRegion string

func handler(ctx context.Context, event events.SQSEvent) error {

	logger := zapplog.NewZapLogger()
	repository := memory.NewSpeedupRepository()
	service := service.NewIndexService(logger, repository)

	logger.Info("Lambda Accepted Request")

	for _, message := range event.Records {

		var mp map[string]string
		json.Unmarshal([]byte(message.Body), &mp)

		doc := &dto.Document{}
		err := json.Unmarshal([]byte(mp["Message"]), doc)

		logger.Info("Documento recebido")
		logger.Info(doc)

		if err != nil {
			logger.Fatal(err.Error())
			return err
		}

		err = service.CreateIndex(doc.Id, doc.Title, doc.Body)

		if err != nil {
			logger.Fatal(err.Error())
			return err
		}

	}

	return nil
}

func main() {
	//TopicArn = "arn:aws:sns:us-east-1:149501088887:mestrado-document-created" //os.Getenv("BAR")
	AwsRegion = "us-east-1"
	TableName = "NIR_Index"
	lambda.Start(handler)
}
