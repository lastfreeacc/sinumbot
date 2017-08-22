package main

import (
	"log"
	"fmt"
	"strings"
	"io/ioutil"
	"encoding/json"
	"github.com/lastfreeacc/sinumbot/store"
	"github.com/lastfreeacc/sinumbot/teleapi"
)

type cmd string

const (
	confFilename = "sinumbot.conf.json"
	startCmd     	cmd = "/start"
	listCmd			cmd = "/l"
	tagCmd 			cmd = "/t"
)

func (c cmd) isMe(msg string) bool {
	return strings.HasPrefix(msg, string(c))
}

var (
	conf     = make(map[string]interface{})
	botToken string
	bot      teleapi.Bot
	botStore = store.NewBoltStore()
)

func main() {

	// inline test for bolt storage
	// botStore := store.NewBoltStore()
	// botStore.SaveMemo(123, "memo1")
	// botStore.SaveMemo(123, "memo2")
	// u, _ := botStore.GetUser(123)
	// log.Printf("%v", u)
	// ...

	myInit()
	upCh := bot.Listen()
	for update := range upCh {
		cmd := update.Message.Text
		switch true {
		case startCmd.isMe(cmd):
			doStrart(update)
		case listCmd.isMe(cmd):
			doList(update)
		case tagCmd.isMe(cmd):
			doTag(update)
		default:
			doFeed(update)
		}
	}

}

func myInit() {
	readMapFromJSON(confFilename, &conf)
	botToken, ok := conf["botToken"]
	if !ok || botToken == "" {
		log.Fatalf("[Error] can not find botToken in config file: %s\n", confFilename)
	}
	bot = teleapi.NewBot(botToken.(string))

}

func readMapFromJSON(filename string, mapVar *map[string]interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("[Warning] can not read file '%s'\n", filename)
	}
	if err := json.Unmarshal(data, mapVar); err != nil {
		log.Fatalf("[Warning] can not unmarshal json from file '%s'\n", filename)
	}
	log.Printf("[Info] read data from file: %s:\n%v\n", filename, mapVar)
}

func doStrart(update *teleapi.Update) {
	msg := fmt.Sprint(
		`Hello, i am ur pocket!
	feed me ur urls, i'll show them later
	Usage:
	share url for read later
	/l - list urls in pocket
	/t <>`)
	bot.SendMessage(update.Message.Chat.ID, msg, false)
}

func doList(update *teleapi.Update) {
	userID := update.Message.From.ID
	u, err := botStore.GetUser(userID)
	if err != nil {
		log.Printf("[Warn] some trobles in doList, err: %s\n", err)
		bot.SendMessage(update.Message.Chat.ID, "Sorry, we have some troubles!", false)
		return
	}
	
	msg := `ur memos is:
	-------
	`
	for _, memo := range u.Memos {
		log.Printf("[Info] memo is: %+v\n", *memo)
		msg = msg + memo.Feed + "\n-------\n" 
	}
	bot.SendMessage(update.Message.Chat.ID, msg, true)
}

func doFeed(update *teleapi.Update) {
	feed := update.Message.Text
	userID := update.Message.From.ID
	tags := teleapi.TagsFromMessage(update.Message)
	memo := store.NewMemo(feed, tags)
	botStore.SaveMemo(userID, &memo)
	msg := "ok, i'll show it later"
	bot.SendMessage(update.Message.Chat.ID, msg, false)
}

func doTag(update *teleapi.Update) {
	tags := strings.Fields(update.Message.Text)
	tags = tags[1:]
	for i, tag := range tags {
		if !strings.HasPrefix(tag, "#") {
			tags[i] = "#" + tag
		}
	}
	log.Printf("[Info] tags: %+v\n", tags)
	userID := update.Message.From.ID
	u, err := botStore.GetUser(userID)
	if err != nil {
		log.Printf("[Warn] some trobles in doTag, err: %s\n", err)
		bot.SendMessage(update.Message.Chat.ID, "Sorry, we have some troubles!", false)
		return
	}
	taggedMemos := containsAllTags(u.Memos, tags)
	msg := `ur memos tagged :
	-------
	`
	for _, memo := range taggedMemos {
		msg = msg + memo.Feed + "\n-------\n" 
	}
	bot.SendMessage(update.Message.Chat.ID, msg, true)
}

func containsAllTags(memos []*store.Memo, tags []string) []*store.Memo{
	res := make([]*store.Memo, 0, len(memos))
	for _, memo := range memos {
		if containsAll(memo.Tags, tags) {
			res = append(res, memo)
		}
	}
	return res
}

func containsAll(srcSl, testSl []string) bool {
	for _, test := range testSl {
		if !contains(srcSl, test) {
			return false
		}
	}
	return true
}

func contains(srcSl []string, e string) bool {
	for _, src := range srcSl {
		if src == e {
			return true
		}
	}
	return false
}