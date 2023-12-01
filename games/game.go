package games

type Player interface {
	GetName() string
	GetScore() int64
	Input()
}

type Game interface {
	AddPlayer(player Player) error

	Init() error

	ReadyToGo() bool

	Start()

	Pause()

	Over()
}
