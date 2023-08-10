// main telegram bot which asks for user inputs and outputs the recommended result
// communicates with recommendationservice through gRPC

package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BotState int

type UserInput struct {
	Mealtime         string   `protobuf:"bytes,1,opt,name=mealtime,proto3" json:"mealtime,omitempty"`
	Mood             string   `protobuf:"bytes,2,opt,name=mood,proto3" json:"mood,omitempty"`
	CuisineDontWants []string `protobuf:"bytes,3,rep,name=cuisine_dont_wants,json=cuisineDontWants,proto3" json:"cuisine_dont_wants,omitempty"`
	Budget           int32    `protobuf:"varint,4,opt,name=budget,proto3" json:"budget,omitempty"`
}

const (
	StateStart     BotState = iota
	StateQuestion1          //mealtime
	StateQuestion2          //budget
	StateQuestion3          //mood
	StateQuestion4          //text cuisine
	StateQuestion5          //personalized
	StateEnd
)

const (
	ADDRESS = "10.0.53.179:83"
)

// caching of relevant info
var (
	userState   = make(map[int64]BotState) //map of key: userID; val: BotState
	userChoices = make(map[int64][]string) //map of key: userID; val: uerChoices[] string array
	userFinalID = make(map[int64]int)      //map of key: userID; val: finalMessageID int
)

func buttonCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	query := update.CallbackQuery
	data := query.Data
	chatID := query.Message.Chat.ID
	userID := query.From.ID
	messageID := query.Message.MessageID

	//logic: gets a data feedback from user, appends it to the choice cache, moves to next state and renders next qn
	switch userState[userID] {
	case StateQuestion1:
		userChoices[userID] = append(userChoices[userID], data)
		userState[userID] = StateQuestion2
		sendQuestion(bot, chatID, messageID, "Whats your average budget per item?",
			[]string{"below $10", "below $20", "below $30", "below $40", "below $50", "Over $50"})

	case StateQuestion2:
		userChoices[userID] = append(userChoices[userID], data)
		userState[userID] = StateQuestion3
		sendQuestion(bot, chatID, messageID, "Whats your mood?",
			[]string{"Comfort", "Casual", "Healthy", "Indulgent", "Fancy"})

	case StateQuestion3:
		userChoices[userID] = append(userChoices[userID], data)
		userState[userID] = StateQuestion4
		sendTextQuestion(bot, chatID, messageID,
			"What cuisines would you not want to eat?\n"+
				"*Reply* to this message with a list of the cuisines with a space between each type",
		)

	// case StateQuestion4:
	// 	userState[userID] = StateQuestion5
	// 	finalMessageID := sendFinalResponse(bot, chatID, userID, userChoices[userID]) //this is state5
	// 	userFinalID[userID] = finalMessageID

	case StateQuestion5:
		userChoices[userID] = append(userChoices[userID], data)
		userState[userID] = StateEnd
		sendPersonalized(bot, chatID, userID, userFinalID[userID], userChoices[userID])

	default:
		sendWelcomeMessage(bot, chatID, userID)
	}
}

// Send a question with inline keyboard choices, edit the current message
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

func startQuestion(bot *tgbotapi.BotAPI, chatID int64, question string, choices []string) int {
	var options [][]tgbotapi.InlineKeyboardButton
	for _, choice := range choices {
		options = append(options, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(choice, choice),
		})
	}

	menu := tgbotapi.NewInlineKeyboardMarkup(options...)
	msg := tgbotapi.NewMessage(chatID, question)
	msg.ReplyMarkup = &menu

	sentMessage, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}

	return sentMessage.MessageID
}

// TODO: add new request method to server
func sendPersonalized(bot *tgbotapi.BotAPI, chatID int64, userID int64, finalMessageID int, userIDChoices []string) {
	//delete the previous message
	del := tgbotapi.NewDeleteMessage(chatID, finalMessageID)
	_, delErr := bot.Send(del)
	if delErr != nil {
		log.Println(delErr)
	}

	selected := userIDChoices[len(userIDChoices)-1]

	//send request to gRPC with selected choice
	//response is an array of objects
	res := generatePersonalised(selected)

	//display as a message
	message := "Based on what you chose previously, "
	message += selected + "\nThese are recommended: \n\n"
	//loop thru response and destructure into message string

	for _, r := range res {
		message += r.Name + "\n"
		message += "Rated " + strconv.Itoa(int(r.Rating)) + "/5 \n"
		message += "Average budget per pax: " + strconv.Itoa(int(r.Budget)) + "\n"
		message += "Cuisine: " + strings.Join(r.Cuisine, ", ")
		message += "\n\n"
	}

	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}

	// Reset the conversation for the user
	delete(userState, userID)
	delete(userChoices, userID)
	delete(userFinalID, userID)
}

func parseBudget(budgetStr string) int32 {
	switch budgetStr {
	case "below $10":
		return 10
	case "below $20":
		return 20
	case "below $30":
		return 30
	case "below $40":
		return 40
	case "below $50":
		return 50
	case "Over $50":
		return 100
	default:
		return 100
	}
}

func parseMealtime(mealtimeStr string) string {
	switch mealtimeStr {
	case "Lunch":
		return "lunch"
	case "Brunch":
		return "brunch"
	case "Dinner":
		return "dinner"
	case "Snack/Dessert":
		return "snack"
	default:
		return "lunch"
	}
}

// Send the final response with the collected choices
func sendFinalResponse(bot *tgbotapi.BotAPI, chatID int64, userID int64, userIDChoices []string) int {
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

	u := UserInput{
		Mealtime:         parseMealtime(userIDChoices[0]),
		Budget:           parseBudget(userIDChoices[1]),
		Mood:             userIDChoices[2],
		CuisineDontWants: strings.Fields(userIDChoices[3]),
	}

	recommendation := generateRecommendation(u)
	var options []string

	for _, r := range recommendation {
		response += r.Name + "\n"
		response += "Rated " + strconv.Itoa(int(r.Rating)) + "/5 \n"
		response += "Average budget per pax: " + strconv.Itoa(int(r.Budget)) + "\n"
		response += "Cuisine: " + strings.Join(r.Cuisine, ", ")
		response += "\n\n"

		options = append(options, r.Name)
	}
	response += "Not satisfied? Choose one of the above for a more personalized recommendation! \n"
	// response += generatePrompt(userIDChoices)
	messageID := startQuestion(bot, chatID, response, options)

	return messageID

}

func sendWelcomeMessage(bot *tgbotapi.BotAPI, chatID int64, userID int64) {
	userState[userID] = StateQuestion1
	userChoices[userID] = nil

	startQuestion(bot, chatID, "Which meal are you looking recommendations for!",
		[]string{"Brunch", "Lunch", "Dinner", "Snack/Dessert"})
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

func generateRecommendation(userChoices UserInput) []*Place {
	conn, err := grpc.Dial(ADDRESS, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic("Error connecting to GRPC server")
	}

	defer conn.Close()

	client := NewMessageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	response, err := client.GetRecommendation(ctx, &RecommendationRequest{
		Mealtime:         userChoices.Mealtime,
		Mood:             userChoices.Mood,
		CuisineDontWants: userChoices.CuisineDontWants,
		Budget:           userChoices.Budget,
	})
	// response, err := client.Echo(context.Background(), &EchoRequest{Message: "hello"})
	if err != nil {
		log.Print("gRPC Request Error: ")
		log.Println(err)
	}

	return response.Recommendations
}

func generatePersonalised(userChoice string) []*Place {
	conn, err := grpc.Dial(ADDRESS, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic("Error connecting to GRPC server")
	}

	defer conn.Close()

	client := NewMessageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	response, err := client.GetPersonalised(ctx, &PersonalRequest{
		Name: userChoice,
	})
	// response, err := client.Echo(context.Background(), &EchoRequest{Message: "hello"})
	if err != nil {
		log.Print("gRPC Request Error: ")
		log.Println(err)
	}

	return response.Recommendations
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
				userState[userID] = StateQuestion5
				userFinalID[userID] = sendFinalResponse(bot, update.Message.Chat.ID, userID, userChoices[userID])
			}
		}
	}
}
