package database

import (
	"log"
)

type TelegramBotReaction struct {
	ID                  int
	BotID               int
	AdditionalMessageID int
	Handle              string
	Answer              string
	InlineMenu          string
	ReplyMenu           string
	Published           bool
}

func (c *TelegramBotReaction) TableName() string {
	return "telegram_bot_reaction"
}

func LoadReactions(dbService *Service) (reactions []*TelegramBotReaction, err error) {
	if err = dbService.DB.Model(&TelegramBotReaction{}).Where("published = ?", true).Find(&reactions).Error; err != nil {
		log.Println("Failed to load reactions:", err)
		return
	}
	log.Println("Reactions found:", len(reactions))
	if len(reactions) == 0 {
		log.Println("No reactions found")
	}
	return
}
