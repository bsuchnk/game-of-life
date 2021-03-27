package main

import (
	"math/rand"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type cell struct {
	drawable uint32 // VAO

	alive     bool
	aliveNext bool

	ageDead int
	color   float32

	x int
	y int
}

func (c *cell) draw() {
	if c.ageDead >= 20 {
		return
	}

	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}

func (c *cell) getNextState(cells [][]*cell) {
	n, col := c.countNeighbours(cells)
	if c.alive {
		if n == 2 || n == 3 {
			c.aliveNext = true
		} else {
			c.aliveNext = false
		}
	} else {
		if n == 3 {
			c.aliveNext = true
			col = col + rand.Float32()*16 - 8 + 360
			if col >= 360 {
				col -= 360
			}
			c.color = col
		} else {
			c.aliveNext = false
		}
	}
}

func (c *cell) updateState() {
	c.alive = c.aliveNext

	if c.alive {
		c.ageDead = 0
	} else {
		c.ageDead++
	}
}

func (c *cell) countNeighbours(cells [][]*cell) (int32, float32) {
	var count int32
	var color float32
	var colors []float32

	for i := c.x - 1; i <= c.x+1; i++ {
		for j := c.y - 1; j <= c.y+1; j++ {
			if i == c.x && j == c.y {
				continue
			}

			if i >= 0 && i < rows && j >= 0 && j < columns {
				if cells[i][j].alive {
					count++
					color += cells[i][j].color
					colors = append(colors, cells[i][j].color)
				}
			}
		}
	}
	if count == 0 {
		return 0, 0
	}

	if count == 2 {
		return count, colors[rand.Int31()%2]
	}

	return count, color / float32(count)
}

func newCell(x, y int) *cell {
	points := make([]float32, len(square))
	copy(points, square)

	var rowRes float32 = 2.0 / rows
	var colRes float32 = 2.0 / columns

	for i := 0; i < len(points); i++ {
		switch i % 3 {
		case 0:
			points[i] = -1 + points[i]*rowRes + float32(x)*rowRes
		case 1:
			points[i] = -1 + points[i]*colRes + float32(y)*colRes
		}
	}

	return &cell{
		drawable: makeVao(points),

		x: x,
		y: y,

		ageDead: 20,
		color:   rand.Float32() * 360,
	}
}
