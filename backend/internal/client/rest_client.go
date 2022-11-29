package client

import (
	"bytes"
	"feedback/internal/logger"
	"io"
	"net/http"
)

var log = logger.Instance()

func Post(url string, reqBody []byte) (io.ReadCloser, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Error(err)
	}
	return resp.Body, err
}
