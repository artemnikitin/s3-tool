package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/artemnikitin/s3-tool/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var contentTypes = map[string]string{
	"txt":  "text/plain",
	"json": "application/json",
	"xml":  "application/xml",
	"pdf":  "application/pdf",
	"html": "text/html",
	"htm":  "text/html",
	"css":  "text/css",
	"js":   "application/javascript",
	"bmp":  "image/bmp",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"tiff": "image/tiff",
	"gif":  "image/gif",
}

// UploadFile will upload file to specific S3 bucket
func UploadFile(session *session.Session, bucket, key string, file *os.File) {
	service := s3manager.NewUploader(session)
	resp, err := service.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(getContentType(file)),
	})
	logger.Process(err, "Can't upload file")
	fmt.Println("File was successfully uploaded! Location:", resp.Location)
}

// UploadDirectory will upload directory and all it's content while keeping it structure
func UploadDirectory(session *session.Session, bucket, key, dir string, keep bool) {
	var wg sync.WaitGroup
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			file, err := os.Open(path)
			if err == nil {
				path := getPathInsideFolder(path, getFolderName(dir), keep)
				go func() {
					wg.Add(1)
					UploadFile(session, bucket, key+path, file)
					file.Close()
					wg.Done()
				}()
			}
		}
		return nil
	})
	wg.Wait()
	logger.Process(err, "Can't process directory")
	fmt.Println("Directory was successfully uploaded!")
}

func getPathInsideFolder(path, folder string, keep bool) string {
	if path == "" || folder == "" {
		return ""
	}
	pos := -1
	var result string
	if keep {
		pos = strings.Index(path, folder)
	} else {
		temp := strings.Index(path, folder) + len(folder)
		if temp <= len(path) && temp > 0 {
			pos = temp
		}
	}
	if pos != -1 {
		if keep {
			result = strings.Join([]string{"/", path[pos:]}, "")
		} else {
			result = path[pos:]
		}
	}
	return result
}

func getFolderName(filepath string) string {
	if filepath == "" {
		return ""
	}
	var result string
	if strings.HasSuffix(filepath, "/") {
		pos := strings.LastIndex(string(filepath[:len(filepath)-1]), "/")
		result = string(filepath[pos+1 : len(filepath)-1])
	} else {
		pos := strings.LastIndex(filepath, "/")
		result = string(filepath[pos+1:])
	}
	return result
}

func getContentType(file *os.File) string {
	result := "binary/octet-stream"
	name := file.Name()
	pos := strings.LastIndex(name, ".")
	if pos != -1 {
		v, ok := contentTypes[name[pos+1:]]
		if ok {
			result = v
		}
	}
	return result
}
