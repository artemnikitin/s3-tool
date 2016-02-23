package command

import (
	"time"

	"github.com/artemnikitin/s3-tool/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// GetPresignedURL generate pre-signed URL, which can be used for download file
func GetPresignedURL(session *session.Session, bucket, key string) (string, error) {
	client := s3.New(session)
	request, _ := client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	url, err := request.Presign(10 * time.Minute)
	logger.Process(err, "Failed to create pre-signed URL")
	return url, err
}
