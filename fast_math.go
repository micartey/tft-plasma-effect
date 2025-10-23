package main

import (
	"math"
)

const pi = float32(math.Pi)

func sin(x float32) float32 {
	const (
		pi       float32 = 3.1415926535
		twoPi    float32 = 6.283185307
		invTwoPi float32 = 0.15915494309
		c        float32 = 0.40528473456
	)

	y := x * invTwoPi
	y = (y - float32(math.Floor(float64(y)))) * twoPi

	var sign float32 = 1.0
	if y > pi {
		y = y - pi
		sign = -1.0
	}

	return sign * c * y * (pi - y)
}

func sqrt(x float32) float32 {
	if x == 0 {
		return 0
	}

	i := (math.Float32bits(x) + 0x3FF00000) >> 1
	guess := math.Float32frombits(i)

	guess = (guess + x/guess) * 0.5
	return guess
}

func clamp(v, min, max float32) float32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func atan2(y, x float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}
