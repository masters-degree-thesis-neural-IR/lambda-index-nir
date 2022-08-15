package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"lambda-index-nir/service/infraestructure/dto"
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

	for _, message := range event.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
	}

	return nil

	//if req.HTTPMethod != "POST" {
	//	return events.APIGatewayProxyResponse{
	//		StatusCode: http.StatusBadRequest,
	//		Headers:    map[string]string{"Content-Type": "text/plain; charset=utf-8"},
	//		Body:       "Invalid HTTP Method",
	//	}, nil
	//}
	//
	//document, err := makeBody(req.Body)
	//
	//if err != nil {
	//	return ErrorHandler(err), nil
	//}
	//
	//awsSession, err := session.NewSession(&aws.Config{
	//	Region: aws.String(AwsRegion)},
	//)
	//
	//if err != nil {
	//	return ErrorHandler(err), nil
	//}
	//
	//repository := dydb.NewDocumentRepository(awsSession, TableName)
	//documentEvent := sns.NewDocumentEvent(awsSession, TopicArn)
	//documentService := service.NewDocumentService(documentEvent, repository)
	//err = documentService.CreateDocument(document.Title, document.Body)
	//
	//if err != nil {
	//	return ErrorHandler(err), nil
	//}
	//
	//return events.APIGatewayProxyResponse{
	//	StatusCode: http.StatusCreated,
	//	Headers:    map[string]string{"Content-Type": "text/plain; charset=utf-8"},
	//	Body:       "Document created",
	//}, nil

}

func main() {
	AwsRegion = "us-east-1"
	TopicArn = "arn:aws:sns:us-east-1:149501088887:mestrado-document-created" //os.Getenv("BAR")
	TableName = "NIR_Document"
	lambda.Start(handler)
}
