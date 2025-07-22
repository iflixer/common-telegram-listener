package database

type Setting struct {
	ID        int
	BotID     int
	Command   string
	Part      string
	Orderby   int
	Content   string
	ImageURL  string
	Link      string
	Published bool
}

func (c *Setting) TableName() string {
	return "telegram_settings"
}

func GetAllSettings(dbService *Service) (settings []*Setting, err error) {
	err = dbService.DB.Where("published=1").Order("command,part,orderby").Find(&settings).Error
	return
}
