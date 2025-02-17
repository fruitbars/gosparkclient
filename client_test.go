package gosparkclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewSparkClient(t *testing.T) {
	tests := []struct {
		name        string
		opts        []ConfigOption
		wantErr     bool
		errContains string
	}{
		{
			name: "valid configuration",
			opts: []ConfigOption{
				WithCredentials("appid", "apikey", "secret"),
				WithURLs("http://example.com", "http://emb.example.com"),
				WithDomain("domain"),
			},
			wantErr: false,
		},
		{
			name: "missing credentials",
			opts: []ConfigOption{
				WithURLs("http://example.com", "http://emb.example.com"),
			},
			wantErr:     true,
			errContains: "AppID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewSparkClient(tt.opts...)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if client == nil {
				t.Error("expected client, got nil")
			}
		})
	}
}

func TestSparkClient_ChatSimple(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
            "header": {
                "code": 0,
                "message": "success",
                "sid": "test-sid",
                "status": 2
            },
            "payload": {
                "choices": {
                    "status": 2,
                    "seq": 1,
                    "text": [
                        {
                            "content": "test response",
                            "role": "assistant",
                            "index": 0
                        }
                    ]
                },
                "usage": {
                    "text": {
                        "question_tokens": 10,
                        "prompt_tokens": 10,
                        "completion_tokens": 20,
                        "total_tokens": 30
                    }
                }
            }
        }`))
	}))
	defer mockServer.Close()

	client, err := NewSparkClient(
		WithCredentials("test-app-id", "test-api-key", "test-secret"),
		WithURLs(mockServer.URL, mockServer.URL),
		WithTimeout(time.Second),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	resp, err := client.ChatSimple(ctx, "Hello")
	if err != nil {
		t.Fatalf("ChatSimple failed: %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}
}

func TestSparkClient_WithNewConfig(t *testing.T) {
	client, err := NewSparkClient(
		WithCredentials("test-app-id", "test-api-key", "test-secret"),
		WithURLs("http://example.com", "http://emb.example.com"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test creating new client with updated config
	newClient, err := client.WithNewConfig(
		WithTimeout(time.Second*60),
		WithUID("new-uid"),
	)
	if err != nil {
		t.Errorf("WithNewConfig failed: %v", err)
	}

	if newClient.config.Timeout != time.Second*60 {
		t.Errorf("expected timeout 60s, got %v", newClient.config.Timeout)
	}
	if newClient.config.UID != "new-uid" {
		t.Errorf("expected UID 'new-uid', got %v", newClient.config.UID)
	}

	// Verify original client remains unchanged
	if client.config.Timeout == time.Second*60 {
		t.Error("original client timeout should not have changed")
	}
	if client.config.UID == "new-uid" {
		t.Error("original client UID should not have changed")
	}
}
