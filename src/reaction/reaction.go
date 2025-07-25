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

	// tbot.Use(func(next telebot.HandlerFunc) telebot.HandlerFunc {
	// 	return func(c telebot.Context) error {
	// 		log.Printf("[DEBUG] ChatType: %s, Message from %s: %s", c.Chat().Type, c.Sender().Username, c.Text())
	// 		return next(c)
	// 	}
	// })

	reactions := s.getAllReactions()
	for _, reaction := range reactions {
		if reaction.BotID == botID && reaction.Handle != "" {
			log.Println("Registering reaction:", reaction.Handle, "for bot ID:", botID)
			// add or update user
			if strings.HasPrefix(reaction.Handle, "/") { // command
				tbot.Handle(reaction.Handle, func(c telebot.Context) error {
					msg := c.Text() // Получаем текст сообщения
					if len(msg) > 200 {
						msg = msg[:200] // Ограничиваем длину сообщения до 200 символов
					}
					database.UpsertUser(s.dbService, botID, c.Sender().ID, msg)
					log.Printf("bot:%d received command from %d (%s): %s", botID, c.Sender().ID, c.Sender().Username, msg)
					for reaction != nil {
						err := s.senderService.SendText(tbot, c.Sender().ID, reaction.Answer, reaction.InlineMenu, reaction.ReplyMenu)
						if err != nil {
							log.Printf("❌ Failed to send message by bot %d to %d: %v", botID, c.Sender().ID, err)
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
			}
		}
	}

	tbot.Handle(telebot.OnText, func(c telebot.Context) error {
		msg := c.Text() // Получаем текст сообщения
		if len(msg) > 200 {
			msg = msg[:200] // Ограничиваем длину сообщения до 200 символов
		}

		chatType := c.Chat().Type
		log.Println("Received message in chat type:", chatType, "from user:", c.Sender().ID, "with text:", msg)

		msgPrefix := ""
		switch chatType {
		case telebot.ChatPrivate:
			// Это личный чат (1:1 с пользователем)
			log.Println("Личное сообщение от пользователя:", c.Sender().Username)
		case telebot.ChatGroup, telebot.ChatSuperGroup:
			// Это группа или супергруппа
			log.Println("Сообщение в группе:", c.Chat().Title)
			msgPrefix = c.Sender().Username + ": "
		default:
			log.Println("Другой тип чата:", chatType)
		}

		database.UpsertUser(s.dbService, botID, c.Sender().ID, msg)
		log.Printf("bot:%d received message from %d (%s): %s", botID, c.Sender().ID, c.Sender().Username, msg)

		reaction := s.getOneByHandle(botID, msg)

		// может это команда без слеша?
		if reaction != nil {
			log.Printf("Non-slash command: %+v", reaction)
			log.Printf("Found reaction for message: %s", msg)
			err := s.senderService.SendText(tbot, c.Sender().ID, msgPrefix+reaction.Answer, reaction.InlineMenu, reaction.ReplyMenu)
			if err != nil {
				log.Printf("❌ Failed to send message by bot %d to %d: %v", botID, c.Sender().ID, err)
			}
			return err
		}

		reaction = s.getOneByHandle(botID, "https://")
		if reaction == nil {
			log.Printf("❌ Failed to load search reaction for botID %d", botID)
			return nil
		}
		inlineMenu, err := s.searchPosts(reaction.Handle, appURL, msg)
		if err != nil {
			log.Printf("❌ Failed to search posts for reaction ID %d: %v", reaction.ID, err)
			answer := "Search error. Try later"
			settingNotFound, err := s.settingsService.GetOne(botID, "search", "not_found")
			if err == nil && settingNotFound != nil {
				answer = settingNotFound.Content
			}
			err = s.senderService.SendText(tbot, c.Sender().ID, msgPrefix+answer, inlineMenu, "")
			if err != nil {
				log.Printf("❌ Failed to send message by bot %d to %d: %v", botID, c.Sender().ID, err)
			}
			return err
		}

		if inlineMenu != "" {
			answer := strings.ReplaceAll(reaction.Answer, "[orig-message]", msg)
			answer = strings.ReplaceAll(answer, "[user-username]", c.Sender().Username)
			err = s.senderService.SendText(tbot, c.Sender().ID, msgPrefix+answer, inlineMenu, reaction.ReplyMenu)
			if err != nil {
				log.Printf("❌ Failed to send message by bot %d to %d: %v", botID, c.Sender().ID, err)
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

func (s *Service) getOneByHandle(botID int, handle string) *database.TelegramBotReaction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, reaction := range s.reactions {
		if reaction.BotID == botID && strings.HasPrefix(reaction.Handle, handle) {
			return reaction
		}
	}
	return nil
}
