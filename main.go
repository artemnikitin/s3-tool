package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	logging = flag.Bool("log", false, "Enable logging")
	region  = flag.String("region", "us-east-1", "Set S3 region")
	bucket  = flag.String("bucket", "", "Name of bucket in S3")
	key     = flag.String("key", "", "Key for object in bucket")
	path    = flag.String("path", "", "Path for download")
)

func main() {
	flag.Parse()
	if *bucket == "" || *key == "" {
		fmt.Println("Please, specify valid parameters!")
		os.Exit(1)
	}

	client := s3.New(session.New(createConfig()))
	request, _ := client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(*bucket),
		Key:    aws.String(*key),
	})

	url, err := request.Presign(10 * time.Minute)
	if err != nil {
		log.Fatal("Failed to create pre-signed URL because of: ", err)
	}

	log.Println("Pre-signed URL:", url)
	downloadFile(url)
	log.Println("File is downloaded!")
}

func downloadFile(url string) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		log.Fatal("Failed to download a file by pre-signed URL because of: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to get a body from HTTP request because of: ", err)
	}
	err = ioutil.WriteFile(getFileName(*key), body, 0777)
	if err != nil {
		log.Fatal("Failed to create a file from a body of HTTP request because of: ", err)
	}
}

func getFileName(filepath string) string {
	index := strings.LastIndex(filepath, "/")
	if index == -1 {
		return filepath
	}
	return filepath[index+1:]
}

func createConfig() *aws.Config {
	config := aws.NewConfig()
	config.WithCredentials(credentials.NewEnvCredentials())
	config.WithRegion(*region)
	if *logging {
		config.WithLogLevel(aws.LogDebugWithHTTPBody)
	}
	return config
}
