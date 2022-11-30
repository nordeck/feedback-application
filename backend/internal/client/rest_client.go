package client

import (
	"bytes"
	"io"
	"net/http"
)

func Post(url string, reqBody []byte) (io.ReadCloser, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}
