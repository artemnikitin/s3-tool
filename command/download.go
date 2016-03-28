package command

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/artemnikitin/s3-tool/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// DownloadFile using result of GetPresignedURL for downloading file by pre-signed URL
func DownloadFile(url, path string) {
	resp, err := http.DefaultClient.Get(url)
	logger.Process(err, "Failed to download a file by pre-signed URL")
	defer func() {
		// Drain and close the body to let the Transport reuse the connection
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	logger.Process(err, "Failed to get a body from HTTP request")
	err = ioutil.WriteFile(path, body, 0777)
	logger.Process(err, "Failed to create a file from a body of HTTP request")
}

// Download using to download file from S3 by bucket and key
func Download(session *session.Session, bucket, key, path string) {
	file, err := os.Create(path)
	logger.Process(err, "Can't create file")
	defer file.Close()
	client := s3manager.NewDownloader(session)
	_, err = client.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	logger.Process(err, "Can't retrieve object from S3")
}
