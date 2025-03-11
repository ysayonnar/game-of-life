package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

//INTERESTING RULES
//b3678/s34678 - day and night
//b3/s - seeds

var (
	LENGTH        = flag.Int("l", 50, "length")
	WIDTH         = flag.Int("w", 50, "width")
	TIME_INTERVAL = flag.Int("i", 100, "interval between render")
	RULE_FORMAT   = flag.String("rule", "b3/s23", "rule of game in format b3/s23")

	White = "\033[37m"
	Reset = "\033[0m"
	Black = "\033[30m"
)

type Rule struct {
	Born    []int
	Survive []int
}

type Field struct {
	Cells  [][]bool
	Length int
	Width  int
	Rule   Rule
}

func ParseRule(codedRule string) (*Rule, error) {
	var InvalidFormatError = errors.New("incorrect rule, must be b<nums>/s<nums>")
	rule := &Rule{Born: []int{}, Survive: []int{}}

	s := strings.Split(codedRule, "/")
	if len(s) != 2 || s[0][0] != 'b' || s[1][0] != 's' {
		return nil, InvalidFormatError
	}

	bornStr, survivestr := s[0][1:], s[1][1:]

	for _, s := range bornStr {
		if !unicode.IsDigit(s) {
			return nil, InvalidFormatError
		} else {
			n, _ := strconv.Atoi(string(s))
			rule.Born = append(rule.Born, n)
		}
	}

	for _, s := range survivestr {
		if !unicode.IsDigit(s) {
			return nil, InvalidFormatError
		} else {
			n, _ := strconv.Atoi(string(s))
			rule.Survive = append(rule.Survive, n)
		}
	}

	return rule, nil
}

func NewField(l, w int, rule Rule) *Field {
	cells := make([][]bool, w)
	for i, _ := range cells {
		cells[i] = make([]bool, l)
	}

	return &Field{Length: l, Width: w, Cells: cells, Rule: rule}
}

func (f *Field) FillRandom() {
	for i := 0; i < f.Width; i++ {
		for j := 0; j < f.Length; j++ {
			cell := rand.IntN(2)
			if cell == 0 {
				f.Cells[i][j] = false
			} else {
				f.Cells[i][j] = true
			}
		}
	}
}

func (f *Field) PrintField(done chan bool) {
	var b bytes.Buffer

	for i := 0; i < f.Width; i++ {
		for j := 0; j < f.Length; j++ {
			cell := f.Cells[i][j]
			if cell {
				fmt.Fprintf(&b, "%s██%s", White, Reset)
			} else {
				fmt.Fprintf(&b, "  ")
			}
		}
		fmt.Fprint(&b, "\n")
	}

	fmt.Print("\033[H")
	fmt.Fprint(os.Stdout, b.String())
	done <- false
}

func (f *Field) RenderField() *Field {
	newField := NewField(f.Length, f.Width, f.Rule)
	directions := [][]int{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
		{-1, 1}, {1, -1}, {1, 1}, {-1, -1},
	}

	for i := 0; i < f.Width; i++ {
		for j := 0; j < f.Length; j++ {
			ctr := 0
			for _, dir := range directions {
				ni := (i + dir[0] + f.Width) % f.Width
				nj := (j + dir[1] + f.Length) % f.Length

				if f.Cells[ni][nj] {
					ctr++
				}
			}

			newField.Cells[i][j] = false
			for _, surviveNum := range f.Rule.Survive {
				if ctr == surviveNum {
					newField.Cells[i][j] = f.Cells[i][j]
					break
				}
			}

			for _, bornNum := range f.Rule.Born {
				if ctr == bornNum {
					newField.Cells[i][j] = true
				}
			}
		}
	}
	return newField
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func main() {
	flag.Parse()

	if *TIME_INTERVAL < 0 {
		fmt.Println("Incorrect interval")
		return
	} else if *LENGTH < 0 || *LENGTH > 1000 {
		fmt.Println("length must be from 0 to 1000")
		return
	} else if *WIDTH < 0 || *WIDTH > 1000 {
		fmt.Println("width must be from 0 to 1000")
		return
	}

	rule, err := ParseRule(*RULE_FORMAT)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	field := NewField(*LENGTH, *WIDTH, *rule)
	field.FillRandom()
	ClearScreen()

	done := make(chan bool)
	go field.PrintField(done)

	for {
		time.Sleep(time.Millisecond * time.Duration(*TIME_INTERVAL))
		newField := field.RenderField()

		<-done

		field.Cells = newField.Cells

		go field.PrintField(done)
	}
}
