package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken string `mapstructure:"tgToken"`
	GptToken      string `mapstructure:"gptToken"`
	Preamble      string `mapstructure:"preamble"`
}

func LoadConfig(path string) (c Config, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)
	return
}

func sendChatGPT(apiKey, sendText string) string {
	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	var response string

	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt:      []string{sendText},
		MaxTokens:   gpt3.IntPtr(100),
		Temperature: gpt3.Float32Ptr(0),
	}, func(res *gpt3.CompletionResponse) {
		response += res.Choices[0].Text
	})
	if err != nil {
		log.Println(err)
		return "ChatGPT is not available"
	}
	return response
}

func main() {
	var userPrompt string
	var gptPrompt string
	// Reading config.yaml
	config, err := LoadConfig(".")
	apiKey := config.GptToken

	if err != nil {
		panic(fmt.Errorf("fatal error with config.yaml: %w", err))
	}

	// Telegram initialization
	bot, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true // set to false for suppress logs in stdout
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Start Telegram long polling update
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	//Check message in updates
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// check if message has '/topic' or '/phrase' as prefix
		if !strings.HasPrefix(update.Message.Text, "/topic") && !strings.HasPrefix(update.Message.Text, "/phrase") {
			continue
		}

		// remove '/topic' or '/phrase' from message
		if strings.HasPrefix(update.Message.Text, "/topic") {
			userPrompt = strings.TrimPrefix(update.Message.Text, "/topic")
			gptPrompt = config.Preamble + "TOPIC: "
		} else if strings.HasPrefix(update.Message.Text, "/phrase") {
			userPrompt = strings.TrimPrefix(update.Message.Text, "/phrase")
			gptPrompt = config.Preamble + "PHRASE: "
		}

		if userPrompt != "" {
			gptPrompt += userPrompt
			// Send request to ChatGPT
			res := sendChatGPT(apiKey, gptPrompt)
			update.Message.Text = res
		} else {
			update.Message.Text = "Please, enter your topic or phrase"
		}

		// Send message to Telegram
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		_, err = bot.Send(msg)
		if err != nil {
			log.Println("Error:", err)
		}
	}
}
