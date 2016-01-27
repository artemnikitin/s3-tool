package command

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// DownloadFile using result of GetPresignedURL for downloading file by pre-signed URL
func DownloadFile(url, path string) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		log.Fatal("Failed to download a file by pre-signed URL because of: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to get a body from HTTP request because of: ", err)
	}
	err = ioutil.WriteFile(path, body, 0777)
	if err != nil {
		log.Fatal("Failed to create a file from a body of HTTP request because of: ", err)
	}
}

// Download using to download file from S3 by bucket and key
func Download(session *session.Session, bucket, key, path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal("Can't create file because of:", err)
	}
	defer file.Close()

	client := s3manager.NewDownloader(session)
	_, err = client.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Fatal("Can't retrieve object from S3 because of:", err)
	}
}
