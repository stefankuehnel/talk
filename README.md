# talk

[![CI](../../actions/workflows/ci.yaml/badge.svg)](../../actions/workflows/ci.yaml)

A Command-Line Interface (CLI) in Go for sending messages to [Nextcloud Talk](https://nextcloud.com/talk/), built with [Cobra](https://github.com/spf13/cobra).

## Installation

### Using Nix

If you have Nix with flakes enabled:

```bash
nix run github:stefankuehnel/talk#talk -- <args>
```

Or install it to your profile:

```bash
nix profile install github:stefankuehnel/talk#talk
```

### Using Go

```bash
go install github.com/stefankuehnel/talk@latest
```

### Building from Source

```bash
git clone https://github.com/stefankuehnel/talk.git
cd talk
go build
```

## Usage

Send a message:

```bash
talk send \
  --server-url "https://<nextcloud_server_url>" \
  --chat-id "<nextcloud_talk_chat_id>" \
  --username "<nextcloud_username>" \
  --password "<nextcloud_app_password>" \
  --message "Deployment {{.Version}} finished" \
  --message-data '{"Version":"v1.2.3"}'
```

> [!NOTE]
> For convenience, the flags `--server-url`, `--chat-id`, `--username`, and `--password` can also be provided via
the environment variables `TALK_SERVER_URL`, `TALK_CHAT_ID`, `TALK_USERNAME`, and `TALK_PASSWORD`.

## Development

This project uses [Task](https://taskfile.dev) as a task runner.

### Available Tasks

```bash
# Run default tasks (lint, build and test)
task

# Build the project
task build

# Run Go tests
task test

# Format Go code
task fmt

# Run Go linter
task lint

# Clean build artifacts
task clean
```

## License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.
