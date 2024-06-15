package main

import (
	"fmt"
	"regexp"
	"strings"
)

type KeyList []Key

type Key struct {
	Code   string
	Index  int
	Column int
}

// Layout is a representation of a keymap layout. It contains the Name as well as a slice containing all the keycodes.
// It will also contain the exact input that it was parsed from and an Error member that contains any parsing error
// that was found while trying to parse or format the Layout.
type Layout struct {
	Name      string
	KeyCodes  KeyList
	OrigLines string
	Columns   int
	Error     error
}

func (l *Layout) Format() string {
	if l.Error != nil {
		return l.OrigLines
	}
	colWidths := getColWidths(l.KeyCodes, l.Columns)
	for _, key := range l.KeyCodes {
		fmtS := fmt.Sprintf("%d")
		fmt.Printf()
	}
	return ""
}

func getColWidths(keys KeyList, nrOfColumns int) []int {
	var widths = make([]int, nrOfColumns)
	for _, key := range keys {
		width := len(key.Code)
		if width > widths[key.Column] {
			widths[key.Column] = width
		}
	}
	return widths
}

func newKeyCodes(layoutKeys []string, keyIdxToRowCol map[int]RowCol) KeyList {
	var keyList KeyList
	for i, key := range layoutKeys {
		keyList = append(keyList, Key{Code: key, Index: i, Column: keyIdxToRowCol[i].Col})
	}
	return keyList
}

// Match the formatted LAYOUT line, e.g [NAME] = LAYOUT (KC_1,KC_2,KC_3....)
var layout_re = regexp.MustCompile(`\[(.*)\]\s*=\s*LAYOUT\s*\((.*)\)`)

func removeTabsNSpaces(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(s), "\t", ""), " ", "")
}

// newLayout create a new keymap layout from an input string. It handles windows file endings.
func newLayout(origS string, keyIdxToRowCol map[int]RowCol, numKeys, nrOfColumns int) *Layout {

	var current_layer string
	for _, line := range strings.Split(strings.ReplaceAll(origS, "\r\n", "\n"), "\n") {
		current_layer += removeTabsNSpaces(line)
	}
	matches := layout_re.FindStringSubmatch(current_layer)
	if len(matches) != 3 {
		return &Layout{Error: fmt.Errorf("layout is badly formatted"), OrigLines: origS}
	}
	name := matches[1]
	layoutKeys := strings.Split(matches[2], ",")
	if len(layoutKeys) != numKeys {
		return &Layout{Error: fmt.Errorf("Layer %s contains %d keys but %d was expected", name, len(layoutKeys), numKeys), OrigLines: origS}
	}
	keyCodes := newKeyCodes(layoutKeys, keyIdxToRowCol)

	return &Layout{Name: name, KeyCodes: keyCodes, OrigLines: origS, Columns: nrOfColumns}
}

type RowCol struct {
	Row int
	Col int
}

// getMapFromIndexToRowCol returns a map from each key index to which row and column it is in
func getMapFromIndexToRowCol(rowlines [][]int) map[int]RowCol {
	var idxToRowCol = make(map[int]RowCol)
	for col := range rowlines[0] {
		for row := range rowlines {
			keyIdx := rowlines[row][col]
			// key index less than 0 indicates no key present
			if keyIdx < 0 {
				continue
			}
			idxToRowCol[keyIdx] = RowCol{Row: row, Col: col}
		}
	}
	return idxToRowCol
}

// GetKeymapsFromLines takes all lines from a keymap.c file and returns a []Layout with the keymap layouts that were found
func GetKeymapsFromLines(lines []string, kbConfig *keyboard_t) []*Layout {
	var paren_count int
	var orig_lines string
	var layouts []*Layout
	var keymap_begin = regexp.MustCompile(`.*=\s*LAYOUT.*`)
	keyIdxToRowCol := getMapFromIndexToRowCol(kbConfig.Rows)
	for _, line := range lines {
		if keymap_begin.MatchString(line) || len(orig_lines) > 0 {
			paren_count += strings.Count(line, "(") - strings.Count(line, ")")
			orig_lines += line + "\n"
		}
		if paren_count < 1 && len(orig_lines) > 0 {
			layouts = append(layouts, newLayout(orig_lines, keyIdxToRowCol, kbConfig.Numkeys, len(kbConfig.Rows[0])))
			orig_lines = ""
		}
	}
	return layouts
}
