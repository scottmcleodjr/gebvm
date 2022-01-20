# GebVM

GebvM is an 8-bit register-based virtual machine.  GebVM has...

**Registers** There are 8 8-bit registers (R0-R7) that store unsigned values.

**Memory** There are 65,536 bytes of memory also storing unsigned values.  The running program is loaded into memory starting at address 0x0000.  A stack occupies the space from 0xFF00 to 0xFFFF.

**Bytecode** GebVM uses its own bytecode specification.  A couple trivial compiled examples are included, and more interesting examples may be added if there is an assembler for the project in the future.

**I/O** Input, output, and errors are written from/to stdin, stdout, and stderr (respectively).

## Running the "Hello World" Program

```
go get github.com/scottmcleodjr/gebvm
go build github.com/scottmcleodjr/gebvm
./gebvm $GOPATH/src/github.com/scottmcleodjr/gebvm/examples/hello_world.geb
```

## Instructions

| Instruction         | Code | Arguments                               | Description                                                       |
|---------------------|------|-----------------------------------------|-------------------------------------------------------------------|
| No Operation        | 0x00 |                                         | Do nothing and continue execution                                 |
| Move Lit Reg        | 0x01 | Literal, Register                       | Copy literal value to register                                    |
| Move Reg Reg        | 0x02 | Source Register, Destination Register   | Copy value from source to destination register                    |
| Move Lit Memory     | 0x03 | Literal, Pointer Register               | Copy literal value to address at pointer register                 |
| Move Reg Memory     | 0x04 | Register, Pointer Register              | Copy value from register to address at pointer register           |
| Move Memory Reg     | 0x05 | Pointer Register, Register              | Copy value from address at pointer register to register           |
| Logical And         | 0x20 | Left Register, Right Register           | Set R0 to logical and of values in left and right registers       |
| Logical Or          | 0x21 | Left Register, Right Register           | Set R0 to logical or of values in left and right registers        |
| Logical Xor         | 0x22 | Left Register, Right Register           | Set R0 to logical xor of values in left and right registers       |
| Logical Bit Clear   | 0x23 | Left Register, Right Register           | Set R0 to logical bit clear of values in left and right registers |
| Logical Shift Left  | 0x24 | Register, Shift Distance                | Logical shift left value in register by distance bytes            |
| Logical Shift Right | 0x25 | Register, Shift Distance                | Logical shift right value in register by distance bytes           |
| Inc                 | 0x40 | Register                                | Increment value in register by 1                                  |
| Dec                 | 0x41 | Register                                | Decrement value in register by 1                                  |
| Add                 | 0x42 | Left Register, Right Register           | Set R0 to sum of values in left and right registers               |
| Subtract            | 0x43 | Left Register, Right Register           | Set R0 to difference of values in left and right registers        |
| Multiply            | 0x44 | Left Register, Right Register           | Set R0 to product of values in left and right registers           |
| Divide              | 0x45 | Left Register, Right Register           | Set R0 to quotient of values in left and right registers          |
| Jump                | 0x60 | Address (High Byte, Low Byte)           | Set IP to address                                                 |
| JumpEqual           | 0x61 | Register, Address (High Byte, Low Byte) | Set IP to address if values in R0 and register are equal          |
| JumpNotEqual        | 0x62 | Register, Address (High Byte, Low Byte) | Set IP to address if values in R0 and register are not equal      |
| StackPushLit        | 0x80 | Literal                                 | Push literal value onto the stack                                 |
| StackPushReg        | 0x81 | Register                                | Push value from register onto the stack                           |
| StackPop            | 0x82 | Register                                | Pop the top value from the stack and store at register            |
| Call                | 0x83 | Address (High Byte, Low Byte)           | Function call                                                     |
| Return              | 0x84 |                                         | Function return                                                   |
| Print               | 0xE0 | Address (High Byte, Low Byte), Length   | Output length chars beginning at address to the writer            |
| ReadInput           | 0xE1 | Register                                | Store a single char from the reader to register                   |
| Halt                | 0xFF |                                         | Halt execution                                                    |

#### Notes:
- A pointer register holds the high byte of the address.  The sequential next register holds the low byte.  R7 cannot be used as a pointer register.
- Registers are represented by the numbers 0x00 (R0) through 0x07 (R7).
