package main

import (
	"fmt"
	"os"
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

func main() {
	myCPU := NewCpu()
	err := myCPU.load_program_to_memory("roms/test_opcode.ch8")

	if err != nil {
		fmt.Printf("failed to load program: %v\n", err)
		return
	}
	fmt.Printf("%+v\n", myCPU.memory)
}
