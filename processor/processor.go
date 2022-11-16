package processor

import (
	"bufio"
	"errors"
	"fmt"
)

const (
	RegisterCount uint8  = 8
	StackStart    uint16 = 0xFF00
	StackLimit    uint16 = 0xFFFF
)

type MemoryDevice interface {
	WriteMemory(address uint16, value uint8)
	ReadMemory(address uint16) uint8
	LoadProgram(program []uint8) error
}

type Processor struct {
	memory             MemoryDevice
	registers          [RegisterCount]uint8
	instructionPointer uint16
	errors             []error       // errors encountered during execution
	stackPointer       uint16        // absolution position of top of stack in memory
	stackSize          uint8         // size of current stack call frame
	reader             *bufio.Reader // reader for input to RIN
	writer             *bufio.Writer // writer for output from PNT
	errorWriter        *bufio.Writer // writer for execution errors
}

func New(m MemoryDevice, r *bufio.Reader, w, ew *bufio.Writer) *Processor {
	return &Processor{
		memory:       m,
		stackPointer: StackStart,
		reader:       r,
		writer:       w,
		errorWriter:  ew,
	}
}

func (p *Processor) InstructionPointer() uint16 {
	return p.instructionPointer
}

func (p *Processor) StackPointer() uint16 {
	return p.stackPointer
}

func (p *Processor) StackSize() uint8 {
	return p.stackSize
}

func (p *Processor) Errors() []error {
	out := make([]error, len(p.errors))
	copy(out, p.errors)
	return out
}

func (p *Processor) RegisterValue(register uint8) uint8 {
	if register >= RegisterCount {
		p.errors = append(p.errors, fmt.Errorf("invalid register access: %d", register))
		return 0x00
	}
	return p.registers[register]
}

func (p *Processor) setRegisterValue(register uint8, value uint8) {
	if register >= RegisterCount {
		p.errors = append(p.errors, fmt.Errorf("invalid register access: %d", register))
		return
	}
	p.registers[register] = value
}

func (p *Processor) registerPointerValue(register uint8) uint16 {
	highByte := p.RegisterValue(register)
	lowByte := p.RegisterValue(register + 1)
	return (uint16(highByte) << 8) + uint16(lowByte)
}

func (p *Processor) fetchInstruction() uint8 {
	instruction := p.memory.ReadMemory(p.instructionPointer)
	p.instructionPointer++
	// Detect instructionPointer overflow
	if p.instructionPointer == 0x0000 {
		p.errors = append(p.errors, errors.New("instruction pointer out of memory bounds"))
		return 0x00
	}
	return instruction
}

func (p *Processor) fetchAddressInstruction() uint16 {
	highByte := p.fetchInstruction()
	lowByte := p.fetchInstruction()
	return (uint16(highByte) << 8) + uint16(lowByte)
}

func (p *Processor) Step() bool {
	instruction := p.fetchInstruction()

	handler, instructionFound := instructions[instruction]
	if !instructionFound {
		p.errors = append(p.errors,
			fmt.Errorf("unknown instruction 0x%X at position 0x%X", instruction, p.instructionPointer-1))
		return false
	}

	continueRunning := handler(p)
	return continueRunning && len(p.errors) == 0
}

func (p *Processor) Run() int {
	for continueRunning := true; continueRunning; continueRunning = p.Step() {
	}
	if len(p.errors) > 0 {
		fmt.Fprintf(p.errorWriter, "** ERRORS:\n")
		for _, e := range p.errors {
			fmt.Fprintf(p.errorWriter, "** %s\n", e)
			p.errorWriter.Flush()
		}
		return 1
	}
	return 0
}
