package games

type Player interface {
	GetName() string
	GetScore() int64
	Input()
}

type Game interface {
	Name() string

	Init() error

	AddPlayer(player Player) error

	ReadyToGo() bool

	Start()

	Pause()

	Over()
}

var games = make(map[string]Game, 1000)

func Register(game Game) {
	games[game.Name()] = game
}
