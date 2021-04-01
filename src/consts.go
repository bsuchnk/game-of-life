package main

var (
	width  = 720
	height = 360

	rows    = 256
	columns = 128

	fps = 15

	threshold = 0.15

	square = []float32{
		0, 1, 0,
		0, 0, 0,
		1, 0, 0,

		0, 1, 0,
		1, 1, 0,
		1, 0, 0,
	}
)
