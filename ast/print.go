package ast

import (
	"fmt"
	"strings"

	"github.com/xiam/sexpr/lexer"
	"github.com/xiam/sexpr/node"
)

func Print(n *node.Node) {
	printLevel(n, 0)
}

func printLevel(n *node.Node, level int) {
	if n == nil {
		fmt.Printf(":nil\n")
		return
	}
	indent := strings.Repeat("    ", level)
	fmt.Printf("%s(%s): ", indent, n.Type())
	switch n.Type() {

	case node.NodeTypeExpression, node.NodeTypeList, node.NodeTypeMap:
		fmt.Printf("(%v)\n", n.Token())
		list := n.List()
		for i := range list {
			printLevel(list[i], level+1)
		}

	case node.NodeTypeValue:
		fmt.Printf("%#v (%v)\n", n.Value(), n.Token())

	default:
		panic("unknown node type")
	}
}

func Compile(n *node.Node) []byte {
	return compileNodeLevel(n, 0)
}

func compileNodeLevel(n *node.Node, level int) []byte {
	if n == nil {
		return []byte(":nil")
	}
	switch n.Type() {
	case node.NodeTypeMap:
		nodes := []string{}
		for i := range n.List() {
			nodes = append(nodes, string(compileNodeLevel(n.List()[i], level+1)))
		}
		return []byte(fmt.Sprintf("{%s}", strings.Join(nodes, " ")))

	case node.NodeTypeList:
		nodes := []string{}
		for i := range n.List() {
			nodes = append(nodes, string(compileNodeLevel(n.List()[i], level+1)))
		}
		return []byte(fmt.Sprintf("[%s]", strings.Join(nodes, " ")))

	case node.NodeTypeExpression:
		nodes := []string{}
		for i := range n.List() {
			nodes = append(nodes, string(compileNodeLevel(n.List()[i], level+1)))
		}
		if level == 0 {
			return []byte(fmt.Sprintf("%s", strings.Join(nodes, " ")))
		}
		return []byte(fmt.Sprintf("(%s)", strings.Join(nodes, " ")))

	case node.NodeTypeValue:
		if n.Token().Is(lexer.TokenString) {
			return []byte(fmt.Sprintf("%q", n.Value()))
		}
		return []byte(fmt.Sprintf("%v", n.Value()))

	default:
		panic("unknown node type")
	}
}
