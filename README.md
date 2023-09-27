# ELI5 Telegram Bot using OpenAI API

## Install

Just a few simple steps:

1. Clone repository
2. Rename file `config.yaml.example` to `config.yaml` and set up your tokens (OpenAI and Telegram) in file
3. Run `go mod tidy`
4. Run `go run main.go`

## Usage

This is a simple Telegram bot that uses the OpenAI API to generate a ELI5 explanation of a given topic or phrase.

### Commands

- `/start` - Start the bot
- `/topic <topic>` - Generate a ELI5 explanation of a given topic
- `/phrase <phrase>` - Generate a ELI5 explanation of a given phrase
