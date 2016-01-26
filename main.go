package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/artemnikitin/s3-tool/command"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	logging  = flag.Bool("log", false, "Enable logging")
	region   = flag.String("region", "us-east-1", "Set S3 region")
	bucket   = flag.String("bucket", "", "Name of bucket in S3")
	key      = flag.String("key", "", "Key for object in bucket")
	path     = flag.String("path", "", "Path for download")
	commands = []string{"presigned"}
)

func main() {
	flag.Parse()
	comm, err := getCommand()
	if err != nil {
		log.Fatal("Incorrect command or command wasn't specified")
	}
	log.Println("Command:", comm)
	if *bucket == "" || *key == "" {
		fmt.Println("Please, specify valid parameters for command!")
		os.Exit(1)
	}

	session := session.New(createConfig())

	switch comm {
	case "presigned":
		url, err := command.GetPresignedURL(session, *bucket, *key)
		if err != nil {
			log.Fatal("Can't generate pre-signed S3 URL because of:", err)
		}
		log.Println("Pre-signed URL:", url)
		dest := *key
		if *path != "" {
			dest = *path + "/" + *key
		}
		command.DownloadFile(url, dest)
		fmt.Println("File is downloaded!")
	}
}

func getCommand() (string, error) {
	var err error
	defer func() {
		rec := recover()
		if rec != nil {
			log.Println("Catch panic:", rec)
			err = errors.New("Received panic while processing command")
		}
	}()
	arg := os.Args[1]
	if !validateCommand(arg) {
		err = errors.New("Unexisted command")
	}
	return arg, err
}

func validateCommand(command string) bool {
	valid := false
	for _, comm := range commands {
		if comm == command {
			valid = true
			break
		}
	}
	return valid
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
