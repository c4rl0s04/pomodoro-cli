package ui

import (
	"strings"
	"testing"
)

func TestCLI_centerBlocks(t *testing.T) {
	c := NewCLI()

	block1 := "Short"
	block2 := "A Very Long Block"

	res1, _ := c.centerBlocks(block1, block2)

	if !strings.HasPrefix(res1, "      ") {
		t.Errorf("expected padding on shorter block, got: %q", res1)
	}
}

func TestCLI_getBlockWidth(t *testing.T) {
	c := NewCLI()

	block := "Line 1\nLonger Line 2\nL3"
	width := c.getBlockWidth(block)
	
	if width != 13 {
		t.Errorf("expected width 13, got %d", width)
	}
}
