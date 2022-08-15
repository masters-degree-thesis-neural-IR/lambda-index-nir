package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"lambda-index-nir/service/application/service"
	"lambda-index-nir/service/infraestructure/dto"
	"lambda-index-nir/service/infraestructure/dydb"
	"log"
)

var TableName string
var AwsRegion string

func handler(ctx context.Context, event events.SQSEvent) error {

	//sess := session.Must(session.NewSessionWithOptions(session.Options{
	//	SharedConfigState: session.SharedConfigEnable,
	//}))

	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(AwsRegion)},
	)

	if err != nil {
		log.Fatalln("error...: ", err)
		return err
	}

	repository := dydb.NewIndexRepository(awsSession, TableName)
	service := service.NewIndexService(repository)

	for _, message := range event.Records {

		var mp map[string]string
		json.Unmarshal([]byte(message.Body), &mp)

		doc := &dto.Document{}
		err := json.Unmarshal([]byte(mp["Message"]), doc)

		if err != nil {
			log.Fatalln("error...: ", err)
			return err
		}

		err = service.CreateIndex(doc.Id, doc.Title, doc.Body)

		if err != nil {
			log.Fatalln("error...: ", err)
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
