package incapsula

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func retryOnRateLimit(fn func() (*http.Response, error), initialDelay time.Duration) (*http.Response, error) {
	delay := initialDelay
	for attempt := 0; attempt < 5; attempt++ {
		resp, err := fn()
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		log.Printf("[DEBUG] API Security 429 response body: %s", string(body))
		if attempt < 4 {
			log.Printf("[WARN] API Security request rate limited (429), retrying in %s (attempt %d/5)", delay, attempt+1)
			time.Sleep(delay)
			delay *= 2
		}
	}
	return nil, fmt.Errorf("API Security request rate limited after 5 attempts")
}

func (c *Client) DoApiSecurityJsonRequest(method, url string, data []byte, operation string) (*http.Response, error) {
	return retryOnRateLimit(func() (*http.Response, error) {
		return c.DoJsonRequestWithHeaders(method, url, data, operation)
	}, 5*time.Second)
}

func (c *Client) DoApiSecurityFormDataRequest(method, url string, data []byte, contentType, operation string) (*http.Response, error) {
	return retryOnRateLimit(func() (*http.Response, error) {
		return c.DoFormDataRequestWithHeaders(method, url, data, contentType, operation)
	}, 5*time.Second)
}
