package talk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestClientSendMessage(t *testing.T) {
	t.Parallel()

	t.Run("posts message to nextcloud talk endpoint", func(t *testing.T) {
		t.Parallel()

		chatID := "chat-id"
		message := "hello world"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("expected method %q, got %q", http.MethodPost, r.Method)
			}

			expectedPath := fmt.Sprintf(ChatEndpoint, "http://"+r.Host, url.PathEscape(chatID))
			expectedPath = expectedPath[len("http://"+r.Host):]
			if r.URL.Path != expectedPath {
				t.Fatalf("expected path %q, got %q", expectedPath, r.URL.Path)
			}

			username, password, ok := r.BasicAuth()
			if !ok {
				t.Fatal("expected basic auth to be set")
			}

			if username != "username" {
				t.Fatalf("expected username %q, got %q", "username", username)
			}

			if password != "password" {
				t.Fatalf("expected password %q, got %q", "password", password)
			}

			if got := r.Header.Get("Content-Type"); got != "application/json" {
				t.Fatalf("expected content-type application/json, got %q", got)
			}

			if got := r.Header.Get("Accept"); got != "application/json" {
				t.Fatalf("expected accept application/json, got %q", got)
			}

			if got := r.Header.Get("OCS-APIRequest"); got != "true" {
				t.Fatalf("expected OCS-APIRequest true, got %q", got)
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}

			var payload map[string]string
			if err := json.Unmarshal(body, &payload); err != nil {
				t.Fatalf("unmarshal body: %v", err)
			}

			if payload["token"] != chatID {
				t.Fatalf("expected token %q, got %q", chatID, payload["token"])
			}

			if payload["message"] != message {
				t.Fatalf("expected message %q, got %q", message, payload["message"])
			}

			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		taskClient := NewClient(server.URL, "username", "password")
		if err := taskClient.SendMessage(context.Background(), chatID, message); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("returns an error for non-2xx responses", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("unauthorized"))
		}))
		defer server.Close()

		taskClient := NewClient(server.URL, "username", "password")
		err := taskClient.SendMessage(context.Background(), "chat-id", "hello")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
