package main

import "testing"

func TestCreateNewLayout(t *testing.T) {
	origLayout := "[TEST_NAME] = LAYOUT(A1,A2,A3)"
	m := map[int]RowCol{0: {0, 0}, 1: {0, 1}, 2: {0, 2}}

	l := newLayout(origLayout, 3, m)
	if l.Error != nil {
		t.Fatal(l.Error)
	}
}

func TestCreateNewLayoutWithLinebreaks(t *testing.T) {
	origLayout := `[TEST_NAME]
	= LAYOUT(
		A1,A2
		,A3

	)`
	m := map[int]RowCol{0: {0, 0}, 1: {0, 1}, 2: {0, 2}}

	l := newLayout(origLayout, 3, m)
	if l.Error != nil {
		t.Fatal(l.Error)
	}
}
