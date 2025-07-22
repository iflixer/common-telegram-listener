package database

type TelegramChannel struct {
	ID      int
	Name    string
	TgID    int64
	Comment string
}

func (c *TelegramChannel) TableName() string {
	return "telegram_channel"
}

func (c *TelegramChannel) Load(dbService *Service, id int64) (err error) {
	return dbService.DB.Where("id=?", id).Limit(1).First(&c).Error
}
