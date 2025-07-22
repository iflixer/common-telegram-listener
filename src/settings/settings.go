package settings

import (
	"fmt"
	"log"
	"sync"
	"telegram-listener/database"
	"time"
)

type Service struct {
	mu           sync.RWMutex
	dbService    *database.Service
	updatePeriod time.Duration
	settings     []*database.Setting
}

func NewService(dbService *database.Service, updatePeriod int) (s *Service, err error) {

	s = &Service{
		dbService:    dbService,
		updatePeriod: time.Duration(updatePeriod),
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

func (s *Service) GetOne(botID int, command, part string) (setting *database.Setting, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, s := range s.settings {
		if s.BotID == botID && s.Command == command && s.Part == part {
			return s, nil
		}
	}
	return nil, fmt.Errorf("setting not found")
}

func (s *Service) loadData() (err error) {
	settings, err := database.GetAllSettings(s.dbService)
	if err != nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.settings = settings
	return
}
