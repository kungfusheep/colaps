package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestPrintTree(t *testing.T) {

	for _, test := range []struct {
		text     string
		expected string
	}{
		{`One
	One.One
	One.Two
Two
	Two.One
		Two.One.One
		Two.One.Two
	Two.Two
Three
Four
	Four.One
	Four.Two
	Four.Three
		Four.Three.One
		Four.Three.Two
Five`, `├── One
│   ├── One.One
│   └── One.Two
├── Two
│   ├── Two.One
│   │   ├── Two.One.One
│   │   └── Two.One.Two
│   └── Two.Two
├── Three
├── Four
│   ├── Four.One
│   ├── Four.Two
│   └── Four.Three
│       ├── Four.Three.One
│       └── Four.Three.Two
└── Five
`},
	} {

		w := new(strings.Builder)
		printTree(w, test.text)
		result := w.String()

		fmt.Println(result)

		if result != test.expected {
			t.Errorf("Expected %s but got %s", test.expected, result)
		}
	}
}
