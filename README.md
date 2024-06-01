
# Slack AI Bot

OpenAI-enabled Slack chat bot for demo purposes.

## Table of Contents

- [Features](#features)
- [Prerequisites](#Prerequisites)
- [Steps](#steps)
- [Usage](#usage)
- [Configuration](#configuration)
- [License](#license)

## Features

- Integrates with OpenAI to provide intelligent responses within Slack.
- Configurable and extendable.
- Easy to deploy using Docker Compose.
- Responds to bot mentions in a thread
- Tracks each thread separately


### Prerequisites

- [Go](https://golang.org/dl/) (version 1.16 or higher)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [github.com/cosmtrek/air](https://github.com/cosmtrek/air) live reloading go development tool
- Slack workspace and app. You can use [slack-app-example.yaml](slack-app-example.yaml) as your slack app definition.


### Steps
1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/slack-ai-bot.git
   cd slack-ai-bot
   ```

1. Install a slack app by following [slack documentation](https://api.slack.com/quickstart#creating).
1. Generate `config.toml`, you can take [config.toml.example](config.toml.example) as refference.
1. Start the app by running `air` 

## Usage

Once the bot is running, you can interact with it within your Slack workspace. Mention the bot in any channel or send it a direct message to receive responses powered by OpenAI.

## Configuration

The bot can be configured using environment variables or configuration files. The main configuration file is `config.go` located in the `internal/config` directory.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
