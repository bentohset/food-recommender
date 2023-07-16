package main

import (
	"context"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type BotState int

const (
	StateStart BotState = iota
	StateQuestion1
	StateQuestion2
	StateQuestion3
	StateQuestion4
	StateEnd
)

var (
	userState   = make(map[int64]BotState) //map of key: userID; val: BotState
	userChoices = make(map[int64][]string) //map of key: userID; val: uerChoices[] string array
)

func buttonCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	query := update.CallbackQuery
	data := query.Data
	chatID := query.Message.Chat.ID
	userID := query.From.ID
	messageID := query.Message.MessageID

	switch userState[userID] {
	case StateQuestion1:
		userChoices[userID] = append(userChoices[userID], data)
		userState[userID] = StateQuestion2
		sendQuestion(bot, chatID, messageID, "Whats your average budget per item?",
			[]string{"below $10", "$10 - $20", "$20 - $30", "$30 - $50", "over $50"})

	case StateQuestion2:
		userChoices[userID] = append(userChoices[userID], data)
		userState[userID] = StateQuestion3
		sendQuestion(bot, chatID, messageID, "Whats your mood?",
			[]string{"Comfort", "Energy", "Healthy", "Indulgent", "Buffet"})

	case StateQuestion3:
		userChoices[userID] = append(userChoices[userID], data)
		userState[userID] = StateQuestion4
		sendTextQuestion(bot, chatID, messageID, "What cuisines would you not want to eat? *Reply* to this message with a list of the cuisines with a space between each type")

	case StateQuestion4:
		userState[userID] = StateEnd
		sendFinalResponse(bot, chatID, userID, userChoices[userID])

	default:
		sendWelcomeMessage(bot, chatID, userID)
	}
}

// Send a question with inline keyboard choices
func sendQuestion(bot *tgbotapi.BotAPI, chatID int64, messageID int, question string, choices []string) {
	var options [][]tgbotapi.InlineKeyboardButton
	for _, choice := range choices {
		options = append(options, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(choice, choice),
		})
	}

	menu := tgbotapi.NewInlineKeyboardMarkup(options...)
	msg := tgbotapi.NewEditMessageText(chatID, messageID, question)
	msg.ReplyMarkup = &menu

	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// delete the questions message and send new message asking for user reply
func sendTextQuestion(bot *tgbotapi.BotAPI, chatID int64, messageID int, question string) {
	del := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, delErr := bot.Send(del)
	if delErr != nil {
		log.Println(delErr)
	}

	msg := tgbotapi.NewMessage(chatID, question)
	msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
	msg.ParseMode = "MarkdownV2"

	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func startQuestion(bot *tgbotapi.BotAPI, chatID int64, question string, choices []string) {
	var options [][]tgbotapi.InlineKeyboardButton
	for _, choice := range choices {
		options = append(options, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(choice, choice),
		})
	}

	menu := tgbotapi.NewInlineKeyboardMarkup(options...)
	msg := tgbotapi.NewMessage(chatID, question)
	msg.ReplyMarkup = &menu

	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// Send the final response with the collected choices
func sendFinalResponse(bot *tgbotapi.BotAPI, chatID int64, userID int64, userIDChoices []string) {
	response := "Based on your preferences below,\n\n"
	titles := [4]string{"Mealtime", "Budget", "Mood", "Cuisine don't-wants"}
	for i, choice := range userIDChoices {
		response += titles[i] + ": " + choice + "\n"
	}

	response += "\nI suggest: \n\n"

	//add openai generated answer to the response
	// resChannel := make(chan string)

	// go func() {
	// 	res := generatePrompt(userIDChoices)
	// 	resChannel <- res
	// }()

	// response += <-resChannel

	response += generatePrompt(userIDChoices)

	msg := tgbotapi.NewMessage(chatID, response)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}

	// Reset the conversation for the user
	delete(userState, userID)
	delete(userChoices, userID)
}

func sendWelcomeMessage(bot *tgbotapi.BotAPI, chatID int64, userID int64) {
	userState[userID] = StateQuestion1
	userChoices[userID] = nil

	startQuestion(bot, chatID, "Which meal are you looking recommendations for!",
		[]string{"Brunch", "Lunch", "Dinner", "Just a snack"})
}

func generatePrompt(prompts []string) string {
	str := prompts[3]
	items := strings.Fields(str)
	cuisines := strings.Join(items, ", ")

	text := "Suggest 3 food places in Singapore with the following parameters: \n" +
		"Mealtime - " + prompts[0] + "\n" +
		"Budget - " + prompts[1] + "\n" +
		"Mood - " + prompts[2] + " foods\n" +
		"Excluding the following cuisines - " + cuisines

	client := openai.NewClient(os.Getenv("OPENAI_APITOKEN"))
	log.Println(text)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: text,
				},
			},
		},
	)

	if err != nil {
		log.Print("OpenAI Error:")
		log.Println(err)
		return err.Error()
	}
	log.Println(resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content
}

func main() {
	loadErr := godotenv.Load(".env")
	if loadErr != nil {
		log.Panic("Could not load .env")
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Print("Error loading bot")
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	// updates := bot.ListenForWebhook("/" + bot.Token)
	// go http.ListenAndServe("0.0.0.0:8080", nil)

	for update := range updates {
		if update.CallbackQuery != nil {
			go buttonCallback(bot, update)
		} else if update.Message != nil && update.Message.IsCommand() && update.Message.Command() == "start" {
			sendWelcomeMessage(bot, update.Message.Chat.ID, update.Message.From.ID)
		} else if update.Message != nil && isMention(update.Message, bot.Self.UserName) {
			log.Printf("Bot mentioned in message: %s", update.Message.Text)
			sendWelcomeMessage(bot, update.Message.Chat.ID, update.Message.From.ID)
		} else if update.Message != nil && update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.ID == bot.Self.ID {
			userID := update.Message.From.ID
			if userState[userID] == StateQuestion4 {
				userChoices[userID] = append(userChoices[userID], update.Message.Text)
				sendFinalResponse(bot, update.Message.Chat.ID, userID, userChoices[userID])
			}
		}
	}
}

func isMention(message *tgbotapi.Message, botUsername string) bool {
	if message.Entities != nil {
		for _, entity := range message.Entities {
			if entity.Type == "mention" && strings.HasPrefix(message.Text, "@"+botUsername) {
				return true
			}
		}
	}
	return false
}
