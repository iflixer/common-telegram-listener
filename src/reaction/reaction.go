package reaction

import (
	"log"
	"strings"
	"sync"
	"telegram-listener/database"
	"telegram-listener/sender"
	"telegram-listener/settings"
	"time"

	"gopkg.in/telebot.v4"
)

type Service struct {
	mu              sync.RWMutex
	reactions       []*database.TelegramBotReaction
	dbService       *database.Service
	updatePeriod    time.Duration
	senderService   *sender.Service
	settingsService *settings.Service
}

func NewService(dbService *database.Service, senderService *sender.Service, settingsService *settings.Service) (s *Service, err error) {
	s = &Service{
		dbService:       dbService,
		senderService:   senderService,
		updatePeriod:    time.Second * 60,
		reactions:       []*database.TelegramBotReaction{},
		settingsService: settingsService,
	}

	err = s.loadData()

	go s.loadWorker()

	return
}

func (s *Service) getAllReactions() []*database.TelegramBotReaction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.reactions
}

func (s *Service) getOne(id int) *database.TelegramBotReaction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, reaction := range s.reactions {
		if reaction.ID == id {
			return reaction
		}
	}
	return nil
}

func (s *Service) RegisterReactions(botID int, appURL string, tbot *telebot.Bot) {
	reactions := s.getAllReactions()
	for _, reaction := range reactions {
		if reaction.BotID == botID && reaction.Handle != "" {
			log.Println("Registering reaction:", reaction.Handle, "for bot ID:", botID)
			// add or update user
			if strings.HasPrefix(reaction.Handle, "/") { // command
				tbot.Handle(reaction.Handle, func(c telebot.Context) error {
					database.UpsertUser(s.dbService, botID, c.Sender().ID, reaction.Handle)
					log.Printf("bot:%d received command from %d (%s): %s", botID, c.Sender().ID, c.Sender().Username, reaction.Handle)
					for reaction != nil {
						err := s.senderService.SendText(tbot, c.Sender().ID, reaction.Answer, reaction.InlineMenu, reaction.ReplyMenu)
						if err != nil {
							log.Printf("❌ Failed to send message ID %d: %v", reaction.ID, err)
							return err
						}
						if reaction.AdditionalMessageID > 0 {
							reaction = s.getOne(reaction.AdditionalMessageID)
						} else {
							break
						}
					}
					return nil
				})
			} else { // text
				tbot.Handle(telebot.OnText, func(c telebot.Context) error {
					msg := c.Text() // Получаем текст сообщения
					database.UpsertUser(s.dbService, botID, c.Sender().ID, msg)
					log.Printf("bot:%d received message from %d (%s): %s", botID, c.Sender().ID, c.Sender().Username, msg)
					log.Printf("Searching for reaction: %+v", reaction)
					inlineMenu, err := s.searchPosts(reaction.Handle, appURL, msg)
					if err != nil {
						log.Printf("❌ Failed to search posts for reaction ID %d: %v", reaction.ID, err)
						answer := "Search error. Try later"
						settingNotFound, err := s.settingsService.GetOne(botID, "search", "not_found")
						if err == nil && settingNotFound != nil {
							answer = settingNotFound.Content
						}
						err = s.senderService.SendText(tbot, c.Sender().ID, answer, inlineMenu, "")
						if err != nil {
							log.Printf("❌ Failed to send message ID %d: %v", reaction.ID, err)
						}
						return err
					}

					if inlineMenu != "" {
						answer := strings.ReplaceAll(reaction.Answer, "[orig-message]", msg)
						answer = strings.ReplaceAll(answer, "[user-username]", c.Sender().Username)
						err = s.senderService.SendText(tbot, c.Sender().ID, answer, inlineMenu, reaction.ReplyMenu)
						if err != nil {
							log.Printf("❌ Failed to send message ID %d: %v", reaction.ID, err)
						}
						return err
					}

					if reaction.AdditionalMessageID > 0 {
						log.Printf("No posts found for reaction ID %d with handle %s", reaction.ID, reaction.Handle)
						reaction2 := s.getOne(reaction.AdditionalMessageID)
						answer := strings.ReplaceAll(reaction2.Answer, "[orig-message]", msg)
						answer = strings.ReplaceAll(answer, "[user-username]", c.Sender().Username)
						err = s.senderService.SendText(tbot, c.Sender().ID, answer, reaction2.InlineMenu, reaction2.ReplyMenu)
						if err != nil {
							log.Printf("❌ Failed to send message ID %d: %v", reaction2.ID, err)
						}
						return err
					}

					return nil
				})
			}
		}
	}
}

func (s *Service) loadWorker() {
	for {
		time.Sleep(time.Second * s.updatePeriod)
		if err := s.loadData(); err != nil {
			log.Println(err)
		}
	}
}

func (s *Service) loadData() (err error) {
	log.Println("Loading reactions from database")
	s.mu.Lock()
	defer s.mu.Unlock()
	reactions, err := database.LoadReactions(s.dbService)
	if err != nil {
		log.Println("Failed to load reactions:", err)
		return err
	}
	s.reactions = reactions

	return
}
