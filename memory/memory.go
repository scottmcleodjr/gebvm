package memory

import (
	"errors"
	"math"
)

const MemorySize = math.MaxUint16 + 1 // Limit for 16-bit addresses

type Memory struct {
	memory [MemorySize]uint8
}

func New() *Memory {
	return &Memory{}
}

func (m *Memory) Write(address uint16, value uint8) {
	m.memory[address] = value
}

func (m *Memory) Read(address uint16) uint8 {
	return m.memory[address]
}

func (m *Memory) LoadProgram(program []uint8) error {
	if len(program) > MemorySize {
		return errors.New("program length exceeds available memory")
	}
	for i, instruction := range program {
		m.Write(uint16(i), instruction)
	}
	return nil
}
