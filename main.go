package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// configuration structure

type Config struct {
	TelegramBotToken  string `json:"telegram_bot_token"`
	WalletLookupToken string `json:"wallet_lookup_token"`
	WalletLookupURL   string `json:"wallet_lookup_url"`
}

func main() {
	color.Cyan("Initializing Bot..")
	config, err := loadConfig()
	if err != nil {
		color.Red("Failed to load configuration, error: %s", err)
		return
	}
	bot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		color.Red("Failed to initialize bot, error: %s", err)
		return
	}
	color.Green("Bot successfully initialized as %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			var username_author string
			if len(update.Message.From.UserName) > 0 {
				username_author = "@" + update.Message.From.UserName
			} else {
				username_author = ""
			}

			switch update.Message.Command() {
			case "start":
				msg.Text = fmt.Sprintf("*Welcome %s %s* 👋\n\nThis bot was created to look up any crypto wallet by address.\nThe bot will automatically detect the coin type by providing a wallet address to check.\n\n👨‍💻 Coder: @Zile42O\n💾 [Source](https://github.com/zile42O/wallet-lookup-tgbot)\n\n✅ *Supported coins:* Bitcoin, Bitcoin Cash, Litecoin, Dogecoin, Dash\n\n*Use command:* /lookup", update.Message.From.FirstName, update.Message.From.LastName)
			case "lookup":
				arg := strings.Fields(update.Message.Text)
				if len(arg) < 2 {
					msg.Text = "❌ *Please provide wallet address.*"
				} else {
					msg.Text = fmt.Sprintf("⏳ *Please wait..*\n\n🔎 Looking for: `%s`", arg[1])
					msg.ParseMode = "Markdown"
					sent_message, err := bot.Send(msg)
					if err == nil {
						msg.Text = lookupAddress(bot, config, arg[1])
						msg.ReplyToMessageID = sent_message.MessageID
					}
				}
			default:
				continue
			}

			msg.ParseMode = "Markdown"
			_, err := bot.Send(msg)
			if err != nil {
				color.Red("Bot failed to reply on command: %s to user: %d (%s)", update.Message.Command(), update.Message.From.ID, username_author)
			} else {
				color.Blue("Bot success to reply on command: %s to user: %d (%s)", update.Message.Command(), update.Message.From.ID, username_author)
			}
		}
	}

}

// functions

func lookupAddress(bot *tgbotapi.BotAPI, config Config, address string) string {
	var return_string string
	query := fmt.Sprintf("%s?key=%s&address=%s", config.WalletLookupURL, config.WalletLookupToken, address)
	response, err := http.Get(query)
	if err != nil {
		return "❌ *API is currently unavailable*"
	}
	defer response.Body.Close()

	var data map[string]interface{}
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&data)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return_string = "❌ *Error while processing request*"
	} else {
		if errorObj, ok := data["error"]; ok {
			errorMessage, ok := errorObj.(string)
			if ok {
				return_string = "❌ *" + errorMessage + "*"
			} else {
				return_string = "❌ *Unexpected error format*"
			}
		} else if data_count, ok := data["data"].(map[string]interface{}); ok {
			if len(data_count) > 0 {
				if ok {
					dataObj, ok := data["data"].(map[string]interface{})
					if ok {
						return_string = fmt.Sprintf("🔎 *Data for* `%s`:\n\n🪙 *Coin:* %s\n\n💹 *Coin Market Price*: $%.2f\n\n🧾 *Transactions:* %d (Sent: %d Received: %d)\n\n💰 *Balance:* (%.8f) $%.2f\n\n🤑 *Recieved:* $%.2f\n\n💸 *Spent:* $%.2f\n\n📅 Dates:\n\n*First Time Transaction:* %s\n*Last Time Transaction:* %s\n",
							address,
							dataObj["coin"].(string),
							dataObj["market_price_usd"].(float64),
							int(dataObj["total_transactions"].(float64)),
							int(dataObj["transactions_sent"].(float64)),
							int(dataObj["transactions_received"].(float64)),
							dataObj["balance"].(float64),
							dataObj["balance_usd"].(float64),
							dataObj["receieved_value_usd"].(float64),
							dataObj["sent_value_usd"].(float64),
							dataObj["first_transaction"].(string),
							dataObj["last_transaction"].(string))
					} else {
						return_string = "❌ *Cannot access API data, try again later*"
					}
				} else {
					return_string = "❌ *No API data currently, try again later*"
				}
			}
		} else {
			return_string = "❌ *Unexpected response format*"
		}
	}
	return return_string
}

func loadConfig() (Config, error) {
	var config Config
	file, err := os.Open("config.json")
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
