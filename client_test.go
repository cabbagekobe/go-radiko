package radiko

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"
	"time"
)

// Should restore defaultHTTPClient if SetHTTPClient is called.
func teardownHTTPClient() {
	SetHTTPClient(&http.Client{Timeout: defaultHTTPTimeout})
}

func TestNew(t *testing.T) {
	_, err := New("")
	if err != nil {
		t.Fatalf("Failed to construct client: %s", err)
	}
}

func TestEmptyHTTPClient(t *testing.T) {
	var c *http.Client

	SetHTTPClient(c)
	defer teardownHTTPClient()

	client, err := New("")
	if err == nil {
		t.Errorf(
			"Should detect HTTPClient is nil.\nclient: %v", client)
	}
}

func TestNewRequest(t *testing.T) {
	client, err := New("")
	if err != nil {
		t.Fatalf("Failed to construct client: %s", err)
	}

	ctx := context.Background()
	_, err = client.newRequest(ctx, "GET", "", &Params{})
	if err != nil {
		t.Error(err)
	}
}

func TestNewRequestWithAuthToken(t *testing.T) {
	const expected = "auth_token"

	client, err := New(expected)
	if err != nil {
		t.Errorf("Failed to construct client: %s", err)
	}

	req, err := client.newRequest(context.Background(), "GET", "", &Params{
		setAuthToken: true,
	})
	if err != nil {
		t.Error(err)
	}
	if actual := req.Header.Get(radikoAuthTokenHeader); actual != expected {
		t.Errorf("expected %s, but %s.", expected, actual)
	}
}

func TestNewRequestWithContext(t *testing.T) {
	client, err := New("")
	if err != nil {
		t.Fatalf("Failed to construct client: %s", err)
	}

	timeout := 100 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = client.newRequest(ctx, "GET", "", &Params{})
	if err != nil {
		t.Error(err)
	}

	select {
	case <-time.After(3 * time.Second):
		t.Fatalf("context: %v", ctx)
	case <-ctx.Done():
	}

	if ctx.Err() == nil {
		t.Errorf("Shoud detect the context deadline exceeded.\n%v", ctx)
	}
}

func TestNewRequestWithEmptyContext(t *testing.T) {
	client, err := New("")
	if err != nil {
		t.Fatalf("Failed to construct client: %s", err)
	}

	var ctx context.Context
	_, err = client.newRequest(ctx, "GET", "", &Params{})
	if err == nil {
		t.Error("Should detect empty context.")
	}
}

func TestSetAuthTokenHeader(t *testing.T) {
	client, err := New("")
	if err != nil {
		t.Errorf("Failed to construct client: %s", err)
	}

	const expected = "test_token"
	client.setAuthTokenHeader(expected)
	if expected != client.authTokenHeader {
		t.Errorf("expected %s, but %s", expected, client.authTokenHeader)
	}
}

func TestSetJar(t *testing.T) {
	client, err := New("")
	if err != nil {
		t.Errorf("Failed to construct client: %s", err)
	}

	if client.httpClient.Jar != nil {
		t.Error("httpClient.Jar should be nil.")
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	client.SetJar(jar)
	if client.httpClient.Jar == nil {
		t.Error("httpClient.Jar is nil.")
	}
}

func TestDo(t *testing.T) {
	client, err := New("")
	if err != nil {
		t.Errorf("Failed to construct client: %s", err)
	}

	ctx := context.Background()
	req, err := client.newRequest(ctx, "GET", "", &Params{})
	if err != nil {
		t.Error(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	const expected = 200
	if actual := resp.StatusCode; actual != expected {
		t.Errorf("expected %d, but StatusCode is %d.", expected, actual)
	}
}

func TestSetHTTPClient(t *testing.T) {
	const expected = 1 * time.Second

	SetHTTPClient(&http.Client{Timeout: expected})
	defer teardownHTTPClient()

	client, err := New("")
	if err != nil {
		t.Errorf("Failed to construct client: %s", err)
	}
	if client.httpClient.Timeout != expected {
		t.Errorf("expected %d, but %d", expected, client.httpClient.Timeout)
	}
}

func TestAPIPath(t *testing.T) {
	const path = "test"
	var apiEndpoint string

	apiEndpoint = apiPath(apiV2, path)
	if !(strings.HasPrefix(apiEndpoint, apiV2+"/") && strings.HasSuffix(apiEndpoint, "/"+path)) {
		t.Errorf("invalid apiEndpoint: %s", apiEndpoint)
	}

	apiEndpoint = apiPath(apiV3, path)
	if !(strings.HasPrefix(apiEndpoint, apiV3+"/") && strings.HasSuffix(apiEndpoint, "/"+path)) {
		t.Errorf("invalid apiEndpoint: %s", apiEndpoint)
	}
}
