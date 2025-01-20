package main

import (
	"fmt"
	"os"

	"example.com/test/ast"
	"example.com/test/codegen"
)

// writeToFile writes the given content to a file.
func writeToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func main() {
	code := ".+.+.+.+.+.+.+.+.+.+."
	tokens := ast.Tokenize(code)
	a, err := ast.Parse(tokens)
	if err != nil {
		panic(err)
	}
	//ast.PrintAST(a, 0)

	llvmIR := codegen.GenerateCode(a)
	fmt.Println(llvmIR)

	// Write LLVM IR to a file
	filename := "brainfuck.ll"
	err = writeToFile(filename, llvmIR)
	if err != nil {
		fmt.Printf("Error writing LLVM IR to file: %v\n", err)
	} else {
		fmt.Printf("LLVM IR written to %s\n", filename)
	}

}
