package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	// aws
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func New(s3Client *s3.Client, bucketName, host string) (IHandler, error) {
	if s3Client == nil {
		return nil, fmt.Errorf("s3Client is nil")
	}
	if bucketName == "" {
		return nil, fmt.Errorf("bucketName is empty")
	}
	if host == "" {
		return nil, fmt.Errorf("host is empty")
	}

	h := handler{
		s3Client:   s3Client,
		bucketName: bucketName,
		host:       host,
	}
	return h, nil
}

func (h handler) putObject(ctx context.Context, body *Body, id string) error {
	bodyMarshalled, err := json.Marshal(body)
	if err != nil {
		return err
	}
	input := s3.PutObjectInput{
		Bucket:        &h.bucketName,
		Key:           aws.String(fmt.Sprintf("models/%s/training_data/input.json", id)),
		Body:          bytes.NewReader(bodyMarshalled),
		ContentLength: int64(len(bodyMarshalled)),
		ContentMD5:    aws.String(GetMD5Hash(bodyMarshalled)),
		ContentType:   aws.String("application/json"),
		Expires:       getExpires(),
		Metadata:      getMetadata(),
		Tagging:       aws.String(getTagging()),
	}
	output, err := h.s3Client.PutObject(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to put object: %s", err.Error())
	}
	fmt.Printf("%+v\n", output)
	return nil
}

func (h handler) getURL(action, id string) string {
	return fmt.Sprintf("https://%s/openc-lifs/models/%s/%s", h.host, id, action)
}

func (h handler) getSuccessResponse(id, name string) (events.APIGatewayProxyResponse, error) {
	response := Response{
		Model: model{
			ID:   id,
			Name: name,
		},
		URLs: urls{
			Status:   h.getURL("status", id),
			Classify: h.getURL("classify", id),
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		return genResponse(err.Error(), http.StatusInternalServerError)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(data),
		StatusCode: http.StatusAccepted,
	}, nil
}

func (h handler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	in, err := formatBody(request.Body)
	if err != nil {
		return genResponse(err.Error(), http.StatusBadRequest)
	}

	id := uuid.New().String()
	err = h.putObject(ctx, in, id)
	if err != nil {
		return genResponse(err.Error(), http.StatusInternalServerError)
	}

	return h.getSuccessResponse(id, in.ModelName)
}
