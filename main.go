package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/NicoNex/echotron/v3"
	"github.com/tjarratt/babble"

	_ "embed"
)

type bot struct {
	id int64
	echotron.API
}

var (
	//go:embed token
	token    string
	words    = readWords()
	api      = echotron.NewAPI(token)
	babbler  = babble.NewBabbler()
	startMsg = "Hi, this is Jiitbot!\nGenerate new Jitsi meetings with the command /new or try the inline mode by typing @jiitbot in any chat!"

	commands = []echotron.BotCommand{
		{Command: "/new", Description: "Generate a new Jitsi meeting."},
	}
)

func newBot(id int64) echotron.Bot {
	return &bot{id, api}
}

func (b *bot) Update(update *echotron.Update) {
	if update.InlineQuery != nil {
		query := update.InlineQuery
		b.AnswerInlineQuery(query.ID, inline(meeting()), &echotron.InlineQueryOptions{CacheTime: 1})
		return
	}

	switch msg := message(update); {
	case strings.HasPrefix(msg, "/start"):
		b.SendMessage(startMsg, b.id, nil)

	case strings.HasPrefix(msg, "/new"):
		b.SendMessage(meeting(), b.id, nil)
	}
}

func inline(s string) []echotron.InlineQueryResult {
	return []echotron.InlineQueryResult{
		echotron.InlineQueryResultArticle{
			Type:        echotron.InlineArticle,
			ID:          fmt.Sprintf("%d", time.Now().Unix()),
			Title:       "Jitsi Meet",
			Description: s,
			InputMessageContent: echotron.InputTextMessageContent{
				MessageText: s,
			},
		},
	}
}

// Returns the message from the given update.
func message(u *echotron.Update) string {
	if u.Message != nil {
		return u.Message.Text
	} else if u.EditedMessage != nil {
		return u.EditedMessage.Text
	} else if u.CallbackQuery != nil {
		return u.CallbackQuery.Data
	}
	return ""
}

func readWords() []string {
	cnt, err := os.ReadFile("/usr/share/dict/words")
	if err != nil {
		return []string{"Error"}
	}
	return strings.Split(string(cnt), "\n")
}

func meeting() string {
	return fmt.Sprintf("https://meet.jit.si/%s", generate("", 4))
}

func generate(sep string, n int) string {
	toks := make([]string, n)
	for i := 0; i < n; i++ {
		toks[i] = strings.Title(words[rand.Intn(len(words))])
	}
	return strings.Join(toks, sep)
}

func main() {
	rand.Seed(time.Now().Unix())
	api.SetMyCommands(nil, commands...)
	dsp := echotron.NewDispatcher(token, newBot)
	log.Fatalln(dsp.Poll())
}
