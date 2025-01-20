package codegen

import (
	"fmt"
	"strings"

	"example.com/test/ast"
)

const (
	header    = "source_filename = \"main.c\"\ntarget datalayout = \"e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128\"\ntarget triple = \"x86_64-pc-linux-gnu\""
	var_tape  = "@tape = dso_local global [4294967296 x i8] zeroinitializer, align 16"
	var_ptr   = "@pointer = dso_local global i32 0, align 4"
	var_str   = `@.str = private unnamed_addr constant [11 x i8] c"Print: %u\0A\00", align 1`
	var_print = `define dso_local void @print(i8 noundef zeroext %0) #0 {
  %2 = alloca i8, align 1
  store i8 %0, ptr %2, align 1
  %3 = load i8, ptr %2, align 1
  %4 = zext i8 %3 to i32
  %5 = call i32 (ptr, ...) @printf(ptr noundef @.str, i32 noundef %4)
  ret void
}

declare i32 @printf(ptr noundef, ...)`
)

type CodeGenContext struct {
	Instructions []string // Generated code
	LoopCounter  int      // Counter for unique loop labels
	VarCounter   int
}

func (ctx *CodeGenContext) emit(instruction string) {
	ctx.Instructions = append(ctx.Instructions, instruction)
}

func (ctx *CodeGenContext) emitIncrementPointer(amount int) {
	instruction := fmt.Sprintf("%%%d = load i32, ptr @pointer\n%%%d = add i32 %%%d, %d\nstore i32 %%%d, ptr @pointer", ctx.VarCounter, ctx.VarCounter+1, ctx.VarCounter, amount, ctx.VarCounter+1)
	ctx.VarCounter += 2
	ctx.Instructions = append(ctx.Instructions, instruction)
}

func (ctx *CodeGenContext) emitDecrementPointer(amount int) {
	instruction := fmt.Sprintf("%%%d = load i32, ptr @pointer\n%%%d = sub i32 %%%d, %d\nstore i32 %%%d, ptr @pointer", ctx.VarCounter, ctx.VarCounter+1, ctx.VarCounter, amount, ctx.VarCounter+1)
	ctx.VarCounter += 2
	ctx.Instructions = append(ctx.Instructions, instruction)
}

func (ctx *CodeGenContext) emitIncrementValue(amount int) {
	ct := ctx.VarCounter
	instruction := fmt.Sprintf("%%%d = load i32, ptr @pointer\n%%%d = zext i32 %%%d to i64\n%%%d = getelementptr inbounds [0 x i8], ptr @tape, i64 0, i64 %%%d\n%%%d = load i8, ptr %%%d\n%%%d = zext i8 %%%d to i32\n%%%d = add nsw i32 %d, %%%d\n%%%d = trunc i32 %%%d to i8\nstore i8 %%%d, ptr %%%d",
		ct, ct+1, ct, ct+2, ct+1, ct+3, ct+2, ct+4, ct+3, ct+5, amount, ct+4, ct+6, ct+5, ct+6, ct+2)
	ctx.VarCounter += 7
	ctx.Instructions = append(ctx.Instructions, instruction)
}

func (ctx *CodeGenContext) emitDecrementValue(amount int) {
	ct := ctx.VarCounter
	instruction := fmt.Sprintf("%%%d = load i32, ptr @pointer\n%%%d = zext i32 %%%d to i64\n%%%d = getelementptr inbounds [0 x i8], ptr @tape, i64 0, i64 %%%d\n%%%d = load i8, ptr %%%d\n%%%d = zext i8 %%%d to i32\n%%%d = sub nsw i32 %d, %%%d\n%%%d = trunc i32 %%%d to i8\nstore i8 %%%d, ptr %%%d",
		ct, ct+1, ct, ct+2, ct+1, ct+3, ct+2, ct+4, ct+3, ct+5, amount, ct+4, ct+6, ct+5, ct+6, ct+2)
	ctx.VarCounter += 7
	ctx.Instructions = append(ctx.Instructions, instruction)
}

func (ctx *CodeGenContext) emitOut() {
	ct := ctx.VarCounter
	instruction := fmt.Sprintf("%%%d = load i32, ptr @pointer\n%%%d = zext i32 %%%d to i64\n%%%d = getelementptr inbounds [4294967296 x i8], ptr @tape, i64 0, i64 %%%d\n%%%d = load i8, ptr %%%d\ncall void @print(i8 noundef zeroext %%%d)",
		ct, ct+1, ct, ct+2, ct+1, ct+3, ct+2, ct+3)
	ctx.VarCounter += 4
	ctx.Instructions = append(ctx.Instructions, instruction)
}

func GenerateNode(ctx *CodeGenContext, node *ast.ASTNode) {
	switch node.Type {
	case ast.NodeIncrementPointer:
		ctx.emitIncrementPointer(1)
	case ast.NodeDecrementPointer:
		ctx.emitDecrementPointer(1)
	case ast.NodeIncrementValue:
		ctx.emitIncrementValue(1)
	case ast.NodeDecrementValue:
		ctx.emitDecrementValue(1)
	case ast.NodeOutput:
		ctx.emitOut()
	case ast.NodeInput:
		ctx.emit("")
	case ast.NodeLoop:
		// Unique loop labels
		loopStart := ctx.LoopCounter
		loopEnd := ctx.LoopCounter + 1
		ctx.LoopCounter += 2

		// Emit loop start
		ctx.emit(fmt.Sprintf("br label %%loop%d", loopStart))
		ctx.emit(fmt.Sprintf("loop%d:", loopStart))
		ctx.emit("call i1 @is_zero()")
		ctx.emit(fmt.Sprintf("br i1 %%retval, label %%loop%d_end, label %%loop%d_body", loopEnd, loopStart))

		// Emit loop body
		ctx.emit(fmt.Sprintf("loop%d_body:", loopStart))
		for _, child := range node.Children {
			GenerateNode(ctx, child)
		}

		// Loop back
		ctx.emit(fmt.Sprintf("br label %%loop%d", loopStart))

		// Emit loop end
		ctx.emit(fmt.Sprintf("loop%d_end:", loopEnd))
	}
}

func GenerateCode(ast []*ast.ASTNode) string {
	ctx := &CodeGenContext{Instructions: []string{}, LoopCounter: 0, VarCounter: 1}
	ctx.Instructions = append(ctx.Instructions, header)
	ctx.Instructions = append(ctx.Instructions, var_ptr)
	ctx.Instructions = append(ctx.Instructions, var_str)
	ctx.Instructions = append(ctx.Instructions, var_tape)
	ctx.Instructions = append(ctx.Instructions, var_print)
	ctx.Instructions = append(ctx.Instructions, "define i32 @main() {")

	// Process each top-level node
	for _, node := range ast {
		GenerateNode(ctx, node)
	}

	ctx.Instructions = append(ctx.Instructions, "ret i32 0\n}")
	// Combine all instructions into a single LLVM IR string
	return strings.Join(ctx.Instructions, "\n")
}
