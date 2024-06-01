
# Slack AI Bot

OpenAI-enabled Slack chat bot for demo purposes.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## Features

- Integrates with OpenAI to provide intelligent responses within Slack.
- Configurable and extendable.
- Easy to deploy using Docker Compose.
- Responds to bot mentions in a thread
- Tracks each thread separately

## Installation
1. Install a slack app by following [slack documentation](https://api.slack.com/quickstart#creating).
2. Generate `config.toml`, you can take [config.toml.example](config.toml.example) as refference.
3. Start the app by running `air` 

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.16 or higher)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [github.com/cosmtrek/air](https://github.com/cosmtrek/air) live reloading go development tool
- Slack workspace and app. You can use 

### Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/slack-ai-bot.git
   cd slack-ai-bot
   ```

2. Build the Docker image:

   ```bash
   docker-compose build
   ```

3. Run the Docker container:

   ```bash
   docker-compose up
   ```

## Usage

Once the bot is running, you can interact with it within your Slack workspace. Mention the bot in any channel or send it a direct message to receive responses powered by OpenAI.

## Configuration

The bot can be configured using environment variables or configuration files. The main configuration file is `config.go` located in the `internal/config` directory.

### Environment Variables

- `SLACK_BOT_TOKEN`: Your Slack bot token.
- `OPENAI_API_KEY`: Your OpenAI API key.

### Example `.env` file

```
SLACK_BOT_TOKEN=xoxb-your-slack-bot-token
OPENAI_API_KEY=your-openai-api-key
```

## Contributing

Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Create a new Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

If you have any questions, feel free to open an issue or reach out to the project maintainer at [your-email@example.com](mailto:your-email@example.com).
