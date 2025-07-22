package listener

import (
	"log"
	"sync"
	"telegram-listener/database"
	"telegram-listener/reaction"
	"telegram-listener/sender"
	"time"
)

type Service struct {
	mu              sync.RWMutex
	bots            map[int]*FlixBot
	dbService       *database.Service
	senderService   *sender.Service
	reactionService *reaction.Service
	updatePeriod    time.Duration
}

func NewService(dbService *database.Service, senderService *sender.Service, reactionService *reaction.Service) (s *Service, err error) {
	s = &Service{
		dbService:       dbService,
		updatePeriod:    time.Second * 60,
		bots:            make(map[int]*FlixBot),
		senderService:   senderService,
		reactionService: reactionService,
	}

	err = s.loadData()

	go s.loadWorker()

	return
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
	s.mu.Lock()
	defer s.mu.Unlock()
	telegramBots, err := database.LoadBots(s.dbService)
	if err != nil {
		log.Println("Failed to load bots:", err)
		return err
	}

	for _, botNew := range telegramBots {
		if botOld, found := s.bots[botNew.ID]; found { // update old bot?
			log.Printf("update bot %d? ", botOld.telegramBot.ID)
			if botOld.telegramBot.UpdatedAt == botNew.UpdatedAt {
				log.Println("no")
				continue // don't need to restart bot
			}
			log.Println("yes")
			botOld.Stop()
		}
		newFlixBot := &FlixBot{
			telegramBot: botNew,
		}
		err = newFlixBot.Register(s.dbService, s.reactionService)
		if err != nil {
			log.Println("failed to register bot:", botNew.ID, botNew.Name, err)
			continue
		}

		s.bots[botNew.ID] = newFlixBot
		log.Println("bot registered and started:", botNew.ID, botNew.Name)
	}
	log.Println("bots loaded:", len(s.bots))
	log.Println("= loading bots...DONE")
	return
}
