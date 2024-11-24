package controllers

import (
	"context"
	"time"
)

type ActionHandel interface {
	Handle(ctx context.Context, action *Action) error
}

type Action struct {
	Name    string
	value   int64
	handles []ActionHandel
}

func (a *Action) SetValue(value int64) {
	a.value = value
}

func (a *Action) BindHandle(handel ActionHandel) {
	a.handles = append(a.handles, handel)
}

type Controller interface {
	ActionMap() map[string]Action
}

type StateHandel interface {
	Start() (chan<- State, error)
}

type State struct {
	Name             string
	value            float64
	onChangeChannels []chan<- State
	timingWheel      time.Ticker
	timingChannels   []chan<- State
}

func (s *State) SetValue(value float64) {
	s.value = value
	s.onChange()
}

func (s *State) Incr(val float64) {
	s.value += val
	s.onChange()
}

func (s *State) Decr(val float64) {
	s.value -= val
	s.onChange()
}

func (s *State) BindOnChangeHandle(handel StateHandel) error {
	ch, err := handel.Start()
	if err != nil {
		return err
	}
	s.onChangeChannels = append(s.onChangeChannels, ch)
	return nil
}

func (s State) onChange() {
	for _, ch := range s.onChangeChannels {
		go func(nch chan<- State) {
			nch <- s
		}(ch)
	}
}

func (s *State) BindTimingHandle(handel StateHandel) error {
	ch, err := handel.Start()
	if err != nil {
		return err
	}
	s.timingChannels = append(s.timingChannels, ch)
	return nil
}
