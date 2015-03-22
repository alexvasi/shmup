package main

import mgl "github.com/go-gl/mathgl/mgl32"

type ShipModel struct {
	color mgl.Vec3
	size  mgl.Vec2
	gun   mgl.Vec2
	hull  []mgl.Vec2
}

var PlayerModel = ShipModel{
	color: mgl.Vec3{45. / 255, 107. / 255, 178. / 255},
	size:  mgl.Vec2{90, 110},
	gun:   mgl.Vec2{0, 1},
	hull: []mgl.Vec2{
		mgl.Vec2{0, 1},
		mgl.Vec2{-1, -1},
		mgl.Vec2{1, -1},
	},
}
