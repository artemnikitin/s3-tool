package command

import (
	"os"
	"testing"
)

func TestGetFolderName(t *testing.T) {
	cases := map[string]struct{ In, Out string }{
		"without slash": {"/ddd/bbb", "bbb"},
		"with slash":    {"/ddd/ccc/", "ccc"},
		"empty input":   {"", ""},
	}

	for key, value := range cases {
		res := getFolderName(value.In)
		if res != value.Out {
			t.Errorf("For case: %s, actual: %s, expected: %s", key, res, value.Out)
		}
	}
}

func TestGetPathInsideFolder(t *testing.T) {
	cases := map[string]struct {
		Path, Folder   string
		KeepRootFolder bool
		Result         string
	}{
		"without slash":               {"/ddd/bbb/f/12.a", "f", true, "/f/12.a"},
		"with slash":                  {"/ddd/ccc/k/12.b/", "k", true, "/k/12.b/"},
		"long folder name":            {"/ddd/ccc/kkkkk/12.b/", "kkkkk", true, "/kkkkk/12.b/"},
		"unexisted folder":            {"/ddd/ccc/c/12.b/", "e", true, ""},
		"empty input":                 {"", "", true, ""},
		"without slash with false":    {"/ddd/bbb/f/12.a", "f", false, "/12.a"},
		"with slash with false":       {"/ddd/ccc/k/12.b/", "k", false, "/12.b/"},
		"long folder name with false": {"/ddd/ccc/kkkkk/12.b/", "kkkkk", false, "/12.b/"},
		"unexisted folder with false": {"/ddd/ccc/c/12.b/", "e", false, ""},
	}

	for k, v := range cases {
		res := getPathInsideFolder(v.Path, v.Folder, v.KeepRootFolder)
		if res != v.Result {
			t.Errorf("For case: %s, actual: %s, expected: %s", k, res, v.Result)
		}
	}
}

func TestGetContentType(t *testing.T) {
	cases := map[string]struct{ Path, Result string }{
		"unknown type": {"upload_test.go", "binary/octet-stream"},
	}

	for k, v := range cases {
		file, _ := os.Open(v.Path)
		res := getContentType(file)
		if res != v.Result {
			t.Errorf("For case: %s, actual: %s, expected: %s", k, res, v.Result)
		}
	}
}
