package main

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/lucasb-eyer/go-colorful"
)

func Max(a float32, numbers ...float32) float32 {
	for _, n := range numbers {
		if n > a {
			a = n
		}
	}
	return a
}

func Min(a float32, numbers ...float32) float32 {
	for _, n := range numbers {
		if n < a {
			a = n
		}
	}
	return a
}

func MaxVec2(v1 mgl.Vec2, vectors ...mgl.Vec2) mgl.Vec2 {
	for _, vec := range vectors {
		v1 = mgl.Vec2{Max(v1.X(), vec.X()), Max(v1.Y(), vec.Y())}
	}
	return v1
}

func MinVec2(v1 mgl.Vec2, vectors ...mgl.Vec2) mgl.Vec2 {
	for _, vec := range vectors {
		v1 = mgl.Vec2{Min(v1.X(), vec.X()), Min(v1.Y(), vec.Y())}
	}
	return v1
}

func HexColor(hex string, alpha float32) mgl.Vec4 {
	c, err := colorful.Hex(hex)
	if err != nil {
		panic(err)
	}
	return mgl.Vec4{float32(c.R), float32(c.G), float32(c.B), alpha}
}

func BlendColors(c1, c2 mgl.Vec4, t float32) mgl.Vec4 {
	from := colorful.Color{
		R: float64(c1[0]),
		G: float64(c1[1]),
		B: float64(c1[2]),
	}
	to := colorful.Color{
		R: float64(c2[0]),
		G: float64(c2[1]),
		B: float64(c2[2]),
	}
	r := from.BlendLab(to, float64(t)).Clamped()
	alpha := mgl.Clamp(c1[3]+(c2[3]-c1[3])*t, 0, 1)

	return mgl.Vec4{float32(r.R), float32(r.G), float32(r.B), alpha}
}

func CheckAABB(a, b mgl.Vec4) bool {
	aMinX, aMinY, aMaxX, aMaxY := a.Elem()
	bMinX, bMinY, bMaxX, bMaxY := b.Elem()

	if aMinX > bMaxX || aMaxX < bMinX {
		return false
	}
	if aMinY > bMaxY || aMaxY < bMinY {
		return false
	}

	return true
}

func SegmentIntersection(a1, a2, b1, b2 mgl.Vec2) (bool, mgl.Vec2) {
	aLine := a1.Vec3(1).Cross(a2.Vec3(1))
	bLine := b1.Vec3(1).Cross(b2.Vec3(1))
	pointH := aLine.Cross(bLine)

	if mgl.FloatEqual(pointH.Z(), 0) {
		return false, mgl.Vec2{} // lines are somewhat parallel
	}

	point := mgl.Vec2{pointH.X() / pointH.Z(), pointH.Y() / pointH.Z()}
	if a1.Sub(point).Dot(a2.Sub(point)) > 0 {
		return false, point // intersection is outside of [a1, a2]
	}
	if b1.Sub(point).Dot(b2.Sub(point)) > 0 {
		return false, point // intersection is outside of [b1, b2]
	}

	return true, point
}
