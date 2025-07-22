package database

type TelegramAudience struct {
	ID        int
	Type      int // 0 - push to users, 1 - push to channel
	Name      string
	Query     string
	ChannelID int64
}

const AudienceTypeUsers = 0
const AudienceTypeChannel = 1

func (c *TelegramAudience) TableName() string {
	return "telegram_audience"
}

func (c *TelegramAudience) Load(dbService *Service, id int) (err error) {
	return dbService.DB.Where("id=?", id).Limit(1).First(&c).Error
}
