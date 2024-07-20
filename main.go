package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type cpu struct {
	opcode      uint16
	memory      [0xFFF]uint8
	v_register  uint8
	I           uint16
	pc          uint16
	gfx         [64 * 32]uint8
	delay_timer uint8
	sound_timer uint8
	stack       [16]uint16
	sp          uint16
	key         [16]uint8
	font_set    [80]uint8
}

func NewCpu() *cpu {
	c := &cpu{
		opcode:      0,
		memory:      [0xFFF]uint8{},
		v_register:  0,
		I:           0,
		pc:          0x200,
		gfx:         [64 * 32]uint8{},
		delay_timer: 0,
		sound_timer: 0,
		stack:       [16]uint16{},
		sp:          0,
		key:         [16]uint8{},
		font_set: [80]uint8{
			0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
			0x20, 0x60, 0x20, 0x20, 0x70, // 1
			0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
			0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
			0x90, 0x90, 0xF0, 0x10, 0x10, // 4
			0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
			0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
			0xF0, 0x10, 0x20, 0x40, 0x40, // 7
			0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
			0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
			0xF0, 0x90, 0xF0, 0x90, 0x90, // A
			0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
			0xF0, 0x80, 0x80, 0x80, 0xF0, // C
			0xE0, 0x90, 0x90, 0x90, 0xE0, // D
			0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
			0xF0, 0x80, 0xF0, 0x80, 0x80, // F
		},
	}

	for i, byteValue := range c.font_set {
		c.memory[0x50+i] = byteValue
	}

	return c
}

func (c *cpu) load_program_to_memory(filePath string) error {
	content, err := os.ReadFile(filePath)

	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	startProgramAddress := 0x200
	for i, byteValue := range content {

		if startProgramAddress+i < len(c.memory) {
			c.memory[startProgramAddress+i] = byteValue
		} else {
			return fmt.Errorf("program size exceeds available memory")
		}
	}

	return nil
}

type app struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	surface  *sdl.Surface
	cpu      *cpu
}

func NewApp() *app {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	const (
		TITLE         = "CPU-8 GO PROJECT"
		SCREEN_WIDTH  = 64
		SCREEN_HEIGHT = 32
		SCREEN_SCALE  = 10
	)

	window, err := sdl.CreateWindow(
		TITLE,
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		SCREEN_WIDTH*SCREEN_SCALE,
		SCREEN_HEIGHT*SCREEN_SCALE,
		sdl.WINDOW_SHOWN,
	)

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE|sdl.RENDERER_PRESENTVSYNC)

	if err != nil {
		panic(err)
	}

	surface, err := window.GetSurface()

	if err != nil {
		panic(err)
	}

	cpu := NewCpu()
	err = cpu.load_program_to_memory("roms/test_opcode.ch8")

	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", cpu.memory)

	return &app{
		window,
		renderer,
		surface,
		cpu,
	}
}

func (a *app) main_loop() {
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}

		sdl.Delay(33)
	}
}

func main() {
	app := NewApp()

	app.surface.FillRect(nil, 0)
	app.main_loop()

	defer app.window.Destroy()
	defer sdl.Quit()

}
