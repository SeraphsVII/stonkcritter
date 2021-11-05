package politstonk

import (
	"strings"

	"github.com/timshannon/badgerhold/v4"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (bot *Bot) follow(msg *tb.Message) {
	topic := msg.Payload

	if len(topic) > 6 {
		bot.Send(msg.Chat, "that looks like a representatives name, let me see if I can find them in my list...")
		names, err := bot.searchReps(topic)

		if err != nil {
			bot.Send(msg.Chat, "sorry, failed to search: "+err.Error())
			return
		}

		if len(names) > 10 {
			bot.Send(msg.Chat, "sorry, I found more than 10 results with the same name, you'll need to be more specific, try /findrep")
			return
		}

		if len(names) > 1 {
			bot.Send(msg.Chat, "OK, I found 10 reps that it could be, can you be more specific?:\n"+strings.Join(names, "\n"))
			return
		}

		bot.Send(msg.Chat, "Cool, found "+names[0]+" in my list, using that")
		topic = names[0]
	} else {
		bot.Send(msg.Chat, "OK that looks like a stock ticker")
	}

	s := Sub{ChatID: int32(msg.Chat.ID), Topic: topic}
	err := bot.store.Insert(s.String(), s)
	if err != nil {
		bot.Send(msg.Chat, "sorry, failed to save that: "+err.Error())
		return
	}

	bot.Send(msg.Chat, "OK saved that")
}

func (bot *Bot) unfollow(msg *tb.Message) {
	bot.Send(msg.Chat, "Looking for that topic in your list of subscriptions...")
	s := Sub{ChatID: int32(msg.Chat.ID), Topic: msg.Payload}
	err := bot.store.Delete(s.String(), s)
	if err != nil {
		bot.Send(msg.Chat, "sorry, failed to delete that: "+err.Error())
		return
	}

	bot.Send(msg.Chat, "OK deleted that")
}

func (bot *Bot) list(msg *tb.Message) {
	subs := []Sub{}
	err := bot.store.Find(&subs, badgerhold.Where("ChatID").Eq(msg.Chat.ID))

	if err != nil {
		bot.Send(msg.Chat, "sorry, failed to search: "+err.Error())
		return
	}

	topics := []string{}
	for _, s := range subs {
		topics = append(topics, s.Topic)
	}

	bot.Send(msg.Chat, "OK, you're following these things:"+strings.Join(topics, "\n"))
}

func (bot *Bot) findrep(msg *tb.Message) {
	search := msg.Payload

	if len(search) < 4 {
		bot.Send(msg.Chat, "sorry, you're gonna need to do more than 3 letters")
		return
	}

	names, err := bot.searchReps(search)

	if err != nil {
		bot.Send(msg.Chat, "sorry, failed to search: "+err.Error())
		return
	}

	if len(names) > 10 {
		bot.Send(msg.Chat, "sorry, too many results.  Please be more specific")
		return
	}

	if len(names) == 1 {
		bot.Send(msg.Chat, "Found them: "+names[0])
		return
	}

	bot.Send(msg.Chat, "OK, it's gotta be one of these:\n"+strings.Join(names, "\n"))
}

func (bot *Bot) setupCommands() {
	bot.Handle("/findrep", bot.findrep)
	bot.Handle("/list", bot.list)
	bot.Handle("/follow", bot.follow)
	bot.Handle("/unfollow", bot.unfollow)
}