package utils

import (
	"io"
	"net/http"
	"regexp"
)

var LinkRegex = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)

func FromWeb(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	return io.ReadAll(res.Body)
}

func FromWebString(url string) (string, error) {
	body, err := FromWeb(url)
	if err != nil {
		return "", err
	}

	return string(body), nil
}