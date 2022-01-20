package memory

import (
	"errors"
	"math"
)

const MemorySize = math.MaxUint16 + 1 // Limit for 16-bit addresses

type MemoryDevice interface {
	WriteMemory(address uint16, value uint8)
	ReadMemory(address uint16) uint8
	LoadProgram(program []uint8) error
}

type Memory struct {
	memory [MemorySize]uint8
}

func New() *Memory {
	return &Memory{}
}

func (m *Memory) WriteMemory(address uint16, value uint8) {
	m.memory[address] = value
}

func (m *Memory) ReadMemory(address uint16) uint8 {
	return m.memory[address]
}

func (m *Memory) LoadProgram(program []uint8) error {
	if len(program) > MemorySize {
		return errors.New("program length exceeds available memory")
	}
	for i, instruction := range program {
		m.WriteMemory(uint16(i), instruction)
	}
	return nil
}
