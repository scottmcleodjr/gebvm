package memory_test

import (
	"testing"

	"github.com/scottmcleodjr/gebvm/memory"
)

func TestWriteAndReadMemory(t *testing.T) {
	tests := []struct {
		address uint16
		value   uint8
	}{
		{address: 0x0000, value: 0x42},
		{address: 0x1234, value: 0xAB},
		{address: 0xFFFF, value: 0xFF},
	}

	for _, test := range tests {
		m := memory.New()
		m.Write(test.address, test.value)
		if m.Read(test.address) != test.value {
			t.Errorf("got 0x%X at address 0x%X, want 0x%X",
				m.Read(test.address), test.address, test.value)
		}
	}
}

func TestLoadProgram(t *testing.T) {
	program := []uint8{0x01, 0xA2, 0xB3, 0xC4}
	m := memory.New()
	m.LoadProgram(program)
	for i, instruction := range program {
		if m.Read(uint16(i)) != instruction {
			t.Errorf("got 0x%X at address 0x%X, want 0x%X", m.Read(uint16(i)), i, instruction)
		}
	}

	program = make([]uint8, memory.MemorySize+1)
	err := memory.New().LoadProgram(program)
	if err == nil {
		t.Error("got nil, want error for program length exceeding available memory")
	}
}
