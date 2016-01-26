package command

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// DownloadFile using result of GetPresignedURL for downloading file by pre-signed URL
func DownloadFile(url, path string) bool {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		log.Fatal("Failed to download a file by pre-signed URL because of: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to get a body from HTTP request because of: ", err)
	}
	err = ioutil.WriteFile(getFileName(path), body, 0777)
	if err != nil {
		log.Fatal("Failed to create a file from a body of HTTP request because of: ", err)
	}
	return true
}

func getFileName(filepath string) string {
	index := strings.LastIndex(filepath, "/")
	if index == -1 {
		return filepath
	}
	return filepath[index+1:]
}
