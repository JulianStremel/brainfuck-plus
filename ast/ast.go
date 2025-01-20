package ast

const (
	NodeIncrementPointer NodeType = iota
	NodeDecrementPointer
	NodeIncrementValue
	NodeDecrementValue
	NodeOutput
	NodeInput
	NodeLoop
	NodeNumeric
	NodeCharacter
	NodeReturn
)

// NodeType represents the type of an AST node.
type NodeType int

// ASTNode represents a single node in the Brainfuck AST.
type ASTNode struct {
	Type     NodeType
	Children []*ASTNode // For loops, to hold nested instructions
	Data     string
}

// Tokenizing the raw string input
func Tokenize(code string) []rune {
	validTokens := map[rune]bool{
		'>': true, '<': true,
		'+': true, '-': true,
		'.': true, ',': true,
		'[': true, ']': true,
		// future features
		'=': false,
		'(': false, ')': false,
		'!': false, '@': false,
	}

	tokens := []rune{}
	for _, char := range code {
		if validTokens[char] {
			tokens = append(tokens, char)
		} else {
			ascii := int(char)
			// a-z + A-Z + '_' + 0-9
			if (123 > ascii && ascii > 64 && (ascii > 94 || ascii < 91) && ascii != 96) || ascii > 47 && ascii < 58 {
				tokens = append(tokens, char)
			}
		}
	}
	return tokens
}

// Parsing the tokenized input
func Parse(tokens []rune) ([]*ASTNode, error) {
	var parseRecursive func(start int) ([]*ASTNode, int, error)
	parseRecursive = func(start int) ([]*ASTNode, int, error) {
		var nodes []*ASTNode
		for i := start; i < len(tokens); i++ {
			switch tokens[i] {
			case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
				nodes = append(nodes, &ASTNode{Type: NodeCharacter, Data: string(tokens[i])})
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				nodes = append(nodes, &ASTNode{Type: NodeNumeric, Data: string(tokens[i])})
			case '>':
				nodes = append(nodes, &ASTNode{Type: NodeIncrementPointer})
			case '<':
				nodes = append(nodes, &ASTNode{Type: NodeDecrementPointer})
			case '+':
				nodes = append(nodes, &ASTNode{Type: NodeIncrementValue})
			case '-':
				nodes = append(nodes, &ASTNode{Type: NodeDecrementValue})
			case '.':
				nodes = append(nodes, &ASTNode{Type: NodeOutput})
			case ',':
				nodes = append(nodes, &ASTNode{Type: NodeInput})
			case '[':
				// Parse nested loop
				childNodes, endIndex, err := parseRecursive(i + 1)
				if err != nil {
					return nil, 0, err
				}
				nodes = append(nodes, &ASTNode{Type: NodeLoop, Children: childNodes})
				i = endIndex
			case ']':
				// End of loop
				return nodes, i, nil
			default:
				panic("unknown value")
			}
		}
		return nodes, len(tokens), nil
	}

	nodes, _, err := parseRecursive(0)
	return nodes, err
}

// Helper to print the AST
func PrintAST(nodes []*ASTNode, depth int) {
	for _, node := range nodes {
		for i := 0; i < depth; i++ {
			print("  ")
		}
		switch node.Type {
		case NodeNumeric, NodeCharacter:
			println(string(node.Data))
		case NodeIncrementPointer:
			println("IncrementPointer")
		case NodeDecrementPointer:
			println("DecrementPointer")
		case NodeIncrementValue:
			println("IncrementValue")
		case NodeDecrementValue:
			println("DecrementValue")
		case NodeOutput:
			println("Output")
		case NodeInput:
			println("Input")
		case NodeLoop:
			println("Loop {")
			PrintAST(node.Children, depth+1)
			for i := 0; i < depth; i++ {
				print("  ")
			}
			println("}")
		}
	}
}
