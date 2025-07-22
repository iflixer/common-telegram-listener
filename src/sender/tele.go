package sender

import (
	"encoding/json"
	"errors"

	"gopkg.in/telebot.v4"
)

type InlineMenu []struct {
	Row []struct {
		Title string `json:"title"`
		Value string `json:"value"`
	} `json:"row"`
}

func (s *Service) createInlineMenu(inlineMenuJson string) (inline *telebot.ReplyMarkup, err error) {
	inline = &telebot.ReplyMarkup{}
	if inlineMenuJson == "" {
		return nil, errors.New("empty js for inline menu")
	}

	inlineMenu := InlineMenu{}

	err = json.Unmarshal([]byte(inlineMenuJson), &inlineMenu)
	if err != nil {
		return nil, err
	}

	rows := []telebot.Row{}
	for _, row := range inlineMenu {
		var rowButtons []telebot.Btn
		for _, btn := range row.Row {
			if btn.Value == "" {
				return nil, errors.New("empty value for button: " + btn.Title)
			}
			rowButtons = append(rowButtons, inline.URL(btn.Title, btn.Value))
		}
		rows = append(rows, inline.Row(rowButtons...))
	}

	inline.Inline(rows...)
	return
}

type ReplyMenu []struct {
	Row []struct {
		Title string `json:"title"`
		Value string `json:"value"`
	} `json:"row"`
}

func (s *Service) createReplyMenu(replyMenuJson string) (inline *telebot.ReplyMarkup, err error) {
	inline = &telebot.ReplyMarkup{}
	if replyMenuJson == "" {
		return nil, errors.New("empty js for reply menu")
	}

	replyMenu := ReplyMenu{}

	err = json.Unmarshal([]byte(replyMenuJson), &replyMenu)
	if err != nil {
		return nil, err
	}

	rows := []telebot.Row{}
	for _, row := range replyMenu {
		var rowButtons []telebot.Btn
		for _, btn := range row.Row {
			if btn.Value == "" {
				return nil, errors.New("empty value for button: " + btn.Title)
			}
			rowButtons = append(rowButtons, inline.URL(btn.Title, btn.Value))
		}
		rows = append(rows, inline.Row(rowButtons...))
	}

	inline.Reply(rows...)

	return
}
