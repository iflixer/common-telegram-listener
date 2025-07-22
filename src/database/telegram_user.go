package database

import (
	"time"

	"gorm.io/gorm/clause"
)

type TelegramUser struct {
	ID             int
	BotID          int
	TgID           int64
	LastCommand    string
	Counter        int
	UpdatedAt      *time.Time
	LastActivityAt *time.Time
	PushID         int
	Disabled       bool
	PushTime       *time.Time
}

func (c *TelegramUser) TableName() string {
	return "telegram_user"
}

func UpdateUsersPushID(dbService *Service, userIds []int64, pushID int) (err error) {
	if len(userIds) == 0 {
		return nil
	}
	return dbService.DB.Model(&TelegramUser{}).Where("tg_id IN (?)", userIds).Update("push_id", pushID).Error
}

func UpsertUser(dbService *Service, botID int, tgID int64, lastCommand string) (err error) {
	user := &TelegramUser{
		BotID:       botID,
		TgID:        tgID,
		LastCommand: lastCommand,
	}
	err = dbService.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "bot_id"}, {Name: "tg_id"}}, // уникальный ключ
		DoUpdates: clause.AssignmentColumns([]string{"last_command"}), // поля для обновления
	}).Create(user).Error

	return
}
