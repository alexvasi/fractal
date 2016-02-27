package main

import "math"

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

func NormalizeAngle(a float32) float32 {
	for a > math.Pi {
		a -= 2 * math.Pi
	}

	for a < math.Pi {
		a += 2 * math.Pi
	}

	return a
}
