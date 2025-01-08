package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func main() {

}

func printTree(w io.Writer, text string) {
	nodes := indentTree(text)

	var printNodes func([]*node, string)
	printNodes = func(nodes []*node, prefix string) {
		for i := 0; i < len(nodes); i++ {
			n := nodes[i]

			lastChild := i == len(nodes)-1
			hasChildren := len(n.children) > 0

			childBar := "│" // give us the ability to remove the bar
			var runes string

			if lastChild {
				runes = "└──"
				if hasChildren {
					childBar = " "
				}
			} else {
				runes = "├──"
			}

			fmt.Fprintf(w, "%s%s%s\n", prefix, runes+" ", n.text) // print node
			printNodes(n.children, prefix+childBar+"   ")         // print children

		}
	}
	printNodes(nodes, "")
}

type stack []*node

func (s *stack) push(n *node) {
	*s = append(*s, n)
}

func (s *stack) pop() *node {
	n := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return n
}

func (s *stack) replaceTop(n *node) {
	(*s)[len(*s)-1] = n
}

func (s *stack) top() *node {
	return (*s)[len(*s)-1]
}

func (s *stack) root() *node {
	return (*s)[0]
}

// indentTree builds up a node tree from whitespace indented text
func indentTree(text string) []*node {

	st := stack{}
	st.push(&node{children: []*node{}, indent: -1})

	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		nn := newNode(scanner.Text())
		n := &nn

		if n.indent > st.top().indent {
			st.top().children = append(st.top().children, n)
			st.push(n)
		} else if n.indent == st.top().indent {
			st.pop()
			st.top().children = append(st.top().children, n)
			st.push(n)
		} else {
			for n.indent <= st.top().indent {
				st.pop()
			}
			st.top().children = append(st.top().children, n)
			st.push(n)
		}

	}

	return st.root().children
}

func newNode(text string) node {
	indent := indentDepth(text, '	')
	return node{text: text[indent:], indent: indent}
}

type node struct {
	text     string
	children []*node
	indent   int
}

// indentDepth returns the number of indent runes at the beginning of the line. this is like the stdlib
func indentDepth(line string, indent rune) int {
	depth := 0
	for _, r := range line {
		if r == indent {
			depth++
		} else {
			break
		}
	}
	return depth
}
