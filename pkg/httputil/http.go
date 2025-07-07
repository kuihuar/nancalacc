package httputil

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

func PostJSON(uri string, jsonData []byte, timeout time.Duration) ([]byte, error) {
	httpClient := &http.Client{Timeout: timeout}
	resp, err := httpClient.Post(uri, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
