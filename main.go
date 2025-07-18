package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/EmilSodergren/go-qmk-keymap/svg"
)

const (
	STATE_HEAD int = iota
	STATE_KEYMAPS
	STATE_KEYMAP
	STATE_TAIL
)

type vizemit_t struct {
	Line  string `json:"line"`
	Layer string `json:"layer"`
}

type keyboard_t struct {
	Name           string            `json:"name"`
	Numkeys        int               `json:"numkeys"`
	Svg            *string           `json:"svg"`
	Rows           [][]int           `json:"rows"`
	Spacing        []int             `json:"spacing"`
	VizWidth       int               `json:"vizcellwidth"`
	VizEmits       []vizemit_t       `json:"vizemits"`
	VizLine        string            `json:"vizline"`
	VizBoard       []string          `json:"vizboard"`
	VizSymbols     map[string]string `json:"vizsymbols"`
	SvgLayers      []string          `json:"svglayers"`
	SvgMapping     [][]int           `json:"svgmapping"`
	SvgSymbolColor map[string]string `json:"svgcolors"`
}

type layer_t struct {
	Name   string
	Keymap []string
	EOLs   []string
}

func is_comment_line(line string) bool {
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, "//") || strings.HasPrefix(line, "*") ||
		strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*/")
}

func sprint_element(element string, separator string, width int) string {
	// left align the element
	// followed by the seperator
	// padded with space until reaching width
	str := make([]rune, 0, width+1)
	cnt := 0

	for _, c := range element {
		str = append(str, c)
		cnt += 1
	}

	for _, c := range separator {
		str = append(str, c)
		cnt += 1
	}

	for cnt < width {
		str = append(str, ' ')
		cnt += 1
	}

	return string(str)
}

func print_formatted(kb *keyboard_t, layer *layer_t) []string {
	output := make([]string, 0, 8)

	width := make([]int, len(kb.Rows[0]))
	for i := range width {
		if len(kb.Spacing) > 0 {
			width[i] = kb.Spacing[0]
		} else {
			width[i] = 8
		}
	}

	for _, row := range kb.Rows {
		for ci, ki := range row {
			if ki >= 0 {
				key := layer.Keymap[ki]
				if len(key) > width[ci] {
					width[ci] = len(key)
				}
			} else {
				ki = -ki
				if ki >= 0 && ki < len(kb.Spacing) {
					if kb.Spacing[ki] > width[ci] {
						width[ci] = kb.Spacing[ki]
					}
				}
			}
		}
	}

	for ri, row := range kb.Rows {
		line := ""

		for i, ki := range row {
			last_column := ki == (kb.Numkeys - 1)
			if ki < 0 {
				line = line + sprint_element("", " ", width[i]+2)
				continue
			}

			if last_column {
				line = line + sprint_element(layer.Keymap[ki], " ", width[i]+2)
			} else {
				line = line + sprint_element(layer.Keymap[ki], ",", width[i]+2)
			}
		}

		// add eol part if present
		if len(layer.EOLs) > ri {
			line = line + layer.EOLs[ri]
		}

		output = append(output, line)
	}

	return output
}

const (
	PARSER_WHITESPACE int = iota
	PARSER_ARRAYITEM
)

// TODO: Create a Regexp for this
func parse_layer_id(line string) string {
	// [ID] = LAYOUT(
	id := make([]rune, 0, 16)

	state := 0
	for _, r := range line {
		if state == 0 {
			if r == '[' {
				state = 1
			}
		} else if state == 1 {
			if r == ']' {
				state = 2
			} else {
				id = append(id, r)
			}
		} else if state == 2 {
			break
		}
	}

	idstr := string(id)
	idstr = strings.TrimSpace(idstr)

	return idstr
}

func parse_elements(line string) ([]string, string) {
	keymap := make([]string, 0, 80)

	eol_pos := len(line)
	for cursor, r := range line {
		if r == '/' || r == '\\' {
			eol_pos = cursor
			break
		}
	}

	// snip the end of the line part and remember it
	end_of_line_part := line[eol_pos:]

	// continue with the trimmed line
	line = line[0:eol_pos]

	state := PARSER_WHITESPACE
	open := 0
	elem := make([]rune, 0, 60)

	for _, r := range line {
		if state == PARSER_WHITESPACE {
			if !unicode.IsSpace(r) {
				state = PARSER_ARRAYITEM
			}
		}

		if state == PARSER_ARRAYITEM {
			if open > 0 {
				if r == '[' || r == '(' {
					open += 1
				} else if r == ']' || r == ')' {
					open -= 1
				}

				elem = append(elem, r)
			} else {
				if r == ',' {
					state = PARSER_WHITESPACE
					elemstr := string(elem)
					elemstr = strings.TrimSpace(elemstr)
					keymap = append(keymap, elemstr)
					elem = elem[:0]
				} else {
					if r == '[' || r == '(' {
						open += 1
					}

					elem = append(elem, r)
				}
			}
		}
	}

	if len(elem) > 0 {
		elemstr := string(elem)

		elemstr = strings.TrimSpace(elemstr)
		if len(elemstr) > 0 {
			keymap = append(keymap, elemstr)
		}
	}

	return keymap, end_of_line_part
}

func parse_viz_layer_names(line string) []string {
	layer_names := make([]string, 0, 4)
	layers_str := make([]rune, 0, 32)

	state := 0
	for _, c := range line {
		if state == 0 {
			if c == '[' {
				state = 1
			}
		} else if state == 1 {
			if c == ']' {
				state = 2
			} else {
				layers_str = append(layers_str, c)
			}
		} else if state == 2 {
			break
		}
	}

	layers := strings.Split(string(layers_str), ",")
	for _, l := range layers {
		layer_name := strings.TrimSpace(l)
		layer_names = append(layer_names, layer_name)
	}

	return layer_names
}

func index_to_viz(index int) string {
	return fmt.Sprintf("_%03d_", index)
}

func print_viz(k *keyboard_t, layer *layer_t) []string {
	keyboardviz := make([]string, 0, 32)
	keyboardviz = append(keyboardviz, k.VizBoard...)

	spacing := "                          "

	for ki, key := range layer.Keymap {
		keysymbol := ""
		if symbol, ok := k.VizSymbols[key]; ok {
			keysymbol = symbol
		} else if len(key) > k.VizWidth {
			keysymbol = key[(len(key) - k.VizWidth):]
		} else if len(key) < k.VizWidth {
			lp := (k.VizWidth - len(key)) / 2
			rp := k.VizWidth - len(key) - lp
			keysymbol = fmt.Sprintf("%s%s%s", spacing[0:lp], key, spacing[0:rp])
		} else {
			keysymbol = key
		}

		marker := index_to_viz(ki)
		for i := 0; i < len(keyboardviz); i++ {
			keyboardviz[i] = strings.Replace(keyboardviz[i], marker, keysymbol, 1)
		}
	}

	return keyboardviz
}

// convert all know string characters to svg escape code
func escape_svg(s string) string {
	output := make([]rune, 0, len(s))
	for _, r := range s {
		if r == '<' {
			output = append(output, '&', 'l', 't', ';')
		} else if r == '>' {
			output = append(output, '&', 'g', 't', ';')
		} else if r == '&' {
			output = append(output, '&', 'a', 'm', 'p', ';')
		} else if r == '"' {
			output = append(output, '&', 'q', 'u', 'o', 't', ';')
		} else if r == '\'' {
			output = append(output, '&', 'a', 'p', 'o', 's', ';')
		} else {
			output = append(output, r)
		}
	}

	return string(output)
}

func print_svg(keyboard *keyboard_t, layers map[string]*layer_t, workdir string) {
	svgLayers := make([]svg.Layer_t, 0, 32)

	for _, layername := range keyboard.SvgLayers {
		layer, has_layer := layers[layername]
		if has_layer {
			svgLayer := svg.Layer_t{}
			svgLayer.Matrix = [][]svg.Key_t{}
			svgLayer.Name = keyboard.VizSymbols[layer.Name]

			for _, row := range keyboard.SvgMapping {
				svgRow := make([]svg.Key_t, 0, 20)
				for _, key := range row {
					svgKey := svg.Key_t{}

					svgKey.Exists = key >= 0
					if key >= 0 {
						svgKey.Key = layer.Keymap[key]
						svgKey.Key = keyboard.VizSymbols[svgKey.Key]
						svgKey.Key = strings.TrimSpace(svgKey.Key)
						svgKey.Key = escape_svg(svgKey.Key)
					}

					svgKeyClass, hasClass := keyboard.SvgSymbolColor[svgKey.Key]
					if hasClass {
						svgKey.Class = svgKeyClass
					} else {
						svgKey.Class = ""
					}

					svgRow = append(svgRow, svgKey)
				}

				svgLayer.Matrix = append(svgLayer.Matrix, svgRow)
			}

			svgLayers = append(svgLayers, svgLayer)
		}
	}

	svgLines := svg.Print(svgLayers)

	// write all lines to a text file
	if len(strings.TrimSpace(*keyboard.Svg)) > 0 {
		f, err := os.Create(filepath.Join(workdir, *keyboard.Svg))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer f.Close()

		for _, line := range svgLines {
			f.WriteString(line)
			f.WriteString("\n")
		}
	}
}

type Args struct {
	WorkingDir string
}

var configPath = regexp.MustCompile(`^/[/\*]\s*(\S+.json)`)

func parseArgs() Args {
	var args Args

	fs := flag.NewFlagSet("go-qmk-keymap", flag.ExitOnError)
	fs.StringVar(&args.WorkingDir, "workdir", "", "If provided, paths to the qmk-keyboard-format.json can be relative this path")
	fs.Parse(os.Args[1:])

	return args
}

func main() {
	args := parseArgs()

	err := run(args)
	if err != nil {
		log.Fatal(err)
	}
}

func findEmit(line string, kb *keyboard_t) int {
	line = strings.TrimSpace(line)
	for emitidx, emitat := range kb.VizEmits {
		if strings.Compare(line, emitat.Line) == 0 {
			return emitidx
		}
	}

	return -1
}

func getConfigData(firstLine string, args Args) (*keyboard_t, error) {
	keymapPath := configPath.FindStringSubmatch(firstLine)
	if len(keymapPath) == 0 {
		return nil, fmt.Errorf("no configuration found on first line of keymap.c")
	}

	// Checks if config path is absolute
	kbdata, err := os.ReadFile(keymapPath[1])
	if err != nil {
		// Checks if config path relative workdir
		kbdata, err = os.ReadFile(filepath.Join(args.WorkingDir, keymapPath[1]))
		if err != nil {
			return nil, fmt.Errorf("no configuration available: %v", err)
		}
	}

	kb, err := UnmarshalKeyboard(kbdata)
	if err != nil {
		return nil, fmt.Errorf("configuration file not parseable: %v", err)
	}

	return kb, nil
}

var test_re = regexp.MustCompile(`\s*\[(?<layer_id>[[:word:]]+)\] = LAYOUT\s*\(`)

func run(args Args) error {
	scanner := bufio.NewScanner(os.Stdin)
	lines := make([]string, 0, 4096)

	scanner.Scan()
	firstLine := scanner.Text()
	lines = append(lines, firstLine)

	var kb, err = getConfigData(firstLine, args)
	if err != nil {
		return err
	}

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	var (
		keymaps_begin = "const uint16_t PROGMEM "
		keymaps_end   = "};"
		keymap_end1   = "),"
		keymap_end2   = ")"
		keymap_begin  = "= LAYOUT"
	)

	output := make([]string, 0, 1024)
	layers := make(map[string]*layer_t)

	state := STATE_HEAD
	var currentLayer *layer_t

	var keymapLines []string
	var keymapLayers []*Layout
	for _, line := range lines {
		// here we check if the line is a '//' comment
		if is_comment_line(line) {
			// If it is NOT a "special" VizLine comment, we keep it.
			// The VizLines will be added later, so skip them to avoid duplicates
			if !strings.HasPrefix(strings.TrimSpace(line), kb.VizLine) {
				output = append(output, line)
			}
			continue
		}

		switch state {
		// Before the KeyMaps, append the lines
		case STATE_HEAD:
			if strings.Contains(line, keymaps_begin) {
				state = STATE_KEYMAPS
			}
			output = append(output, line)
		// Found where the KeyMaps definitions starts
		case STATE_KEYMAPS:
			if strings.Contains(line, keymap_begin) {
				layer_name := parse_layer_id(line)
				currentLayer = &layer_t{Name: layer_name, Keymap: make([]string, 0, kb.Numkeys), EOLs: make([]string, len(kb.Rows))}
				layers[layer_name] = currentLayer
				state = STATE_KEYMAP
				keymapLines = append(keymapLines, line)
			} else if strings.HasSuffix(strings.TrimSpace(line), keymaps_end) {
				keymapLayers = GetKeymapsFromLines(keymapLines, kb)
				state = STATE_TAIL
			}
			output = append(output, line)
		// Parsing a KeyMap definiton
		case STATE_KEYMAP:
			if strings.TrimSpace(line) == keymap_end1 || strings.TrimSpace(line) == keymap_end2 {
				// do we have a parsed keymap, if so append the formatted lines to output
				formatted := print_formatted(kb, currentLayer)
				output = append(output, formatted...)
				currentLayer = nil
				state = STATE_KEYMAPS

				output = append(output, line)
			} else {
				// collect the elements from these lines
				elems, eol_part := parse_elements(line)
				currentLayer.Keymap = append(currentLayer.Keymap, elems...)
				currentLayer.EOLs = append(currentLayer.EOLs, eol_part)
			}
			keymapLines = append(keymapLines, line)
		// Found where the KeyMaps definitins ends
		case STATE_TAIL:
			output = append(output, line)
		}
	}

	if len(kb.VizEmits) > 0 {
		// iterate over all output lines and identify locations where we need to output ascii-art for a layer
		doviz := 0
		output_temp := make([]string, 0, 1024)

		for _, line := range output {
			emitindex := findEmit(line, kb)
			if emitindex >= 0 {
				layer_name := kb.VizEmits[emitindex].Layer
				if layer, ok := layers[layer_name]; ok {
					layer_viz := print_viz(kb, layer)
					output_temp = append(output_temp, layer_viz...)
				}

				output_temp = append(output_temp, line)
				doviz = 0
			} else if doviz == 0 {
				output_temp = append(output_temp, line)
			}
		}

		output = output_temp
	}

	for _, l := range keymapLayers {
		l.Format() // <-- print this to use the new formatter
	}

	for _, l := range output {
		fmt.Println(strings.TrimRight(l, " "))
	}

	if kb.Svg != nil {
		print_svg(kb, layers, args.WorkingDir)
	}

	return nil
}

func UnmarshalKeyboard(data []byte) (*keyboard_t, error) {
	r := &keyboard_t{}
	err := json.Unmarshal(data, &r)

	return r, err
}

func (r *keyboard_t) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
