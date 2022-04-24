package handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type IHandler interface {
	Handle(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayV2HTTPResponse, error)
}

type handler struct {
	s3Client   *s3.Client
	bucketName string
	host       string
}

type sample struct {
	Class  string    `json:"class"`
	Counts []float64 `json:"counts"`
}

type Body struct {
	ModelName  string            `json:"model_name"`
	Wavelength []float64         `json:"wavelength"`
	Samples    map[string]sample `json:"samples"`
}

type model struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type urls struct {
	Status   string `json:"status"`
	Classify string `json:"classify"`
}
type Response struct {
	Model model `json:"model"`
	URLs  urls  `json:"urls"`
}
