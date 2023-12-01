package games

type ShootGame interface {
	Game
}

type Position struct {
	X, Y, Z int
}

type Bullet struct {
	Position Position
	Speed    int
	Forward  Position
}

type WeaponOperates interface {
	Do()
}

type Weapon struct {
	Name          string                    `json:"name"`
	Type          int                       `json:"type"`
	IsConsumables bool                      `json:"is_consumables"`
	Stock         int                       `json:"stock"`
	Operates      map[string]WeaponOperates `json:"operates"`
}

func (w *Weapon) Consume(count int) {
	if !w.IsConsumables {
		return
	}
	w.Stock -= count
	if w.Stock < 0 {
		panic("cost more than stocks")
	}
}

type PlayerUnit interface {
	AddWeapon(w Weapon, operate string)
	Attack(weaponIdx int)
}

type ShootControl interface {
}
