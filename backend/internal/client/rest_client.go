package client

import (
	"bytes"
	"feedback/internal"
	"io"
	"net/http"
)

func Post(config *internal.Configuration, reqBody []byte) (io.ReadCloser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", config.OidcValidationUrl, bytes.NewBuffer(reqBody))
	req.Header.Add("Authorization", "Bearer "+config.UvsAuthToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}
