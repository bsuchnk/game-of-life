package main

const (
	fps = 15

	width  = 720
	height = 360

	rows    = 360
	columns = 180

	threshold = 0.15
)

var (
	square = []float32{
		0, 1, 0,
		0, 0, 0,
		1, 0, 0,

		0, 1, 0,
		1, 1, 0,
		1, 0, 0,
	}
)
