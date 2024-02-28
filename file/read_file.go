package file

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func isURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

func ReadFile(path string) (string, error) {
	var bytes []byte
	var err error
	if isURL(path) {
		bytes, err = get(path)
	} else {
		bytes, err = os.ReadFile(path)
	}
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func get(path string) ([]byte, error) {

	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func ReadFileLines(path string) ([]string, error) {
	str, err := ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(str, "\n")
	return lines, nil
}

func ReadFileLinesTrimSpace(path string) ([]string, error) {
	str, err := ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(str, "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return lines, nil
}
