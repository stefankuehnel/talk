package talk

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var ChatEndpoint = "%s/ocs/v2.php/apps/spreed/api/v1/chat/%s"

// Types

// ClientOption configures a Client.
type ClientOption func(*Client)

// Client sends messages to Nextcloud Talk.
type Client struct {
	baseURL  string
	username string
	password string

	httpClient *http.Client

	insecure bool
}

// Constructors

// NewClient creates a new Nextcloud Talk client.
func NewClient(baseURL, username, password string, clientOptions ...ClientOption) *Client {
	client := &Client{
		baseURL:  strings.TrimRight(baseURL, "/"),
		username: username,
		password: password,

		httpClient: http.DefaultClient,
	}

	for _, clientOption := range clientOptions {
		clientOption(client)
	}

	if client.insecure {
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // optional CLI flag for self-signed certs

		clonedHTTPClient := *client.httpClient
		clonedHTTPClient.Transport = transport
		client.httpClient = &clonedHTTPClient
	}

	return client
}

// Setters

// WithHTTPClient sets the HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(client *Client) {
		if httpClient != nil {
			client.httpClient = httpClient
		}
	}
}

// WithInsecure disables TLS certificate verification.
func WithInsecure(insecure bool) ClientOption {
	return func(client *Client) {
		client.insecure = insecure
	}
}

// Methods

// SendMessage posts a message into a Nextcloud Talk chat.
func (client *Client) SendMessage(ctx context.Context, chatId, message string) (err error) {
	payload, err := json.Marshal(map[string]string{
		"token":   chatId,
		"message": message,
	})
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf(ChatEndpoint, client.baseURL, url.PathEscape(chatId))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.SetBasicAuth(client.username, client.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("OCS-APIRequest", "true")

	res, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, res.Body.Close())
	}()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("nextcloud talk request failed with status %d", res.StatusCode)
		}

		return fmt.Errorf("nextcloud talk request failed with status %d: %s", res.StatusCode, string(body))
	}

	return nil
}
