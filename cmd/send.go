package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"stefanco.de/talk/internal/message"
	"stefanco.de/talk/internal/talk"
)

// NewSendCmd creates and returns the send command.
func NewSendCmd() *cobra.Command {
	var serverURL string
	var chatID string
	var username string
	var password string
	var messageTemplate string
	var messageData string
	var insecure bool

	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send a message to a Nextcloud Talk chat",
		RunE: func(cmd *cobra.Command, args []string) error {
			if serverURL == "" {
				serverURL = os.Getenv("TALK_SERVER_URL")
			}
			if serverURL == "" {
				return fmt.Errorf("server URL must be set via --server-url or TALK_SERVER_URL")
			}

			if chatID == "" {
				chatID = os.Getenv("TALK_CHAT_ID")
			}
			if chatID == "" {
				return fmt.Errorf("chat ID must be set via --chat-id or TALK_CHAT_ID")
			}

			if username == "" {
				username = os.Getenv("TALK_USERNAME")
			}
			if username == "" {
				return fmt.Errorf("username must be set via --username or TALK_USERNAME")
			}

			if password == "" {
				password = os.Getenv("TALK_PASSWORD")
			}
			if password == "" {
				return fmt.Errorf("password must be set via --password or TALK_PASSWORD")
			}

			messageDataUnmarshaled := map[string]any{}
			if messageData != "" {
				if err := json.Unmarshal([]byte(messageData), &messageDataUnmarshaled); err != nil {
					return fmt.Errorf("invalid message data: %w", err)
				}
			}

			renderedMessage, err := message.Render(messageTemplate, messageDataUnmarshaled)
			if err != nil {
				return fmt.Errorf("render message template: %w", err)
			}

			talkClient := talk.NewClient(
				serverURL,
				username,
				password,
				talk.WithInsecure(insecure),
			)

			if err := talkClient.SendMessage(cmd.Context(), chatID, renderedMessage); err != nil {
				return err
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), "message sent")
			return err
		},
	}

	sendCmd.Flags().StringVar(&serverURL, "server-url", "", "Nextcloud server URL")
	sendCmd.Flags().StringVar(&chatID, "chat-id", "", "Nextcloud Talk chat ID")
	sendCmd.Flags().StringVar(&username, "username", "", "Nextcloud username")
	sendCmd.Flags().StringVar(&password, "password", "", "Nextcloud app password")
	sendCmd.Flags().StringVar(&messageTemplate, "message", "", "Message text or Go text/template")
	sendCmd.Flags().StringVar(&messageData, "message-data", "", "JSON object used as message data")
	sendCmd.Flags().BoolVar(&insecure, "insecure", false, "Disable TLS certificate verification")

	_ = sendCmd.MarkFlagRequired("message")

	return sendCmd
}
