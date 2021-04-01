package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	flag.IntVar(&width, "w", 512, "width")
	flag.IntVar(&height, "h", 512+48, "height")
	flag.IntVar(&rows, "r", 128, "rows")
	flag.IntVar(&columns, "c", 128, "columns")
	flag.IntVar(&fps, "fps", 15, "fps of the simulation")
	flag.Float64Var(&threshold, "t", 0.15, "randomness [0-1]")
	flag.Parse()

	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()
	gl.ClearColor(0.2, 0.1, 0.0, 1.0)

	cells := makeCells()

	for !window.ShouldClose() {
		t := time.Now()

		for x := range cells {
			for _, c := range cells[x] {
				c.updateState()
			}
		}

		for x := range cells {
			for _, c := range cells[x] {
				c.getNextState(cells)
			}
		}

		draw(cells, window, program)

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

func draw(cells [][]*cell, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	for _, row := range cells {
		for _, c := range row {
			ageUni := gl.GetUniformLocation(program, gl.Str("age\x00"))
			gl.Uniform1i(ageUni, int32(c.ageDead))
			colUni := gl.GetUniformLocation(program, gl.Str("H\x00"))
			gl.Uniform1f(colUni, c.color)

			c.draw()
		}
	}

	glfw.PollEvents()
	window.SwapBuffers()
}

func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

//--- grid ---

func makeCells() [][]*cell {
	rand.Seed(time.Now().UnixNano())

	cells := make([][]*cell, rows)
	for i := 0; i < rows; i++ {
		cells[i] = make([]*cell, columns)
		for j := 0; j < columns; j++ {
			cells[i][j] = newCell(i, j)

			cells[i][j].alive = false
			cells[i][j].aliveNext = rand.Float64() < threshold
		}
	}
	return cells
}

//--- opengl init ---

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(width, height, "Game of life", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetPos(100, 100)

	return window
}

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.LinkProgram(prog)
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
