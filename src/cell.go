package main

import "github.com/go-gl/gl/v3.3-core/gl"

type cell struct {
	drawable uint32 // VAO

	alive     bool
	aliveNext bool

	ageDead int

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
	n := c.countNeighbours(cells)
	if c.alive {
		if n == 2 || n == 3 {
			c.aliveNext = true
		} else {
			c.aliveNext = false
		}
	} else {
		if n == 3 {
			c.aliveNext = true
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

func (c *cell) countNeighbours(cells [][]*cell) int {
	var count int

	for i := c.x - 1; i <= c.x+1; i++ {
		for j := c.y - 1; j <= c.y+1; j++ {
			if i == c.x && j == c.y {
				continue
			}

			if i >= 0 && i < rows && j >= 0 && j < columns {
				if cells[i][j].alive {
					count++
				}
			}
		}
	}

	return count
}

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
