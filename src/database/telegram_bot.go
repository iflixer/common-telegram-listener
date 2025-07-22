package database

import (
	"log"
)

type TelegramBot struct {
	ID           int
	Name         string
	ListenURL    string
	Type         string
	BotURL       string
	AppURL       string
	Description  string
	Token        string
	GaTrackingID string
	GaSecret     string
	SearchURL    string
	UpdatedAt    string
	Published    bool
}

func (TelegramBot) TableName() string {
	return "telegram_bot"
}

func LoadBots(dbService *Service) (bots []TelegramBot, err error) {
	if err = dbService.DB.Model(&TelegramBot{}).Where("published=1").Where("listen_url!=''").Find(&bots).Error; err != nil {
		log.Println("Failed to load bots:", err)
		return
	}
	log.Println("bots found:", len(bots))
	if len(bots) == 0 {
		log.Println("No published bots found.")
	}
	return
}

func (c *TelegramBot) Load(dbService *Service, id int) (err error) {
	return dbService.DB.Where("id=?", id).Limit(1).First(&c).Error
}
