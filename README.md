# FlexFrog Telegram Bot
<img src="img.png" width="256" alt=""/>

FlexFrog is a Telegram bot designed to interact with users via the Telegram API. It provides features such as sending messages, downloading files, and checking user permissions in a chat.

## Features

- Long-polling to fetch updates from Telegram.
- Send text messages and payloads to users.
- Download files from Telegram servers.
- Check if a user is an administrator in a chat.

## Project Structure

- `tg-bot-api/`: Contains the core logic for interacting with the Telegram API.
- `.github/workflows/build.yml`: GitHub Actions workflow for CI/CD, including building, testing, and Docker image deployment.

## Prerequisites

- Go 1.24 or later
- Docker (for building and pushing images)
- A Telegram bot token (set as an environment variable `BOT_TOKEN`)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/<your-username>/flexfrog-telegram-bot.git
   cd flexfrog-telegram-bot