package command

import "testing"

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
	cases := map[string]struct{ Path, Folder, Result string }{
		"without slash":    {"/ddd/bbb/f/12.a", "f", "/f/12.a"},
		"with slash":       {"/ddd/ccc/k/12.b/", "k", "/k/12.b/"},
		"unexisted folder": {"/ddd/ccc/c/12.b/", "e", ""},
		"empty input":      {"", "", ""},
	}

	for key, value := range cases {
		res := getPathInsideFolder(value.Path, value.Folder)
		if res != value.Result {
			t.Errorf("For case: %s, actual: %s, expected: %s", key, res, value.Result)
		}
	}
}
