package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"time"
)

var (
	LENGTH        = flag.Int("l", 50, "length")
	WIDTH         = flag.Int("w", 50, "width")
	TIME_INTERVAL = flag.Int("i", 100, "interval between render")

	White = "\033[37m"
	Reset = "\033[0m"
	Black = "\033[30m"
)

type Field struct {
	Cells  [][]bool
	Length int
	Width  int
}

func NewField(l, w int) *Field {
	cells := make([][]bool, w)
	for i, _ := range cells {
		cells[i] = make([]bool, l)
	}

	return &Field{Length: l, Width: w, Cells: cells}
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
				fmt.Fprintf(&b, "%s██%s", Black, Reset)
			} else {
				fmt.Fprintf(&b, "%s██%s", White, Reset)
			}
		}
		fmt.Fprint(&b, "\n")
	}

	fmt.Print("\033[H")
	fmt.Fprint(os.Stdout, b.String())
	done <- false
}

func (f *Field) RenderField() *Field {
	newField := NewField(f.Length, f.Width)
	directions := [][]int{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
		{-1, 1}, {1, -1}, {1, 1}, {-1, -1},
	}

	for i := 1; i < f.Width-1; i++ {
		for j := 1; j < f.Length-1; j++ {
			ctr := 0
			for _, dir := range directions {
				if f.Cells[i+dir[0]][j+dir[1]] {
					ctr++
				}
			}

			if ctr == 3 {
				newField.Cells[i][j] = true
			} else if ctr == 2 {
				newField.Cells[i][j] = f.Cells[i][j]
			} else {
				newField.Cells[i][j] = false
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
	}

	done := make(chan bool)
	field := NewField(*LENGTH, *WIDTH)
	field.FillRandom()
	ClearScreen()
	go field.PrintField(done)

	for {
		time.Sleep(time.Millisecond * time.Duration(*TIME_INTERVAL))
		newField := field.RenderField()

		<-done

		field.Cells = newField.Cells

		go field.PrintField(done)
	}
}
