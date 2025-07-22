package database

import (
	"time"
)

type TelegramPush struct {
	ID            int
	Type          string // "text" or "photo"
	BotID         int
	StartedAt     *time.Time
	EndAt         *time.Time
	AudienceID    int
	Affected      int
	Command       string
	Status        string
	InlineButtons string
	MenuButtons   string
	ImageURL      string
	Text          string
}

func (c *TelegramPush) TableName() string {
	return "telegram_push"
}

func (c *TelegramPush) Load(dbService *Service, id int) (err error) {
	return dbService.DB.Where("id=?", id).Limit(1).First(&c).Error
}

func (c *TelegramPush) save(dbService *Service) (err error) {
	return dbService.DB.Save(c).Error
}

func (c *TelegramPush) SearchTasks(dbService *Service) (err error) {
	return dbService.DB.Where("started_at<NOW() AND status='ready'").Limit(1).Find(&c).Error
}

func (c *TelegramPush) SetStatusStart(dbService *Service) (err error) {
	c.Status = "started"
	now := time.Now().UTC()
	c.StartedAt = &now
	return c.save(dbService)
}

func (c *TelegramPush) SetStatusError(dbService *Service) (err error) {
	c.Status = "error"
	now := time.Now().UTC()
	c.EndAt = &now
	return c.save(dbService)
}

func (c *TelegramPush) SetStatusDone(dbService *Service) (err error) {
	c.Status = "done"
	now := time.Now().UTC()
	c.EndAt = &now
	return c.save(dbService)
}

func (c *TelegramPush) UpdateAffected(dbService *Service, affected int, increase int) (err error) {
	if increase > 0 {
		c.Affected += increase
	} else if affected > 0 {
		c.Affected = affected
	} else {
		return nil // No change in affected count
	}
	return c.save(dbService)
}
