package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	defaultOpen   bool
	defaultFormat string
)

func main() {

	flag.BoolVar(&defaultOpen, "open", false, "open all nodes by default")
	flag.StringVar(&defaultFormat, "format", "tree", "format to use")
	flag.Parse()

	// read from standard input
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read from standard input: %v", err)
	}

	text := string(data)

	p := tea.NewProgram(initialModel(text))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	tree         []*node
	visibleNodes []*node
	cursor       int
}

func (m *model) VisibleNode(index int) *node {
	if index < 0 || index > len(m.visibleNodes)-1 {
		return nil
	}

	return m.visibleNodes[index]
}

func (m *model) NumVisibleNodes() int {
	return len(m.visibleNodes)
}

func initialModel(text string) *model {
	return &model{
		tree:         indentTree(text),
		visibleNodes: []*node{},
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Grocery List")
}

// setup a logger looking at a file
var (
	logFile, _ = os.OpenFile("./logfile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	logger     = log.New(logFile, "", log.LstdFlags)
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.tree = nil
			return m, tea.Quit

		case "k", "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
			}

		case "j", "down", "ctrl+n":
			if m.cursor < m.NumVisibleNodes()-1 {
				m.cursor++
			}

		case "h", "left":
			node := m.VisibleNode(m.cursor)
			if node == nil {
				break
			}
			logger.Printf("par: %v", node.parent.text)

			if node.open && len(node.children) > 0 {
				node.open = false
			} else if node.parent != nil {
				var i int
				for i = 0; i < len(node.parent.children); i++ {
					if node.parent.children[i] == node {
						break
					}
				}
				m.cursor -= i + 1
				node.parent.open = !node.parent.open
			}

		case "tab":
			node := m.VisibleNode(m.cursor)
			if node == nil || len(node.children) == 0 {
				break
			}
			node.open = !node.open

		case "l", "right":
			node := m.VisibleNode(m.cursor)
			if node == nil {
				break
			}
			node.open = len(node.children) > 0
		}
	}

	return m, nil
}

func (m *model) View() string {

	if m.tree == nil {
		return ""
	}

	switch defaultFormat {
	case "tree":
		return treeView(m)
	case "folder":
		return folderView(m)
	}

	return "No style found"
}

// treeView returns a string representation of the tree for rendering
func treeView(m *model) string {
	w := &strings.Builder{}

	lines := 0
	visibleNodes := m.visibleNodes[:0]

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

			var text string
			if lines == m.cursor {
				text = colorGreen(n.text)
			} else {
				text = colorWhite(n.text)
			}

			count := ""
			if !n.open && len(n.children) > 0 {
				count = fmt.Sprintf(" (%v)", len(n.children))
			}

			fmt.Fprintf(w, "%s%s%s%v\n", prefix, runes+" ", text, count) // print node
			lines++
			visibleNodes = append(visibleNodes, n)

			if n.open {
				printNodes(n.children, prefix+childBar+"   ") // print children
			}

		}
	}
	printNodes(m.tree, "")
	m.visibleNodes = visibleNodes

	return w.String()
}

// folderView returns a string representation of a typical folder structure
func folderView(m *model) string {
	w := &strings.Builder{}

	lines := 0
	visibleNodes := m.visibleNodes[:0]

	var printNodes func([]*node, string)
	printNodes = func(nodes []*node, prefix string) {
		for i := 0; i < len(nodes); i++ {
			n := nodes[i]
			hasChildren := len(n.children) > 0

			childBar := "" // give us the ability to remove the bar
			var runes string

			if n.open {
				runes = "▾"
			} else if hasChildren {
				runes = "▸"
			} else {
				runes = " "
			}

			var text string
			if lines == m.cursor {
				text = colorGreen(n.text)
			} else {
				text = colorWhite(n.text)
			}

			count := ""
			if !n.open && len(n.children) > 0 {
				count = fmt.Sprintf(" (%v)", len(n.children))
			}

			fmt.Fprintf(w, "%s%s%s%v\n", prefix, runes+" ", text, count) // print node
			lines++
			visibleNodes = append(visibleNodes, n)

			if n.open {
				printNodes(n.children, prefix+childBar+"   ") // print children
			}

		}
	}
	printNodes(m.tree, "")
	m.visibleNodes = visibleNodes

	return w.String()
}
func colorWhite(s string) string {
	return fmt.Sprintf("\x1b[97m%s\x1b[0m", s)
}

func colorGreen(s string) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[0m", s)
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
		n := newNode(scanner.Text())

		if n.indent > st.top().indent {
			st.top().Append(n)
			st.push(n)
		} else if n.indent == st.top().indent {
			st.pop()
			st.top().Append(n)
			st.push(n)
		} else {
			for n.indent <= st.top().indent {
				st.pop()
			}
			st.top().Append(n)
			st.push(n)
		}

	}

	return st.root().children
}

func newNode(text string) *node {
	indent := indentDepth(text, '	')
	return &node{
		text:   text[indent:],
		open:   defaultOpen,
		indent: indent,
	}
}

type node struct {
	parent   *node
	text     string
	children []*node
	open     bool
	indent   int
}

func (n *node) Append(child *node) {
	n.children = append(n.children, child)
	child.parent = n
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
