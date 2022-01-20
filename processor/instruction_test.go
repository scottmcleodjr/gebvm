package processor_test

import (
	"testing"

	"github.com/scottmcleodjr/gebvm/processor"
)

func TestExecuteNoop(t *testing.T) {
	p, _ := newTestProcessorWithPogram([]uint8{processor.Noop})
	stepAndCheckContinueValue(t, p, true)
}

func TestExecuteMoveLitReg(t *testing.T) {
	literals := []uint8{1, 13, 42, 137}
	for register := R0; register < processor.RegisterCount; register++ {
		for _, literal := range literals {
			p, _ := newTestProcessorWithPogram([]uint8{processor.MoveLitReg, literal, register})
			stepAndCheckContinueValue(t, p, true)
			if p.RegisterValue(register) != literal {
				t.Errorf("got 0x%X at R%d, want 0x%X", p.RegisterValue(register), register, literal)
			}
		}
	}
}

func TestExecuteMoveRegReg(t *testing.T) {
	for src := R0; src < processor.RegisterCount; src++ {
		for dst := R0; dst < processor.RegisterCount; dst++ {
			p, _ := newTestProcessorWithPogram([]uint8{
				processor.MoveLitReg, 0x42, src,
				processor.MoveRegReg, src, dst,
			})
			p.Step()
			stepAndCheckContinueValue(t, p, true)
			if p.RegisterValue(dst) != 0x42 {
				t.Errorf("got 0x%X at R%d, want 0x42", p.RegisterValue(dst), dst)
			}
		}
	}
}

func TestExecuteMoveLitMem(t *testing.T) {
	tests := []struct {
		literal uint8
		address uint16
	}{
		{literal: 0x42, address: 0x0000},
		{literal: 0xAB, address: 0x1234},
		{literal: 0xCD, address: 0xFFFF},
	}

	for _, test := range tests {
		for register := R0; register < processor.RegisterCount-1; register++ {
			p, m := newTestProcessorWithPogram([]uint8{
				processor.MoveLitReg, highByte(test.address), register,
				processor.MoveLitReg, lowByte(test.address), register + 1,
				processor.MoveLitMem, test.literal, register,
			})
			p.Step()
			p.Step()
			stepAndCheckContinueValue(t, p, true)
			if m.ReadMemory(test.address) != test.literal {
				t.Errorf("got 0x%X at address 0x%X, want 0x%X",
					m.ReadMemory(test.address), test.address, test.literal)
			}
		}
	}
}

func TestExecuteMoveRegMem(t *testing.T) {
	addresses := []uint16{0x0000, 0x1234, 0xFFFF}
	for _, address := range addresses {
		p, m := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, 0x42, R1,
			processor.MoveLitReg, highByte(address), R2,
			processor.MoveLitReg, lowByte(address), R3,
			processor.MoveRegMem, R1, R2,
		})
		p.Step()
		p.Step()
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if m.ReadMemory(address) != 0x42 {
			t.Errorf("got 0x%X at address 0x%X, want 0x42", m.ReadMemory(address), address)
		}
	}
}

func TestExecuteMoveMemReg(t *testing.T) {
	addresses := []uint16{0xABCD, 0x1234, 0xFFFF}
	for _, address := range addresses {
		p, m := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, highByte(address), R1,
			processor.MoveLitReg, lowByte(address), R2,
			processor.MoveMemReg, R1, R3})
		m.WriteMemory(address, 0x42)
		p.Step()
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if p.RegisterValue(R3) != 0x42 {
			t.Errorf("got 0x%X at R3, want 0x42", p.RegisterValue(R3))
		}
	}
}

func TestExecuteLogicalAnd(t *testing.T) {
	tests := []struct {
		inputLeft, inputRight, expected uint8
	}{
		{inputLeft: 0x5C, inputRight: 0xD2, expected: 0x50},
		{inputLeft: 0xD6, inputRight: 0x13, expected: 0x12},
		{inputLeft: 0xCE, inputRight: 0x47, expected: 0x46},
	}

	for _, test := range tests {
		p, _ := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, test.inputLeft, R1,
			processor.MoveLitReg, test.inputRight, R2,
			processor.LogicalAnd, R1, R2,
		})
		p.Step()
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if p.RegisterValue(0) != test.expected {
			t.Errorf("got 0x%X at R0, want 0x%X", p.RegisterValue(0), test.expected)
		}
	}
}

func TestExecuteLogicalOr(t *testing.T) {
	tests := []struct {
		inputLeft, inputRight, expected uint8
	}{
		{inputLeft: 0x5C, inputRight: 0xD2, expected: 0xDE},
		{inputLeft: 0xD6, inputRight: 0x13, expected: 0xD7},
		{inputLeft: 0xCE, inputRight: 0x47, expected: 0xCF},
	}

	for _, test := range tests {
		p, _ := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, test.inputLeft, R1,
			processor.MoveLitReg, test.inputRight, R2,
			processor.LogicalOr, R1, R2,
		})
		p.Step()
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if p.RegisterValue(0) != test.expected {
			t.Errorf("got 0x%X at R0, want 0x%X", p.RegisterValue(0), test.expected)
		}
	}
}

func TestExecuteLogicalXor(t *testing.T) {
	tests := []struct {
		inputLeft, inputRight, expected uint8
	}{
		{inputLeft: 0x5C, inputRight: 0xD2, expected: 0x8E},
		{inputLeft: 0xD6, inputRight: 0x13, expected: 0xC5},
		{inputLeft: 0xCE, inputRight: 0x47, expected: 0x89},
	}

	for _, test := range tests {
		p, _ := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, test.inputLeft, R1,
			processor.MoveLitReg, test.inputRight, R2,
			processor.LogicalXor, R1, R2,
		})
		p.Step()
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if p.RegisterValue(0) != test.expected {
			t.Errorf("got 0x%X at R0, want 0x%X", p.RegisterValue(0), test.expected)
		}
	}
}

func TestExecuteLogicalBitClear(t *testing.T) {
	tests := []struct {
		inputLeft, inputRight, expected uint8
	}{
		{inputLeft: 0x5C, inputRight: 0xD2, expected: 0x0C},
		{inputLeft: 0xD6, inputRight: 0x13, expected: 0xC4},
		{inputLeft: 0xCE, inputRight: 0x47, expected: 0x88},
	}

	for _, test := range tests {
		p, _ := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, test.inputLeft, R1,
			processor.MoveLitReg, test.inputRight, R2,
			processor.LogicalBitClear, R1, R2,
		})
		p.Step()
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if p.RegisterValue(0) != test.expected {
			t.Errorf("got 0x%X at R0, want 0x%X", p.RegisterValue(0), test.expected)
		}
	}
}

func TestExecuteLogicalShiftLeft(t *testing.T) {
	tests := []struct {
		input, distance, expected uint8
	}{
		{input: 0x55, distance: 1, expected: 0xAA},
		{input: 0xD6, distance: 4, expected: 0x60},
		{input: 0xCE, distance: 8, expected: 0x00},
	}

	for _, test := range tests {
		p, _ := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, test.input, R1,
			processor.LogicalShiftLeft, R1, test.distance,
		})
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if p.RegisterValue(0) != test.expected {
			t.Errorf("got 0x%X at R0, want 0x%X", p.RegisterValue(0), test.expected)
		}
	}
}

func TestExecuteLogicalShiftRight(t *testing.T) {
	tests := []struct {
		input, distance, expected uint8
	}{
		{input: 0x55, distance: 1, expected: 0x2A},
		{input: 0xD6, distance: 4, expected: 0x0D},
		{input: 0xCE, distance: 8, expected: 0x00},
	}

	for _, test := range tests {
		p, _ := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, test.input, R1,
			processor.LogicalShiftRight, R1, test.distance,
		})
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if p.RegisterValue(0) != test.expected {
			t.Errorf("got 0x%X at R0, want 0x%X", p.RegisterValue(0), test.expected)
		}
	}
}

func TestExecuteInc(t *testing.T) {
	values := []uint8{0, 41, 255}
	for _, value := range values {
		p, _ := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, value, R1,
			processor.Inc, R1,
		})
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if p.RegisterValue(1) != value+1 {
			t.Errorf("got 0x%X at R1, want 0x%X", p.RegisterValue(1), value+1)
		}
	}
}

func TestExecuteDec(t *testing.T) {
	values := []uint8{0, 41, 255}
	for _, value := range values {
		p, _ := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, value, R1,
			processor.Dec, R1,
		})
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if p.RegisterValue(1) != value-1 {
			t.Errorf("got 0x%X at R1, want 0x%X", p.RegisterValue(1), value-1)
		}
	}
}

func TestExecuteAdd(t *testing.T) {
	tests := []struct {
		inputLeft  uint8
		inputRight uint8
		expected   uint8
	}{
		{inputLeft: 20, inputRight: 22, expected: 42},
		{inputLeft: 150, inputRight: 150, expected: 44},
		{inputLeft: 0, inputRight: 11, expected: 11},
	}

	for _, test := range tests {
		p, _ := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, test.inputLeft, R1,
			processor.MoveLitReg, test.inputRight, R2,
			processor.Add, R1, R2,
		})
		p.Step()
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if p.RegisterValue(0) != test.expected {
			t.Errorf("got 0x%X at R0, want 0x%X", p.RegisterValue(0), test.expected)
		}
	}
}

func TestExecuteSubtract(t *testing.T) {
	tests := []struct {
		inputLeft  uint8
		inputRight uint8
		expected   uint8
	}{
		{inputLeft: 20, inputRight: 22, expected: 254},
		{inputLeft: 150, inputRight: 150, expected: 0},
		{inputLeft: 11, inputRight: 10, expected: 1},
	}

	for _, test := range tests {
		p, _ := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, test.inputLeft, R1,
			processor.MoveLitReg, test.inputRight, R2,
			processor.Subtract, R1, R2,
		})
		p.Step()
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if p.RegisterValue(0) != test.expected {
			t.Errorf("got 0x%X at R0, want 0x%X", p.RegisterValue(0), test.expected)
		}
	}
}

func TestExecuteMultiply(t *testing.T) {
	tests := []struct {
		inputLeft  uint8
		inputRight uint8
		expected   uint8
	}{
		{inputLeft: 21, inputRight: 2, expected: 42},
		{inputLeft: 150, inputRight: 150, expected: 228},
		{inputLeft: 1, inputRight: 10, expected: 10},
	}

	for _, test := range tests {
		p, _ := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, test.inputLeft, R1,
			processor.MoveLitReg, test.inputRight, R2,
			processor.Multiply, R1, R2,
		})
		p.Step()
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if p.RegisterValue(0) != test.expected {
			t.Errorf("got 0x%X at R0, want 0x%X", p.RegisterValue(0), test.expected)
		}
	}
}

func TestExecuteDivide(t *testing.T) {
	tests := []struct {
		inputLeft  uint8
		inputRight uint8
		expected   uint8
	}{
		{inputLeft: 21, inputRight: 2, expected: 10},
		{inputLeft: 150, inputRight: 150, expected: 1},
		{inputLeft: 1, inputRight: 10, expected: 0},
	}

	for _, test := range tests {
		p, _ := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, test.inputLeft, R1,
			processor.MoveLitReg, test.inputRight, R2,
			processor.Divide, R1, R2,
		})
		p.Step()
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if p.RegisterValue(0) != test.expected {
			t.Errorf("got 0x%X at R0, want 0x%X", p.RegisterValue(0), test.expected)
		}
	}

	p, _ := newTestProcessorWithPogram([]uint8{
		processor.MoveLitReg, 0x01, R1,
		processor.MoveLitReg, 0x00, R2,
		processor.Divide, R1, R2,
	})
	p.Step()
	p.Step()
	stepAndCheckContinueValue(t, p, false)
	if len(p.Errors()) == 0 {
		t.Error("no errors listed after divide by zero")
	}
}

func TestExecuteJump(t *testing.T) {
	addresses := []uint16{0xABCD, 0x1234, 0xFFFF}
	for _, address := range addresses {
		p, _ := newTestProcessorWithPogram([]uint8{processor.Jump, highByte(address), lowByte(address)})
		stepAndCheckContinueValue(t, p, true)
		if p.InstructionPointer() != address {
			t.Errorf("got 0x%X at IP, want 0x%X", p.InstructionPointer(), address)
		}
	}
}

func TestExecuteJumpEqual(t *testing.T) {
	addresses := []uint16{0xABCD, 0x1234, 0xFFFF}

	// When Equal
	for _, address := range addresses {
		p, _ := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, 0x01, R0,
			processor.MoveLitReg, 0x01, R1,
			processor.JumpEqual, R1, highByte(address), lowByte(address),
		})
		p.Step()
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if p.InstructionPointer() != address {
			t.Errorf("got 0x%X at IP, want 0x%X", p.InstructionPointer(), address)
		}
	}

	// When Not Equal
	for _, address := range addresses {
		p, _ := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, 0x00, R0,
			processor.MoveLitReg, 0x01, R1,
			processor.JumpEqual, R1, highByte(address), lowByte(address)})
		p.Step()
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if p.InstructionPointer() == address {
			t.Errorf("executeJumpEqual changed IP when values were not equal.")
		}
	}
}

// Jump Not Equal
func TestExecuteJumpNotEqual(t *testing.T) {
	addresses := []uint16{0xABCD, 0x1234, 0xFFFF}

	// When Not Equal
	for _, address := range addresses {
		p, _ := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, 0x00, R0,
			processor.MoveLitReg, 0x01, R1,
			processor.JumpNotEqual, R1, highByte(address), lowByte(address),
		})
		p.Step()
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if p.InstructionPointer() != address {
			t.Errorf("got 0x%X at IP, want 0x%X", p.InstructionPointer(), address)
		}
	}

	// When Equal
	for _, address := range addresses {
		p, _ := newTestProcessorWithPogram([]uint8{
			processor.MoveLitReg, 0x01, R0,
			processor.MoveLitReg, 0x01, R1,
			processor.JumpNotEqual, R1, highByte(address), lowByte(address),
		})
		p.Step()
		p.Step()
		stepAndCheckContinueValue(t, p, true)
		if p.InstructionPointer() == address {
			t.Errorf("executeJumpNotEqual changed IP when values were equal.")
		}
	}
}

func TestExecuteStackPushLit(t *testing.T) {
	p, m := newTestProcessorWithPogram([]uint8{
		processor.StackPushLit, 0x13,
		processor.StackPushLit, 0x42,
	})
	stepAndCheckContinueValue(t, p, true)
	stepAndCheckContinueValue(t, p, true)
	if m.ReadMemory(0xFF00) != 0x13 {
		t.Errorf("got 0x%X at 0xFF00, want 0x13", m.ReadMemory(0xFF00))
	}
	if m.ReadMemory(0xFF01) != 0x42 {
		t.Errorf("got 0x%X at 0xFF01, want 0x42", m.ReadMemory(0xFF00))
	}
	if p.StackPointer() != 0xFF02 {
		t.Errorf("got 0x%X at SP, want 0xFF02", p.StackPointer())
	}
	if p.StackSize() != 2 {
		t.Errorf("got %d at stack size, want 2", p.StackSize())
	}
}

func TestExecuteStackPushReg(t *testing.T) {
	p, m := newTestProcessorWithPogram([]uint8{
		processor.MoveLitReg, 0x13, R1,
		processor.MoveLitReg, 0x42, R2,
		processor.StackPushReg, R1,
		processor.StackPushReg, R2,
	})
	p.Step()
	p.Step()
	stepAndCheckContinueValue(t, p, true)
	stepAndCheckContinueValue(t, p, true)
	if m.ReadMemory(0xFF00) != 0x13 {
		t.Errorf("got 0x%X at 0xFF00, want 0x13", m.ReadMemory(0xFF00))
	}
	if m.ReadMemory(0xFF01) != 0x42 {
		t.Errorf("got 0x%X at 0xFF01, want 0x42", m.ReadMemory(0xFF00))
	}
	if p.StackPointer() != 0xFF02 {
		t.Errorf("got 0x%X at SP, want 0xFF02", p.StackPointer())
	}
	if p.StackSize() != 2 {
		t.Errorf("got %d at stack size, want 2", p.StackSize())
	}
}

func TestStackOverflow(t *testing.T) {
	program := make([]uint8, 512)
	for x := 0; x < 512; x++ {
		program[x] = processor.StackPushLit
		x++
		program[x] = 0x00
	}
	p, _ := newTestProcessorWithPogram(program)
	for x := 0; x < 255; x++ {
		p.Step()
	}
	stepAndCheckContinueValue(t, p, false)
	if len(p.Errors()) == 0 {
		t.Error("no errors listed after stack overflow")
	}
}

func TestExecuteStackPop(t *testing.T) {
	p, _ := newTestProcessorWithPogram([]uint8{
		processor.StackPushLit, 0x13,
		processor.StackPushLit, 0x42,
		processor.StackPop, R1,
		processor.StackPop, R2,
		processor.StackPop, R3,
	})
	p.Step()
	p.Step()
	// The stack has two items
	stepAndCheckContinueValue(t, p, true)
	stepAndCheckContinueValue(t, p, true)
	if p.RegisterValue(R1) != 0x42 {
		t.Errorf("got 0x%X at R1, want 0x42", p.RegisterValue(R1))
	}
	if p.RegisterValue(R2) != 0x13 {
		t.Errorf("got 0x%X at R2, want 0x13", p.RegisterValue(R2))
	}
	if p.StackSize() != 0 {
		t.Errorf("got %d at stack size, want 0", p.StackSize())
	}
	// The stack is now empty
	stepAndCheckContinueValue(t, p, false)
	if len(p.Errors()) == 0 {
		t.Error("no errors listed after pop from empty stack")
	}

	// Pop from empty stack after Call
	p, _ = newTestProcessorWithPogram([]uint8{
		processor.Call, 0x00, 0x03, // 0x0000->0x0002
		processor.StackPop, R1,
	})
	p.Step()
	stepAndCheckContinueValue(t, p, false)
	if len(p.Errors()) == 0 {
		t.Error("no errors listed after pop from empty stack after call")
	}
}

func TestExecuteCall(t *testing.T) {
	// Only doing some basic checks here.  We'll make sure the
	// correct values come back when we do the return test.
	p, _ := newTestProcessorWithPogram([]uint8{processor.Call, 0xAB, 0xCD})
	stepAndCheckContinueValue(t, p, true)
	if p.StackPointer() != 0xFF09 {
		t.Errorf("got 0x%X at SP, want 0xFF09", p.StackPointer())
	}
	if p.StackSize() != 0 {
		t.Errorf("got %d at stack size, want 0", p.StackSize())
	}
	if p.InstructionPointer() != 0xABCD {
		t.Errorf("got 0x%X at IP, want 0xABCD", p.InstructionPointer())
	}
}

func TestExecuteReturn(t *testing.T) {
	p, _ := newTestProcessorWithPogram([]uint8{

		/**
		 * This is failing when i have the stacksize check in?
		 **/

		processor.MoveLitReg, 0x10, R0, // 0x0000->0x0002
		processor.MoveLitReg, 0x11, R1, // 0x0003->0x0005
		processor.MoveLitReg, 0x12, R2, // 0x0006->0x0008
		processor.Call, 0x00, 0x0D, // 0x0009->0x000B (Jump over the halt)
		processor.Halt, // 0x000C
		processor.StackPushLit, 0x01,
		processor.MoveLitReg, 0x20, R0,
		processor.MoveLitReg, 0x21, R1,
		processor.MoveLitReg, 0x22, R2,
		processor.Return,
	})
	// Stepping through the instructiosn before the return
	for i := 0; i < 8; i++ {
		p.Step()
	}
	stepAndCheckContinueValue(t, p, true) // Checking at the Return
	if p.StackPointer() != 0xFF00 {
		t.Errorf("got 0x%X at SP, want 0xFF00", p.StackPointer())
	}
	if p.StackSize() != 0 {
		t.Errorf("got %d at stack size, want 0", p.StackSize())
	}
	if p.InstructionPointer() != 0x000C {
		t.Errorf("got 0x%X at IP, want 0x000C", p.InstructionPointer())
	}
	if p.RegisterValue(R0) != 0x20 {
		t.Errorf("got 0x%X at R0, want 0x20", p.RegisterValue(R0))
	}
	if p.RegisterValue(R1) != 0x21 {
		t.Errorf("got 0x%X at R1, want 0x21", p.RegisterValue(R1))
	}
	if p.RegisterValue(R2) != 0x12 {
		t.Errorf("got 0x%X at R2, want 0x12", p.RegisterValue(R2))
	}
}

func TestExecuteHalt(t *testing.T) {
	p, _ := newTestProcessorWithPogram([]uint8{processor.Halt})
	stepAndCheckContinueValue(t, p, false)
}
