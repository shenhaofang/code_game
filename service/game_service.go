package service

import (
	"sync"

	svc "github.com/judwhite/go-svc"
)

type GameService struct {
	once sync.Once
}

func NewGameService() *GameService {
	return new(GameService)
}

// Init implements svc.Service.
func (s *GameService) Init(osEnv svc.Environment) error {

}

// Start implements svc.Service.
func (s *GameService) Start() error {
	panic("unimplemented")
}

// Stop implements svc.Service.
func (s *GameService) Stop() error {
	s.once.Do(func() {

	})
	return nil
}
