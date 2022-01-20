package processor_test

import (
	"errors"
	"testing"

	"github.com/scottmcleodjr/gebvm/processor"
)

func TestErrorsReturnsCopy(t *testing.T) {
	p, _ := newTestProcessorWithPogram([]uint8{0x0F, processor.Halt}) // Bad Instruction
	p.Step()
	errs := p.Errors()
	errs[0] = errors.New("a different error")
	if errs[0] == p.Errors()[0] {
		t.Error("processor.Errors() did not return copy of errors slice")
	}
}

func TestUnknownRegisterError(t *testing.T) {
	badRegisters := []uint8{9, 11, 42, 70}
	for _, reg := range badRegisters {

		// Testing with a Move Lit->Reg instruction
		p, _ := newTestProcessorWithPogram([]uint8{processor.MoveLitReg, 0x42, reg})
		stepAndCheckContinueValue(t, p, false)
		if len(p.Errors()) == 0 {
			t.Error("no errors listed after setting value to invalid register.")
		}

		// Testing with a Move Reg->Mem instruction
		p, _ = newTestProcessorWithPogram([]uint8{processor.MoveRegMem, reg, 0x12, 0x34})
		stepAndCheckContinueValue(t, p, false)
		if len(p.Errors()) == 0 {
			t.Error("no errors listed after reading from invalid register.")
		}
	}
}

func TestInstructionPointerOutOfBoundsError(t *testing.T) {
	p, _ := newTestProcessorWithPogram([]uint8{processor.Jump, 0xFF, 0xFF})
	stepAndCheckContinueValue(t, p, true)
	stepAndCheckContinueValue(t, p, false)
	if len(p.Errors()) == 0 {
		t.Error("no errors listed after IP moved out of bounds.")
	}
}

func TestUnknownInstructionError(t *testing.T) {
	bad_instructions := []uint8{0x0F, 0x33, 0x57, 0xCD}
	for _, instruction := range bad_instructions {
		p, _ := newTestProcessorWithPogram([]uint8{instruction})
		stepAndCheckContinueValue(t, p, false)
		if len(p.Errors()) == 0 {
			t.Error("no errors listed after invalid instruction.")
		}
	}
}

func TestRun(t *testing.T) {
	p, _ := newTestProcessorWithPogram([]uint8{processor.Noop, processor.Noop, processor.Halt})
	status := p.Run()
	if status != 0 {
		t.Errorf("got status %d, want 0", status)
	}

	p, _ = newTestProcessorWithPogram([]uint8{0x0F, processor.Halt}) // Bad Instruction
	status = p.Run()
	if status != 1 {
		t.Errorf("got status %d, want 1", status)
	}
}
