package main

import (
	"board"
	"fmt"
	"generator"
	"image/png"
	"os"
	"painter"
	"rand"
	"strconv"
	"time"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s width height [output]\n", os.Args[0])
}

func getIntArg(index int, name string) (val int, error os.Error) {
	val, error = strconv.Atoi(os.Args[index])
	if error != nil {
		fmt.Fprintf(os.Stderr, "Invalid %s: %v\n", name, error)
		printUsage()
	}
	return
}

func drawToFile(b board.Board, fileName string) os.Error {
	solution, error := b.Walk(true)
	if error != nil {
		return error
	}
	img := painter.Paint(b, solution, 10, 2)
	file, error := os.Create(fileName)
	defer file.Close()
	if error != nil {
		fmt.Fprintln(os.Stderr, error)
		return error
	}
	error = png.Encode(file, img)
	if error != nil {
		fmt.Fprintln(os.Stderr, error)
		return error
	}
	return nil
}

func main() {
	rand.Seed(time.Nanoseconds())
	if len(os.Args) < 3 || len(os.Args) > 4 {
		printUsage()
		return
	}
	width, error := getIntArg(1, "width")
	if error != nil {
		return
	}
	height, error := getIntArg(2, "height")
	if error != nil {
		return
	}
	b := generator.Generate(width, height)
	if len(os.Args) == 4 {
		error = drawToFile(b, os.Args[3])
		if error != nil {
			fmt.Fprintf(os.Stderr,
				"Error while drawing the maze: %v", error)
		}
	} else {
		fmt.Println(b.PrettyString())
	}
}
