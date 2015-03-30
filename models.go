package main

import mgl "github.com/go-gl/mathgl/mgl32"

type ShipModel struct {
	Size         mgl.Vec2
	Speed        float32
	Hp           int
	Color1       mgl.Vec4
	Color2       mgl.Vec4
	DmgColor     mgl.Vec4
	BlowupFactor float32
	Guns         []GunModel
	Engines      []EngineModel
	Hull         []mgl.Vec2
}

type GunModel struct {
	Pos        mgl.Vec2
	Rate       float32
	Speed      float32
	Size       mgl.Vec2
	Color      mgl.Vec4
	Sound      string
	SoundGain  float32
	SoundPitch float32
}

type EngineModel struct {
	Pos          mgl.Vec2
	Size         float32
	Rate         float32
	ParticleSize mgl.Vec2
	Color        mgl.Vec4
	TTL          float32
	MinVelocity  float32
}

var PlayerModel = ShipModel{
	Size:         mgl.Vec2{40, 50},
	Speed:        600,
	Hp:           1,
	Color1:       HexColor("#333333", 1),
	Color2:       HexColor("#555555", 1),
	DmgColor:     WhiteColor,
	BlowupFactor: 1,
	Guns: []GunModel{
		{
			Pos:        mgl.Vec2{-0.2, 0.8},
			Rate:       10,
			Speed:      1000,
			Size:       mgl.Vec2{7.5, 5},
			Color:      HexColor("#ffff00", 1),
			Sound:      "shoot_human",
			SoundGain:  0.4,
			SoundPitch: 0.8,
		},
		{
			Pos:   mgl.Vec2{0.2, 0.8},
			Rate:  10,
			Speed: 1000,
			Size:  mgl.Vec2{7.5, 5},
			Color: HexColor("#ffff00", 1),
		},
	},
	Engines: []EngineModel{
		{
			Pos:          mgl.Vec2{0, -0.1}, // near e,f
			Rate:         60,
			Size:         10,
			ParticleSize: mgl.Vec2{3, 8},
			Color:        HexColor("#01ffe6", 0.2),
			TTL:          0.1,
			MinVelocity:  0,
		},
		{
			Pos:          mgl.Vec2{-0.8, -1.0}, //d
			Rate:         30,
			Size:         1,
			ParticleSize: mgl.Vec2{1, 3},
			Color:        HexColor("#01ffe6", 0.2),
			TTL:          0.1,
			MinVelocity:  1,
		},
		{
			Pos:          mgl.Vec2{0.8, -1.0}, //g
			Rate:         30,
			Size:         1,
			ParticleSize: mgl.Vec2{1, 3},
			Color:        HexColor("#01ffe6", 0.2),
			TTL:          0.1,
			MinVelocity:  1,
		},
	},
	Hull: []mgl.Vec2{
		{+0.0, +1.0}, //a
		{-0.2, -0.2}, //e
		{+0.2, -0.2}, //f
		{+0.0, +1.0}, //a
		{-0.4, +0.8}, //b
		{-0.2, -0.2}, //e
		{+0.0, +1.0}, //a
		{+0.2, -0.2}, //f
		{+0.4, +0.8}, //j
		{-0.4, +0.8}, //b
		{-1.0, -0.3}, //c
		{-0.2, -0.2}, //e
		{+0.4, +0.8}, //j
		{+0.2, -0.2}, //f
		{+1.0, -0.3}, //i
		{-1.0, -0.3}, //c
		{-0.8, -1.0}, //d
		{-0.2, -0.2}, //e
		{+0.2, -0.2}, //f
		{+0.8, -1.0}, //g
		{+1.0, -0.3}, //i
	},
}

var CargoModel = ShipModel{
	Size:         mgl.Vec2{55, 60},
	Speed:        600,
	Hp:           5,
	Color1:       HexColor("#4f6952", 1),
	Color2:       HexColor("#7aff93", 1),
	DmgColor:     WhiteColor,
	BlowupFactor: 1,
	Engines: []EngineModel{
		{
			Pos:          mgl.Vec2{0, 0.1}, // near e,f
			Rate:         60,
			Size:         15,
			ParticleSize: mgl.Vec2{3, 16},
			Color:        HexColor("#b181ff", 0.2),
			TTL:          0.2,
			MinVelocity:  0,
		},
	},
	Hull: []mgl.Vec2{
		{+0.0, +1.0}, //a
		{-0.4, -0.0}, //e
		{+0.4, -0.0}, //f
		{+0.0, +1.0}, //a
		{-0.6, +0.8}, //b
		{-0.4, -0.0}, //e
		{+0.0, +1.0}, //a
		{+0.4, -0.0}, //f
		{+0.6, +0.8}, //j
		{-0.6, +0.8}, //b
		{-1.0, -0.1}, //c
		{-0.4, -0.0}, //e
		{+0.6, +0.8}, //j
		{+0.4, -0.0}, //f
		{+1.0, -0.1}, //i
		{-1.0, -0.1}, //c
		{-0.6, -1.0}, //d
		{-0.4, -0.0}, //e
		{+0.4, -0.0}, //f
		{+0.6, -1.0}, //g
		{+1.0, -0.1}, //i
	},
}

var CargoModel2 = ShipModel{
	Size:         mgl.Vec2{59, 65},
	Speed:        CargoModel.Speed,
	Hp:           7,
	Color1:       HexColor("#6b524d", 1),
	Color2:       HexColor("#ecff3e", 1),
	DmgColor:     CargoModel.DmgColor,
	BlowupFactor: 1.3,
	Engines:      CargoModel.Engines,
	Hull:         CargoModel.Hull,
}

var FighterModel = ShipModel{
	Size:         mgl.Vec2{40, 55},
	Speed:        600,
	Hp:           10,
	Color1:       HexColor("#a2904f", 1),
	Color2:       HexColor("#f3e52e", 1),
	DmgColor:     WhiteColor,
	BlowupFactor: 1,
	Guns: []GunModel{
		{
			Pos:        mgl.Vec2{-0.5, 0.9},
			Rate:       0.4,
			Speed:      200,
			Size:       mgl.Vec2{9, 6},
			Color:      HexColor("#e50c0c", 1),
			Sound:      "shoot",
			SoundGain:  0.7,
			SoundPitch: 1,
		},
		{
			Pos:   mgl.Vec2{0.5, 0.9},
			Rate:  0.4,
			Speed: 200,
			Size:  mgl.Vec2{9, 6},
			Color: HexColor("#e50c0c", 1),
		},
	},
	Hull: []mgl.Vec2{
		{+0.0, -1.0}, //a
		{-0.4, +0.0}, //e
		{+0.4, +0.0}, //f
		{+0.0, -1.0}, //a
		{-0.6, -0.8}, //b
		{-0.4, +0.0}, //e
		{+0.0, -1.0}, //a
		{+0.4, +0.0}, //f
		{+0.6, -0.8}, //j
		{-0.6, -0.8}, //b
		{-1.0, +0.1}, //c
		{-0.4, +0.0}, //e
		{+0.6, -0.8}, //j
		{+0.4, +0.0}, //f
		{+1.0, +0.1}, //i
		{-1.0, +0.1}, //c
		{-0.6, +1.0}, //d
		{-0.4, +0.0}, //e
		{+0.4, +0.0}, //f
		{+0.6, +1.0}, //g
		{+1.0, +0.1}, //i
	},
}

var ShooterModel = ShipModel{
	Size:         mgl.Vec2{40, 55},
	Speed:        600,
	Hp:           10,
	Color1:       HexColor("#a64b4b", 1),
	Color2:       HexColor("#35d7d0", 1),
	DmgColor:     WhiteColor,
	BlowupFactor: 1,
	Guns: []GunModel{
		{
			Pos:        mgl.Vec2{0, 0.7},
			Rate:       4,
			Speed:      200,
			Size:       mgl.Vec2{7, 7},
			Color:      HexColor("#ff1818", 1),
			Sound:      "shoot",
			SoundGain:  0.3,
			SoundPitch: 1,
		},
	},
	Hull: FighterModel.Hull,
}

var PapaModel = ShipModel{
	Size:         mgl.Vec2{120, 100},
	Speed:        600,
	Hp:           300,
	Color1:       HexColor("#5b5a59", 1),
	Color2:       HexColor("#c14848", 1),
	DmgColor:     WhiteColor,
	BlowupFactor: 2,
	Guns: []GunModel{
		{
			Pos:        mgl.Vec2{0, 0.75},
			Rate:       1,
			Speed:      800,
			Size:       mgl.Vec2{70, 45},
			Color:      HexColor("#fffd6a", 1),
			Sound:      "shoot_big",
			SoundGain:  1,
			SoundPitch: 1,
		},
	},
	Engines: []EngineModel{
		{
			Pos:          mgl.Vec2{0, -0.8}, // near a
			Rate:         150,
			Size:         100,
			ParticleSize: mgl.Vec2{3, 24},
			Color:        HexColor("#c14848", 0.6),
			TTL:          0.3,
			MinVelocity:  0,
		},
	},
	Hull: FighterModel.Hull,
}

var StarModel = []mgl.Vec2{
	{0, 1},
	{-1, -0.5},
	{1, -0.5},
	{0, -1},
	{1, 0.5},
	{-1, 0.5},
}
