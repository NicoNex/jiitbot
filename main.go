package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/NicoNex/echotron/v3"

	_ "embed"
)

type bot struct {
	id int64
	echotron.API
}

var (
	//go:embed token
	token string
	//go:embed words.txt
	wordFile string
	words    = strings.Split(wordFile, "\n")
	api      = echotron.NewAPI(token)
	startMsg = "Hi, this is Jiitbot!\nGenerate new Jitsi meetings with the command /new or try the inline mode by typing @jiitbot in any chat!"

	commands = []echotron.BotCommand{
		{Command: "/new", Description: "Generate a new Jitsi meeting."},
	}

	iqOpts = &echotron.InlineQueryOptions{CacheTime: 1}
)

func newBot(id int64) echotron.Bot {
	return &bot{id, api}
}

func (b *bot) handleInlineQuery(iq *echotron.InlineQuery) {
	if iq.Query != "" {
		check(b.AnswerInlineQuery(iq.ID, inline(meeting(iq.Query)), iqOpts))
		return
	}
	check(b.AnswerInlineQuery(iq.ID, inline(meeting(generate("", 4))), iqOpts))
}

func (b *bot) Update(update *echotron.Update) {
	if update.InlineQuery != nil {
		b.handleInlineQuery(update.InlineQuery)
		return
	}

	switch msg := message(update); {
	case strings.HasPrefix(msg, "/start"):
		check(b.SendMessage(startMsg, b.id, nil))

	case strings.HasPrefix(msg, "/new"):
		check(b.SendMessage(meeting(generate("", 4)), b.id, nil))
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

func check[T any](first T, err error, a ...any) T {
	if err != nil {
		log.Println(append([]any{err}, a...)...)
	}
	return first
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

func readOSDictWords() []string {
	cnt, err := os.ReadFile("/usr/share/dict/words")
	if err != nil {
		return []string{"Error"}
	}
	return strings.Split(string(cnt), "\n")
}

func meeting(s string) string {
	return fmt.Sprintf("https://meet.jit.si/%s", s)
}

func generate(sep string, n int) string {
	var toks = make([]string, n)
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
