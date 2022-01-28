package main

import (
	"context"
	"log"
	"os"

	// aws
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	// local
	"github.com/kinloth/open-c-lifs/apigateway/new_model/handler"
)

func main() {
	// environment variables
	bucketName := os.Getenv("BUCKET_NAME")
	host := os.Getenv("SERVER_HOST")

	// load config
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load aws config: %s", err.Error())
	}

	// load clients
	s3Client := s3.NewFromConfig(cfg)

	// create handler
	h, err := handler.New(s3Client, bucketName, host)
	if err != nil {
		log.Fatalf("failed to create handler: %s", err.Error())
	}

	// run lambda
	lambda.Start(h.Handle)
}
