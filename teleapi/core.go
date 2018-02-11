package teleapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type method string

const (
	apiURL          string = "https://api.telegram.org/bot"
	sendMessageMthd method = "sendMessage"
	getUpdates      method = "getUpdates"
)

// TODO: move channel from bot stuct... bot can have more than one update subscriptions for different message types
type bot struct {
	token string
}

func (bot *bot) makeURL(m method) string {
	return fmt.Sprintf("%s%s/%s", apiURL, bot.token, m)
}

// Bot ...
type Bot interface {
	SendMessage(SendMessageReq) error
	Listen() <-chan *Update
}

// NewBot ...
func NewBot(t string) Bot {
	bot := bot{token: t}
	return &bot
}

// SendMessage ...
// func (bot *bot) SendMessage(chatID int64, text string, disableWebPagePreview bool) error {
func (bot *bot) SendMessage(sendMessageReq SendMessageReq) error {
	// sendMessageReq := sendMessageReq{ChatID: chatID, Text: text, DisableWebPagePreview: disableWebPagePreview}
	jsonReq, err := json.Marshal(sendMessageReq)
	log.Printf("message to send: %s\n", jsonReq)
	if err != nil {
		log.Printf("[Error] SendMessage: can not marshal json request: %s\n", err)
		return err
	}
	endPnt := bot.makeURL(sendMessageMthd)
	req, err := http.NewRequest(http.MethodPost, endPnt, bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Printf("[Error] in build req: %s", err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Error] in send req: %s", err.Error())
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Warning] can not read api answer: {method: %s, data:%s}, err: %s", sendMessageMthd, jsonReq, err)
	}
	return nil
}

func (bot *bot) Listen() <-chan *Update {
	updateCh := make(chan *Update, 100)
	go doUpdates(bot, updateCh)
	return updateCh
}

func doUpdates(bot *bot, updateCh chan<- *Update) {
	endPnt := bot.makeURL(getUpdates)
	var currenOffset int64
	for {
		jsonStr := fmt.Sprintf(`{"offset":%d, "timeout": 60}`, currenOffset+1)
		jsonBlob := []byte(jsonStr)
		req, err := http.NewRequest(http.MethodPost, endPnt, bytes.NewBuffer(jsonBlob))
		if err != nil {
			log.Printf("[Warning] can not getUpdates: %s", err.Error())
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[Warning] in send req: %s", err.Error())
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Printf("[Warning] http status != 200, statusCode: %d\n", resp.StatusCode)
			continue
		}
		respBlob, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[Warning] can not read api answer: {method: %s, data:%s}, err: %s\n", getUpdates, jsonBlob, err)
		}
		log.Printf("[Debug] update data is: %s\n", respBlob)
		var result getUpdatesResp
		err = json.Unmarshal(respBlob, &result)
		if err != nil {
			log.Printf("[Warning] can not unmarshal resp: %s\n", err.Error())
			log.Printf("[Data] json is: %s\n", respBlob)
			continue
		}
		if !result.Ok {
			log.Printf("[Warning] result not ok\n")
			log.Printf("[Data] json is: %+v\n", result)
		}
		for _, update := range result.Result {
			updateCh <- update
			if update.UpdateID > currenOffset {
				currenOffset = update.UpdateID
			}
		}

	}
}
