package main

import (
	"fmt"
	"generator"
	"os"
	"strconv"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s width height\n", os.Args[0])
}

func getIntArg(index int, name string) (val int, error os.Error) {
	val, error = strconv.Atoi(os.Args[index])
	if error != nil {
		fmt.Fprintf(os.Stderr,"Invalid %s: %v\n", name, error)
		printUsage()
	}
	return
}

func main() {
	if len(os.Args) < 3 {
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
	fmt.Println(generator.Generate(width, height).PrettyString())
}
