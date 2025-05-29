package transport

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Middleware func(raw []byte) ([]byte, error)

type API struct {
	defaultHeaders map[string]string
	urlPath        string
	middlewares    []Middleware
}

func NewAPI(basePath string) *API {
	return &API{
		defaultHeaders: map[string]string{
			"Content-Type":    "application/json",
			"Accept":          "application/json",
			"Accept-Encoding": "gzip",
		},
		urlPath:     basePath,
		middlewares: []Middleware{},
	}
}

func (api *API) WithMiddleware(mw ...Middleware) {
	api.middlewares = append(api.middlewares, mw...)
}

func (api *API) Do(data []byte, headers map[string]string) ([]byte, error) {
	res, err := api.doCall(data, headers)
	if err != nil {
		return nil, err
	}

	for _, mw := range api.middlewares {
		res, err = mw(res)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (api *API) doCall(data []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, api.urlPath, bytes.NewBuffer(data))

	if err != nil {
		return nil, err
	}

	if headers == nil {
		headers = make(map[string]string)
	}

	api.mergeHeaders(req, headers)
	client := &http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2: false,
		},
		Timeout: 5 * time.Second,
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var bodyReader io.Reader = response.Body
	if strings.Contains(response.Header.Get("Content-Encoding"), "gzip") {
		var gz *gzip.Reader

		gz, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress response body: %w", err)
		}
		defer gz.Close()
		bodyReader = gz
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

func (api *API) mergeHeaders(r *http.Request, extraHeaders map[string]string) {
	for key, value := range api.defaultHeaders {
		r.Header.Set(key, value)
	}

	for key, value := range extraHeaders {
		r.Header.Set(key, value)
	}
}
