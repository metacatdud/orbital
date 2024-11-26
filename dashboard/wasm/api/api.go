package api

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"time"
)

type API struct {
	defaultHeaders map[string]string
	urlPath        string
}

func NewAPI(basePath string) *API {
	return &API{
		defaultHeaders: map[string]string{
			"Content-Type":    "application/json",
			"Accept":          "application/json",
			"Accept-Encoding": "gzip, br",
		},
		urlPath: basePath,
	}
}

func (api *API) Do(data []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, api.urlPath, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	if headers == nil {
		headers = make(map[string]string)
	}

	reqHeaders := api.mergeHeaders(headers)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for k, v := range reqHeaders {
		req.Header.Set(k, v)
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var bodyReader io.ReadCloser = response.Body
	if response.Header.Get("Content-Encoding") == "gzip" {

		bodyReader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress response body: %w", err)
		}

		defer bodyReader.Close()
	}

	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status code: %d, body: %s", response.StatusCode, body)
	}

	return body, nil
}

func (api *API) mergeHeaders(extraHeaders map[string]string) map[string]string {
	headersMerged := make(map[string]string)
	for key, value := range api.defaultHeaders {
		headersMerged[key] = value
	}

	for key, value := range extraHeaders {
		headersMerged[key] = value
	}

	return headersMerged
}
