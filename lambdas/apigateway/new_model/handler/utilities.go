package handler

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

func formatBody(body string) (*Body, error) {
	fmt.Println("Body: ", body)

	var in Body
	err := json.Unmarshal([]byte(body), &in)
	if err != nil {
		return nil, err
	}

	if in.ModelName == "" {
		return nil, fmt.Errorf("model name is required")
	}
	if len(in.Wavelength) == 0 {
		return nil, fmt.Errorf("wavelength is required")
	}
	if len(in.Samples) == 0 {
		return nil, fmt.Errorf("samples is required")
	}
	if len(in.Samples) != len(in.Wavelength) {
		return nil, fmt.Errorf("samples and wavelength must be the same length")
	}

	return &in, nil
}

type ResponseError struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"error_code"`
}

func genResponse(body string, statusCode int) (events.APIGatewayProxyResponse, error) {
	r := ResponseError{
		Message:   body,
		ErrorCode: statusCode,
	}
	marshalled, _ := json.Marshal(r)

	return events.APIGatewayProxyResponse{
		Body:       string(marshalled),
		StatusCode: statusCode,
	}, nil
}

func GetMD5Hash(content []byte) string {
	md := md5.New() //nolint:gosec
	md.Write(content)
	sum := md.Sum(nil)
	contentMd5 := make([]byte, base64.StdEncoding.EncodedLen(len(sum)))
	base64.StdEncoding.Encode(contentMd5, sum)
	return string(contentMd5)
}

func getExpires() *time.Time {
	exp := time.Now().Add(time.Hour * 24)
	return &exp
}

func getMetadata() map[string]string { return nil }
func getTagging() string             { return "" }
