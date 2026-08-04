package incapsula

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func makeResponse(statusCode int) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       http.NoBody,
	}
}

func TestRetryOnRateLimit_SuccessFirstAttempt(t *testing.T) {
	calls := 0
	fn := func() (*http.Response, error) {
		calls++
		return makeResponse(http.StatusOK), nil
	}

	resp, err := retryOnRateLimit(fn, 0)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestRetryOnRateLimit_SuccessAfterOneRetry(t *testing.T) {
	calls := 0
	fn := func() (*http.Response, error) {
		calls++
		if calls == 1 {
			return makeResponse(http.StatusTooManyRequests), nil
		}
		return makeResponse(http.StatusOK), nil
	}

	resp, err := retryOnRateLimit(fn, time.Millisecond)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}

func TestRetryOnRateLimit_SuccessAfterFourRetries(t *testing.T) {
	calls := 0
	fn := func() (*http.Response, error) {
		calls++
		if calls <= 4 {
			return makeResponse(http.StatusTooManyRequests), nil
		}
		return makeResponse(http.StatusOK), nil
	}

	resp, err := retryOnRateLimit(fn, time.Millisecond)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if calls != 5 {
		t.Errorf("expected 5 calls, got %d", calls)
	}
}

func TestRetryOnRateLimit_ExhaustsRetries(t *testing.T) {
	calls := 0
	fn := func() (*http.Response, error) {
		calls++
		return makeResponse(http.StatusTooManyRequests), nil
	}

	resp, err := retryOnRateLimit(fn, time.Millisecond)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "API Security request rate limited after 5 attempts" {
		t.Errorf("unexpected error message: %s", err)
	}
	if resp != nil {
		t.Errorf("expected nil response, got %v", resp)
	}
	if calls != 5 {
		t.Errorf("expected 5 calls, got %d", calls)
	}
}

func TestRetryOnRateLimit_NetworkError(t *testing.T) {
	calls := 0
	networkErr := fmt.Errorf("connection refused")
	fn := func() (*http.Response, error) {
		calls++
		return nil, networkErr
	}

	resp, err := retryOnRateLimit(fn, time.Millisecond)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != networkErr {
		t.Errorf("expected network error, got: %s", err)
	}
	if resp != nil {
		t.Errorf("expected nil response, got %v", resp)
	}
	if calls != 1 {
		t.Errorf("expected 1 call (no retry on network error), got %d", calls)
	}
}

func TestRetryOnRateLimit_NonRateLimitError(t *testing.T) {
	calls := 0
	fn := func() (*http.Response, error) {
		calls++
		return makeResponse(http.StatusInternalServerError), nil
	}

	resp, err := retryOnRateLimit(fn, time.Millisecond)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
	if calls != 1 {
		t.Errorf("expected 1 call (no retry on non-429), got %d", calls)
	}
}
