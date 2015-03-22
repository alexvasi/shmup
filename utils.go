package main

import mgl "github.com/go-gl/mathgl/mgl32"

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
