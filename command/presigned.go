package command

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// GetPresignedURL generate pre-signed URL, which can be used for download file
func GetPresignedURL(session *session.Session, bucket, key string) (string, error) {
	var url string
	var err error
	client := s3.New(session)
	request, _ := client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	url, err = request.Presign(10 * time.Minute)
	if err != nil {
		log.Fatal("Failed to create pre-signed URL because of: ", err)
	}
	return url, err
}
