package processor

import (
	"errors"
	"fmt"
)

const (
	Noop              uint8 = 0x00 // NOP
	MoveLitReg        uint8 = 0x01 // MLR
	MoveRegReg        uint8 = 0x02 // MRR
	MoveLitMem        uint8 = 0x03 // MLM
	MoveRegMem        uint8 = 0x04 // MRM
	MoveMemReg        uint8 = 0x05 // MMR
	LogicalAnd        uint8 = 0x20 // LND
	LogicalOr         uint8 = 0x21 // LOR
	LogicalXor        uint8 = 0x22 // LXR
	LogicalBitClear   uint8 = 0x23 // LBC
	LogicalShiftLeft  uint8 = 0x24 // LSL
	LogicalShiftRight uint8 = 0x25 // LSR
	Inc               uint8 = 0x40 // INC
	Dec               uint8 = 0x41 // DEC
	Add               uint8 = 0x42 // ADD
	Subtract          uint8 = 0x43 // SUB
	Multiply          uint8 = 0x44 // MUL
	Divide            uint8 = 0x45 // DIV
	Jump              uint8 = 0x60 // JMP
	JumpEqual         uint8 = 0x61 // JEQ
	JumpNotEqual      uint8 = 0x62 // JNE
	StackPushLit      uint8 = 0x80 // SPL
	StackPushReg      uint8 = 0x81 // SPR
	StackPop          uint8 = 0x82 // STP
	Call              uint8 = 0x83 // CLL
	Return            uint8 = 0x84 // RET
	Print             uint8 = 0xE0 // PNT
	ReadInput         uint8 = 0xE1 // RIN
	Halt              uint8 = 0xFF // HLT
)

// Returns bool indicating if the program should continue running
type instruction func(*Processor) bool

var instructions = map[uint8]instruction{
	Noop:              (*Processor).executeNoop,
	MoveLitReg:        (*Processor).executeMoveLitReg,
	MoveRegReg:        (*Processor).executeMoveRegReg,
	MoveLitMem:        (*Processor).executeMoveLitMem,
	MoveRegMem:        (*Processor).executeMoveRegMem,
	MoveMemReg:        (*Processor).executeMoveMemReg,
	LogicalAnd:        (*Processor).executeLogicalAnd,
	LogicalOr:         (*Processor).executeLogicalOr,
	LogicalXor:        (*Processor).executeLogicalXor,
	LogicalBitClear:   (*Processor).executeLogicalBitClear,
	LogicalShiftLeft:  (*Processor).executeLogicalShiftLeft,
	LogicalShiftRight: (*Processor).executeLogicalShiftRight,
	Inc:               (*Processor).executeInc,
	Dec:               (*Processor).executeDec,
	Add:               (*Processor).executeAdd,
	Subtract:          (*Processor).executeSubtract,
	Multiply:          (*Processor).executeMultiply,
	Divide:            (*Processor).executeDivide,
	Jump:              (*Processor).executeJump,
	JumpEqual:         (*Processor).executeJumpEqual,
	JumpNotEqual:      (*Processor).executeJumpNotEqual,
	StackPushLit:      (*Processor).executeStackPushLit,
	StackPushReg:      (*Processor).executeStackPushReg,
	StackPop:          (*Processor).executeStackPop,
	Call:              (*Processor).executeCall,
	Return:            (*Processor).executeReturn,
	Print:             (*Processor).executePrint,
	ReadInput:         (*Processor).executeReadInput,
	Halt:              (*Processor).executeHalt,
}

func (p *Processor) executeNoop() bool {
	return true
}

/*********
 * MOVES *
 *********/

func (p *Processor) executeMoveLitReg() bool {
	literal := p.fetchInstruction()
	register := p.fetchInstruction()
	p.setRegisterValue(register, literal)
	return true
}

func (p *Processor) executeMoveRegReg() bool {
	srcRegister := p.fetchInstruction()
	dstRegister := p.fetchInstruction()
	p.setRegisterValue(dstRegister, p.RegisterValue(srcRegister))
	return true
}

func (p *Processor) executeMoveLitMem() bool {
	literal := p.fetchInstruction()
	register := p.fetchInstruction()
	address := p.registerPointerValue(register)
	p.memory.Write(address, literal)
	return true
}

func (p *Processor) executeMoveRegMem() bool {
	srcRegister := p.fetchInstruction()
	addressRegister := p.fetchInstruction()
	address := p.registerPointerValue(addressRegister)
	p.memory.Write(address, p.RegisterValue(srcRegister))
	return true
}

func (p *Processor) executeMoveMemReg() bool {
	addressRegister := p.fetchInstruction()
	address := p.registerPointerValue(addressRegister)
	dstRegister := p.fetchInstruction()
	p.setRegisterValue(dstRegister, p.memory.Read(address))
	return true
}

/*********
 * LOGIC *
 *********/

func (p *Processor) executeLogicalAnd() bool {
	registerLeft := p.fetchInstruction()
	registerRight := p.fetchInstruction()
	p.setRegisterValue(0, p.RegisterValue(registerLeft)&p.RegisterValue(registerRight))
	return true
}

func (p *Processor) executeLogicalOr() bool {
	registerLeft := p.fetchInstruction()
	registerRight := p.fetchInstruction()
	p.setRegisterValue(0, p.RegisterValue(registerLeft)|p.RegisterValue(registerRight))
	return true
}

func (p *Processor) executeLogicalXor() bool {
	registerLeft := p.fetchInstruction()
	registerRight := p.fetchInstruction()
	p.setRegisterValue(0, p.RegisterValue(registerLeft)^p.RegisterValue(registerRight))
	return true
}

func (p *Processor) executeLogicalBitClear() bool {
	registerLeft := p.fetchInstruction()
	registerRight := p.fetchInstruction()
	p.setRegisterValue(0, p.RegisterValue(registerLeft)&^p.RegisterValue(registerRight))
	return true
}

func (p *Processor) executeLogicalShiftLeft() bool {
	register := p.fetchInstruction()
	shiftDistance := p.fetchInstruction()
	p.setRegisterValue(0, p.RegisterValue(register)<<shiftDistance)
	return true
}

func (p *Processor) executeLogicalShiftRight() bool {
	register := p.fetchInstruction()
	shiftDistance := p.fetchInstruction()
	p.setRegisterValue(0, p.RegisterValue(register)>>shiftDistance)
	return true
}

/*********
 * MATHS *
 *********/

func (p *Processor) executeInc() bool {
	register := p.fetchInstruction()
	p.setRegisterValue(register, p.RegisterValue(register)+1)
	return true
}

func (p *Processor) executeDec() bool {
	register := p.fetchInstruction()
	p.setRegisterValue(register, p.RegisterValue(register)-1)
	return true
}

func (p *Processor) executeAdd() bool {
	registerLeft := p.fetchInstruction()
	registerRight := p.fetchInstruction()
	sum := p.RegisterValue(registerLeft) + p.RegisterValue(registerRight)
	p.setRegisterValue(0, sum)
	return true
}

func (p *Processor) executeSubtract() bool {
	registerLeft := p.fetchInstruction()
	registerRight := p.fetchInstruction()
	diff := p.RegisterValue(registerLeft) - p.RegisterValue(registerRight)
	p.setRegisterValue(0, diff)
	return true
}

func (p *Processor) executeMultiply() bool {
	registerLeft := p.fetchInstruction()
	registerRight := p.fetchInstruction()
	product := p.RegisterValue(registerLeft) * p.RegisterValue(registerRight)
	p.setRegisterValue(0, product)
	return true
}

func (p *Processor) executeDivide() bool {
	registerLeft := p.fetchInstruction()
	registerRight := p.fetchInstruction()
	if p.RegisterValue(registerRight) == 0x00 {
		p.errors = append(p.errors, errors.New("divide by zero"))
		return false
	}
	quotient := p.RegisterValue(registerLeft) / p.RegisterValue(registerRight)
	p.setRegisterValue(0, quotient)
	return true
}

/*********
 * JUMPS *
 *********/

func (p *Processor) executeJump() bool {
	address := p.fetchAddressInstruction()
	p.instructionPointer = address
	return true
}

func (p *Processor) executeJumpEqual() bool {
	register := p.fetchInstruction()
	address := p.fetchAddressInstruction()
	if p.RegisterValue(0) == p.RegisterValue(register) {
		p.instructionPointer = address
	}
	return true
}

func (p *Processor) executeJumpNotEqual() bool {
	register := p.fetchInstruction()
	address := p.fetchAddressInstruction()
	if p.RegisterValue(0) != p.RegisterValue(register) {
		p.instructionPointer = address
	}
	return true
}

/*********
 * STACK *
 *********/

func (p *Processor) stackPush(value uint8) {
	if p.stackPointer == StackLimit {
		p.errors = append(p.errors, errors.New("stack overflow"))
		return
	}
	p.memory.Write(p.stackPointer, value)
	p.stackPointer++
	p.stackSize++
}

func (p *Processor) stackPop() uint8 {
	p.stackPointer--
	p.stackSize--
	return p.memory.Read(p.stackPointer)
}

func (p *Processor) executeStackPushLit() bool {
	value := p.fetchInstruction()
	p.stackPush(value)
	return true
}

func (p *Processor) executeStackPushReg() bool {
	register := p.fetchInstruction()
	p.stackPush(p.RegisterValue(register))
	return true
}
func (p *Processor) executeStackPop() bool {
	/* Checking this here instead of in stackPop
	 * Because executeReturn pops from stackSize==0
	 * when it fetches the IP and R values */
	if p.stackSize == 0 || p.stackPointer == StackStart {
		p.errors = append(p.errors, errors.New("stack underflow"))
		return false
	}
	register := p.fetchInstruction()
	p.setRegisterValue(register, p.stackPop())
	return true
}

func (p *Processor) executeCall() bool {
	address := p.fetchAddressInstruction()
	p.stackPush(p.stackSize)
	// Only push R2->R7 for return
	for r := uint8(2); r < 8; r++ {
		p.stackPush(p.RegisterValue(r))
	}
	p.stackPush(uint8(p.instructionPointer >> 8))
	p.stackPush(uint8(p.instructionPointer))
	p.stackSize = 0
	p.instructionPointer = address
	return true
}

func (p *Processor) executeReturn() bool {
	for i := uint8(0); i < p.stackSize; i++ {
		p.stackPop() // Current stack falls out of scope
	}
	ip := uint16(p.stackPop())
	ip += uint16(p.stackPop()) << 8
	p.instructionPointer = ip
	for r := uint8(7); r > 1; r-- {
		p.setRegisterValue(r, p.stackPop())
	}
	p.stackSize = p.stackPop()
	return true
}

/*********
 * OTHER *
 *********/

func (p *Processor) executePrint() bool {
	address := p.fetchAddressInstruction()
	length := p.fetchInstruction()
	for i := uint16(0); i < uint16(length); i++ {
		fmt.Fprintf(p.writer, "%c", p.memory.Read(address+i))
	}
	p.writer.Flush()
	return true
}

func (p *Processor) executeReadInput() bool {
	register := p.fetchInstruction()
	c, err := p.reader.ReadByte()
	if err != nil {
		p.errors = append(p.errors, fmt.Errorf("error reading input: %s", err))
	}
	p.setRegisterValue(register, c)
	return true
}

func (p *Processor) executeHalt() bool {
	return false
}
