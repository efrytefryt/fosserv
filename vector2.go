package main

import (
	"fmt"
	"math"
)

type Vector2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func Add(v1, v2 Vector2) (result Vector2) {
	result.X = v1.X + v2.X
	result.Y = v1.Y + v2.Y
	return
}

func Subtract(v1, v2 Vector2) (result Vector2) {
	result.X = v1.X - v2.X
	result.Y = v1.Y - v2.Y
	return
}

func Multiply(v Vector2, m float64) (result Vector2) {
	result.X = v.X * m
	result.Y = v.Y * m
	return
}

func Distance(v1, v2 Vector2) (result float64) {
	result = math.Sqrt( math.Pow(v1.X - v2.X, 2) + math.Pow(v1.Y - v2.Y, 2))
	return
}

func Magnitude(v Vector2) (result float64) {
	var emptyVector Vector2
	result = Distance(v, emptyVector)
	return
}

func Normalize(v *Vector2) {
	var magnitude float64 = Magnitude(*v)
	if magnitude != 0 {
		v.X = v.X / magnitude
		v.Y = v.Y / magnitude
	}
}

func (vector Vector2) String() string {
    return fmt.Sprintf("[%f, %f]", vector.X, vector.Y)
}