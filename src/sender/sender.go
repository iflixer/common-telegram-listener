package sender

import (
	"log"
	"telegram-listener/database"
	"telegram-listener/helper"

	"gopkg.in/telebot.v4"
)

type Service struct {
	// mu        sync.RWMutex
	dbService *database.Service
}

func NewService(dbService *database.Service) (s *Service, err error) {
	s = &Service{
		dbService: dbService,
	}
	return
}

func (s *Service) SendText(tbot *telebot.Bot, chatID int64, msg string, inlineMenuJson string, replyMenuJson string) (err error) {
	recipient := &telebot.User{ID: chatID}

	msg = helper.SanitizeTelegramHTML(msg)
	if inlineMenuJson != "" {
		inlineMenu, err := s.createInlineMenu(inlineMenuJson)
		if err != nil {
			log.Println("Failed to create inline menu:", err)
		} else {
			_, err = tbot.Send(recipient, msg, inlineMenu, telebot.ModeHTML)
			return err
		}
	}
	if replyMenuJson != "" {
		replyMenu, err := s.createReplyMenu(replyMenuJson)
		if err != nil {
			log.Println("Failed to create reply menu:", err)
		} else {
			_, err = tbot.Send(recipient, msg, replyMenu, telebot.ModeHTML)
			return err
		}
	}

	_, err = tbot.Send(recipient, msg, telebot.ModeHTML)
	return
}

func (s *Service) SendPhoto(tbot *telebot.Bot, inlineMenuJson string, chatID int64, msg string, url string) (err error) {
	recipient := &telebot.User{ID: chatID}
	photo := &telebot.Photo{
		File:    telebot.FromURL(url),
		Caption: helper.SanitizeTelegramHTML(msg),
	}

	if inlineMenuJson != "" {
		inlineMenu, err := s.createInlineMenu(inlineMenuJson)
		if err == nil {
			_, err = tbot.Send(recipient, photo, inlineMenu, telebot.ModeHTML)
			return err
		}
	}
	_, err = tbot.Send(recipient, photo, telebot.ModeHTML)
	return
}

func (s *Service) SendVideo(tbot *telebot.Bot, inlineMenuJson string, chatID int64, msg string, url string) (err error) {
	recipient := &telebot.User{ID: chatID}
	video := &telebot.Video{
		File:    telebot.FromURL(url),
		Caption: helper.SanitizeTelegramHTML(msg),
	}

	if inlineMenuJson != "" {
		inlineMenu, err := s.createInlineMenu(inlineMenuJson)
		if err == nil {
			_, err = tbot.Send(recipient, video, inlineMenu, telebot.ModeHTML)
			return err
		}
	}
	_, err = tbot.Send(recipient, video, telebot.ModeHTML)
	return
}
