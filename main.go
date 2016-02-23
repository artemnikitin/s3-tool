package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/artemnikitin/aws-config"
	"github.com/artemnikitin/s3-tool/command"
	"github.com/artemnikitin/s3-tool/logger"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	bucket     = flag.String("bucket", "", "Name of bucket in S3")
	key        = flag.String("key", "", "Key for object in bucket")
	downloadTo = flag.String("download", "", "Path for download")
	url        = flag.String("url", "", "Pre-signed URL for downloading")
	pathToFile = flag.String("path", "", "Path to file")
	rename     = flag.String("rename", "", "Set a new name for file")
	uploadTo   = flag.String("upload", "", "Set a specific path for a file inside S3 bucket")
	commands   = []string{"presigned", "download"}
)

func main() {
	comm, err := getCommand()
	logger.Process(err, "Incorrect command or command wasn't specified")
	log.Println("Command:", comm)

	flag.CommandLine.Parse(os.Args[2:])
	if *bucket == "" || *key == "" {
		fmt.Println("Please, specify valid parameters for command!")
		os.Exit(1)
	}

	session := session.New(awsconfig.New())

	switch comm {
	case "presigned":
		link, err := command.GetPresignedURL(session, *bucket, *key)
		logger.Process(err, "Can't generate pre-signed S3 URL")
		log.Println("Pre-signed URL:", link)
	case "download":
		log.Println("Start downloading file...")
		dest := *key
		if *downloadTo != "" {
			dest = *downloadTo + "/" + *key
		}
		if *url != "" {
			command.DownloadFile(*url, dest)
		} else {
			command.Download(session, *bucket, *key, dest)
		}
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
	if !validCommand(arg) {
		err = errors.New("Unexisted command")
	}
	return arg, err
}

func validCommand(command string) bool {
	for _, comm := range commands {
		if comm == command {
			return true
		}
	}
	return false
}
