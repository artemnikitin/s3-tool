package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/artemnikitin/s3-tool/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// UploadFile will upload file to specific S3 bucket
func UploadFile(session *session.Session, bucket, key string, file io.Reader) {
	service := s3manager.NewUploader(session)
	resp, err := service.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	logger.Process(err, "Can't upload file")
	fmt.Println("File was successfully uploaded! Location:", resp.Location)
}

// UploadDirectory will upload directory and all it's content while keeping it structure
func UploadDirectory(session *session.Session, bucket, param string) {
	var wg sync.WaitGroup
	err := filepath.Walk(param, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			file, err := os.Open(path)
			if err == nil {
				path := getPathInsideFolder(path, getFolderName(path))
				wg.Add(1)
				go func() {
					UploadFile(session, bucket, createKey(param, path), file)
					wg.Done()
					file.Close()
				}()
			}
		}
		return nil
	})
	wg.Wait()
	logger.Process(err, "Can't process directory")
	fmt.Println("Directory was successfully uploaded!")
}

func getPathInsideFolder(path, folder string) string {
	pos := strings.Index(path, folder)
	var result string
	if pos != -1 {
		result = string(path[pos-1:])
	}
	return result
}

func getFolderName(filepath string) string {
	var result string
	if endWith(filepath, "/") {
		pos := strings.LastIndex(string(filepath[:len(filepath)-1]), "/")
		result = string(filepath[pos+1 : len(filepath)-1])
	} else {
		pos := strings.LastIndex(filepath, "/")
		result = string(filepath[pos+1:])
	}
	return result
}

func startWith(original, substring string) bool {
	if len(substring) > len(original) {
		return false
	}
	str := string(original[0:len(substring)])
	return str == substring
}

func endWith(original, substring string) bool {
	if len(substring) > len(original) {
		return false
	}
	str := string(original[len(original)-len(substring):])
	return str == substring
}

func createKey(param, path string) string {
	var buffer bytes.Buffer
	buffer.WriteString(param)
	if param == "/" {
		if startWith(path, "/") {
			return path
		}
		buffer.WriteString(path)
	} else {
		if !endWith(param, "/") && !startWith(path, "/") {
			buffer.WriteString("/")
		}
		if endWith(param, "/") && startWith(path, "/") {
			buffer.WriteString(string(path[1:]))
		} else {
			buffer.WriteString(path)
		}
	}
	return buffer.String()
}
