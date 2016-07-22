package command

import (
	"os"
	"strings"
	"testing"

	"github.com/artemnikitin/aws-config"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	s      = session.New(awsconfig.New())
	key    = os.Getenv("AWS_S3_KEY")
	bucket = os.Getenv("AWS_S3_BUCKET")
)

func TestGetPresignedURL(t *testing.T) {
	url, err := GetPresignedURL(s, bucket, key)
	if err != nil {
		t.Error("Can't get presigned URL: ", err)
	}
	if !strings.Contains(url, bucket) {
		t.Error("URL generated incorrectly: ", url)
	}
}
