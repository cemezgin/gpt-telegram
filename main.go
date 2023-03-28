package main

import (
	"os/signal"
	"context"
	"os"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"act-as-chat/prompts/awesomechatgptprompts"
	"github.com/sashabaranov/go-openai"
	"fmt"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	b, _ := bot.New(os.Getenv("TELEGRAM_BOT_KEY"), opts...)
	b.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	//todo fetch via cron and store in memory
	prompts := awesomechatgptprompts.Exec()

	if update.Message.Text == "/start" {

		kb := inline.New(b)
		row := kb.Row()
		i := 0
		for name, prompt := range prompts {
			if i == 2 {
				row = kb.Row()
				i = 0
			}
			row.Button(name, []byte(prompt), onInlineKeyboardSelect)
			i++
		}

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "Select the prompt you want to use:",
			ReplyMarkup: kb,
		})

		if err != nil {
			fmt.Printf("chat completion error: %v\n", err)
			return
		}
	} else {
		gpt(ctx, b, update.Message)
	}
}

func onInlineKeyboardSelect(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   string(data),
	})

	gpt(ctx, b, msg)

	fmt.Printf("(Request) Client: %s\n", msg.Text)

	if err != nil {
		return
	}

}

func gpt(ctx context.Context, b *bot.Bot, msg *models.Message) *models.Message {
	client := openai.NewClient(os.Getenv("OPEN_AI_KEY"))

	//followup conversation
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg.Text,
				},
			},
		},
	)

	if err != nil {
		return nil
	}

	msg, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: msg.Chat.ID,
		Text:   "(Chat GPT) : " + resp.Choices[0].Message.Content,
	})

	fmt.Printf("(Response) Chat GPT: %s\n", resp.Choices[0].Message.Content)

	return msg
}
