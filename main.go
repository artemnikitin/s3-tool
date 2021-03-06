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
	bucket         = flag.String("bucket", "", "Name of bucket in S3")
	key            = flag.String("key", "", "Key for object in bucket")
	url            = flag.String("url", "", "Pre-signed URL for downloading")
	pathToFile     = flag.String("path", "", "Path to file")
	keepRootFolder = flag.Bool("keepRootFolder", false, "Keep root folder in S3 bucket")
	commands       = []string{"presigned", "download", "upload"}
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	comm, err := getCommand()
	logger.Process(err, "")
	fmt.Println("Command:", comm)

	err = flag.CommandLine.Parse(os.Args[2:])
	logger.Process(err, "Can't parse arguments")
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
		if *pathToFile != "" {
			dest = *pathToFile + "/" + *key
		}
		if *url != "" {
			command.DownloadFile(*url, dest)
		} else {
			command.Download(session, *bucket, *key, dest)
		}
		fmt.Println("File is downloaded!")
	case "upload":
		if *pathToFile == "" {
			fmt.Println("Can't proceed without path to file")
			os.Exit(1)
		}
		file, err := os.Open(*pathToFile)
		logger.Process(err, "Failed to open a file")
		defer file.Close()
		info, err := file.Stat()
		logger.Process(err, "Failed to get info about file")
		switch mode := info.Mode(); {
		case mode.IsDir():
			command.UploadDirectory(session, *bucket, *key, *pathToFile, *keepRootFolder)
		case mode.IsRegular():
			command.UploadFile(session, *bucket, *key, file)
		}
	}
}

func getCommand() (string, error) {
	var err error
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println("Please, provide valid command.")
			os.Exit(1)
		}
	}()
	arg := os.Args[1]
	if !validCommand(arg) {
		err = errors.New("Unexisted command")
	}
	return arg, err
}

func validCommand(command string) bool {
	result := false
	for _, comm := range commands {
		if comm == command {
			result = true
			break
		}
	}
	return result
}
