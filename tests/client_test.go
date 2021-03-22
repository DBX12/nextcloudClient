package tests

import (
	"nextcloudClient/nextcloudClient"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := nextcloudClient.NewClient("http://example.local", "the-user", "the-secret-password")
	if c == nil {
		t.Fatal("Client was not created")
	}
	if c.HTTPClient == nil {
		t.Fatal("http.Client was not created")
	}
}
