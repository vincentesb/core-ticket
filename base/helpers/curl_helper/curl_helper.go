package curl_helper

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// CurlRequest performs an HTTP request with the specified parameters and returns the decoded response object.
func CurlRequest[T any](url string, headers map[string]string, body []byte, method string, requestTimeout int, returnJson bool) (*T, error) {
	timeOut := 600000
	if requestTimeout > 0 {
		timeOut = requestTimeout
	}
	client := &http.Client{
		Timeout: time.Duration(timeOut) * time.Millisecond,
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	respCode := resp.StatusCode

	if !returnJson {
		if respCode >= 200 && respCode <= 299 {
			return nil, nil
		} else {
			return nil, nil
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Decode JSON response
	var obj T
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		return nil, err
	}

	return &obj, nil
}
