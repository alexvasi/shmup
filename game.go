package main

type Game struct {
	world *World
	input *Input
	ship  *Ship

	stage     Stage
	stages    []Stage
	nextStage int
}

type Stage interface {
	Init(world *World)
	Update(dt float32, world *World) bool
}

func NewGame(world *World, input *Input) *Game {
	game := &Game{
		world: world,
		input: input,
	}

	game.ship = NewShip(Human, &PlayerModel)
	game.world.AddShips(game.ship)

	starfield := NewStarfield(world.Size.X(), world.Size.Y())
	game.world.AddObjects(starfield)

	game.stages = []Stage{
		&IntroStage{Stars: starfield, Ship: game.ship},
		&CargoStage{MinCount: 2, MaxCount: 2},
		&CargoStage{MinCount: 5, MaxCount: 8},
		&CargoStage{MinCount: 6, MaxCount: 8},
		&FighterStage{},
		&CargoStage{MinCount: 6, MaxCount: 8},
		&FinalStage{},
		&OutroStage{Stars: starfield, Ship: game.ship},
	}

	return game
}

func (game *Game) Update(dt float32) {
	if game.ship.Race == Human {
		game.ship.Control(game.input.Dir, game.input.Fire)
	}

	if game.stage == nil {
		game.initNextStage()
	}

	if game.ship.IsDead {
		game.nextStage = 1 // skip intro
		game.stage = &DeathStage{Ship: game.ship}
		game.stage.Init(game.world)
	}

	ok := game.stage.Update(dt, game.world)
	if !ok {
		game.stage = nil
	}

	game.world.Update(dt)
}

func (game *Game) initNextStage() {
	game.stage = game.stages[game.nextStage]
	game.nextStage = (game.nextStage + 1) % len(game.stages)

	game.stage.Init(game.world)
}
