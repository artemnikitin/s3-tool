package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/artemnikitin/aws-config"
	"github.com/artemnikitin/s3-tool/command"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	bucket   = flag.String("bucket", "", "Name of bucket in S3")
	key      = flag.String("key", "", "Key for object in bucket")
	path     = flag.String("path", "", "Path for download")
	url      = flag.String("url", "", "Pre-signed URL for downloading")
	filespath  = flag.String("path", "", "Path to file")
	rename     = flag.String("rename", "", "Set a new name for file")
	uploadpath = flag.String("uploadto", "", "Set a specific path for a file inside S3 bucket")
	commands = []string{"presigned", "download"}
)

func main() {
	comm, err := getCommand()
	if err != nil {
		log.Fatal("Incorrect command or command wasn't specified")
	}
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
		if err != nil {
			log.Fatal("Can't generate pre-signed S3 URL because of:", err)
		}
		log.Println("Pre-signed URL:", link)
	case "download":
		log.Println("Start downloading file...")
		dest := *key
		if *path != "" {
			dest = *path + "/" + *key
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
