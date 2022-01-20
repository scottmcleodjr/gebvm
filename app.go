package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/scottmcleodjr/gebvm/memory"
	"github.com/scottmcleodjr/gebvm/processor"
)

const (
	helpText string = `Provide the bytecode file to execute as an argument.
	
  Example: ./gebvm hello_world.geb

`
)

func main() {
	if len(os.Args) == 1 {
		fmt.Print(helpText)
		os.Exit(0)
	}

	filename := os.Args[1]
	program, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input file: %s", err)
		os.Exit(1)
	}

	m := memory.Memory{}
	err = m.LoadProgram(program)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	proc := processor.New(
		&m,
		bufio.NewReader(os.Stdin),
		bufio.NewWriter(os.Stdout),
		bufio.NewWriter(os.Stderr),
	)
	status := proc.Run()
	os.Exit(status)
}
