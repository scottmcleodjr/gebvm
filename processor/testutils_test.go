package processor_test

import (
	"bufio"
	"os"
	"testing"

	"github.com/scottmcleodjr/gebvm/memory"
	"github.com/scottmcleodjr/gebvm/processor"
)

const (
	R0 uint8 = 0
	R1 uint8 = 1
	R2 uint8 = 2
	R3 uint8 = 3
	R4 uint8 = 4
	R5 uint8 = 5
	R6 uint8 = 6
	R7 uint8 = 7
)

func newTestProcessorWithPogram(program []uint8) (*processor.Processor, *memory.Memory) {
	m := memory.New()
	m.LoadProgram(program)
	return processor.New(
		m,
		bufio.NewReader(os.Stdin),
		bufio.NewWriter(os.Stdout),
		bufio.NewWriter(os.Stderr),
	), m
}

func stepAndCheckContinueValue(t *testing.T, p *processor.Processor, expected bool) {
	t.Helper()
	actual := p.Step()
	if expected != actual {
		t.Errorf("got %t, want %t", actual, expected)
	}
}

func highByte(value uint16) uint8 {
	return uint8(value >> 8)
}

func lowByte(value uint16) uint8 {
	return uint8(value)
}

func TestHighByte(t *testing.T) {
	tests := []struct {
		input    uint16
		expected uint8
	}{
		{input: 0x1234, expected: 0x12},
		{input: 0xCDEF, expected: 0xCD},
	}
	for _, test := range tests {
		actual := highByte(test.input)
		if actual != test.expected {
			t.Errorf("got 0x%X, want 0x%X", actual, test.expected)
		}
	}
}

func TestLowByte(t *testing.T) {
	tests := []struct {
		input    uint16
		expected uint8
	}{
		{input: 0x1234, expected: 0x34},
		{input: 0xCDEF, expected: 0xEF},
	}
	for _, test := range tests {
		actual := lowByte(test.input)
		if actual != test.expected {
			t.Errorf("got 0x%X, want 0x%X", actual, test.expected)
		}
	}
}
