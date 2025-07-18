package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type KeyList struct {
	Keys map[RowCol]Key
	Rows int
	Cols int
}

type Key struct {
	Code   string
	Index  int
	RowCol RowCol
}

// Layout is a representation of a keymap layout. It contains the Name as well as a slice containing all the keycodes.
// It will also contain the exact input that it was parsed from and an Error member that contains any parsing error
// that was found while trying to parse or format the Layout.
type Layout struct {
	Name          string
	KeyCodes      KeyList
	OrigLines     string
	RowColToIndex map[RowCol]int
	LayoutWidth   int
	Error         error
}

func (l *Layout) Format() string {
	if l.Error != nil {
		fmt.Fprintln(os.Stderr, l.Error)
		return l.OrigLines
	}

	formattedString := ""
	// This formats correctly
	// But it requires 2-space tab stop
	colWidths := getColWidths(l.KeyCodes)
	for row := range l.KeyCodes.Rows {
		for col := range l.KeyCodes.Cols {
			keycode := l.KeyCodes.Keys[RowCol{row, col}].Code

			formattedString += fmt.Sprintf(" %-[2]*[1]s", keycode, colWidths[col])
			if len(keycode) == 0 {
				formattedString += " "
			} else {
				formattedString += ","
			}
		}

		formattedString += "\n"
	}

	return formattedString
}

func getColWidths(keys KeyList) []int {
	var widths = make([]int, keys.Cols)
	for _, key := range keys.Keys {
		width := len(key.Code)
		widths[key.RowCol.Col] = max(widths[key.RowCol.Col], width)
	}

	return widths
}

func newKeyCodes(layoutKeys []string, keyIdxToRowCol map[int]RowCol) KeyList {
	var keyList = KeyList{
		Keys: make(map[RowCol]Key),
		Rows: 0,
		Cols: 0,
	}

	for i, key := range layoutKeys {
		rowcol := keyIdxToRowCol[i]
		keyList.Rows = max(keyList.Rows, rowcol.Row+1)
		keyList.Cols = max(keyList.Cols, rowcol.Col+1)
		keyList.Keys[rowcol] = Key{Code: key, Index: i, RowCol: rowcol}
	}

	return keyList
}

// Match the formatted LAYOUT line, e.g [NAME] = LAYOUT (KC_1,KC_2,KC_3....)
var layout_re = regexp.MustCompile(`\[([[:word:]]+)\]\s*=\s*LAYOUT\s*\((.*)\)`)

func removeTabsNSpaces(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(s), "\t", ""), " ", "")
}

// newLayout create a new keymap layout from an input string. It handles windows file endings.
func newLayout(origS string, keyIdxToRowCol map[int]RowCol, rowColToKeyIdx map[RowCol]int, numKeys int) *Layout {
	var current_layer string

	for _, line := range strings.Split(strings.ReplaceAll(origS, "\r\n", "\n"), "\n") {
		// TODO: We loose the line comment here!!!
		lineWithComment := strings.SplitN(line, "//", 2)
		current_layer += removeTabsNSpaces(lineWithComment[0])
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

	return &Layout{Name: name, KeyCodes: keyCodes, OrigLines: origS, RowColToIndex: rowColToKeyIdx}
}

type RowCol struct {
	Row int
	Col int
}

// getMapFromIndexToRowCol returns a map from each key index to which row and column it is in
func getMapFromIndexToRowCol(rowlines [][]int) (map[int]RowCol, map[RowCol]int) {
	var (
		idxToRowCol = make(map[int]RowCol)
		rowColToIdx = make(map[RowCol]int)
	)
	// nrOfRows := len(rowlines)
	// nrOfCols := lines[row]0

	for row := range rowlines {
		for col := 0; col < len(rowlines[row]); col++ {
			// nrOfCols = max(nrOfCol, len())
			keyIdx := rowlines[row][col]
			// key index less than 0 indicates no key present
			if keyIdx < 0 {
				continue
			}

			rowcol := RowCol{Row: row, Col: col}
			idxToRowCol[keyIdx] = rowcol
			rowColToIdx[rowcol] = keyIdx
		}
	}

	return idxToRowCol, rowColToIdx
}

// GetKeymapsFromLines takes all lines from a keymap.c file and returns a []Layout with the keymap layouts that were found
func GetKeymapsFromLines(lines []string, kbConfig *keyboard_t) []*Layout {
	var (
		paren_count  int
		orig_lines   string
		layouts      []*Layout
		keymap_begin = regexp.MustCompile(`.*=\s*LAYOUT.*`)
	)

	keyIdxToRowCol, rowColToKeyIdx := getMapFromIndexToRowCol(kbConfig.Rows)

	for _, line := range lines {
		if keymap_begin.MatchString(line) || len(orig_lines) > 0 {
			paren_count += strings.Count(line, "(") - strings.Count(line, ")")
			orig_lines += line + "\n"
		}

		if paren_count < 1 && len(orig_lines) > 0 {
			layout := newLayout(orig_lines, keyIdxToRowCol, rowColToKeyIdx, kbConfig.Numkeys)
			layouts = append(layouts, layout)
			orig_lines = ""
		}
	}

	return layouts
}
