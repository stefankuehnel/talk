package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"testing"

	"net/http"
	"net/http/httptest"

	"stefanco.de/talk/internal/talk"
)

func TestSendCmd(t *testing.T) {
	t.Run("sends rendered message", func(t *testing.T) {
		t.Parallel()

		chatID := "chat-id"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			baseURL := "http://" + r.Host
			expectedURL, err := url.Parse(fmt.Sprintf(talk.ChatEndpoint, baseURL, url.PathEscape(chatID)))
			if err != nil {
				t.Fatalf("parse expected URL: %v", err)
			}
			expectedPath := expectedURL.Path
			if r.URL.Path != expectedPath {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}

			var payload map[string]string
			if err := json.Unmarshal(body, &payload); err != nil {
				t.Fatalf("unmarshal body: %v", err)
			}
			if payload["message"] != "Hello Stefan" {
				t.Fatalf("expected rendered message %q, got %q", "Hello Stefan", payload["message"])
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		out := new(bytes.Buffer)

		rootCmd := NewRootCmd()
		rootCmd.SetOut(out)
		rootCmd.SetErr(out)
		rootCmd.SetArgs([]string{
			"send",
			"--server-url", server.URL,
			"--chat-id", chatID,
			"--username", "user",
			"--password", "pass",
			"--message", "Hello {{.Name}}",
			"--message-data", `{"Name":"Stefan"}`,
		})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if out.String() != "message sent\n" {
			t.Fatalf("expected output %q, got %q", "message sent\n", out.String())
		}
	})

	t.Run("returns error for invalid template data", func(t *testing.T) {
		t.Parallel()

		rootCmd := NewRootCmd()
		rootCmd.SetArgs([]string{
			"send",
			"--server-url", "https://nextcloud.example.com",
			"--chat-id", "chat-id",
			"--username", "user",
			"--password", "pass",
			"--message", "Hello",
			"--message-data", "{not-json}",
		})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns error for missing template data keys", func(t *testing.T) {
		t.Parallel()

		requestCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		rootCmd := NewRootCmd()
		rootCmd.SetArgs([]string{
			"send",
			"--server-url", server.URL,
			"--chat-id", "chat-id",
			"--username", "user",
			"--password", "pass",
			"--message", "Hello {{.Name}}",
			"--message-data", "{}",
		})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if requestCount != 0 {
			t.Fatalf("expected no request to be sent, got %d", requestCount)
		}
	})

	t.Run("uses TALK_* environment variables when flags are omitted", func(t *testing.T) {
		chatID := "chat-id"
		t.Setenv("TALK_PASSWORD", "pass")

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			baseURL := "http://" + r.Host
			expectedURL, err := url.Parse(fmt.Sprintf(talk.ChatEndpoint, baseURL, url.PathEscape(chatID)))
			if err != nil {
				t.Fatalf("parse expected URL: %v", err)
			}
			expectedPath := expectedURL.Path
			if r.URL.Path != expectedPath {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		t.Setenv("TALK_SERVER_URL", server.URL)
		t.Setenv("TALK_CHAT_ID", chatID)
		t.Setenv("TALK_USERNAME", "user")

		rootCmd := NewRootCmd()
		rootCmd.SetArgs([]string{
			"send",
			"--message", "Hello",
		})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}
