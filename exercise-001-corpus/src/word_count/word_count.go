package main

import (
	"fmt"
	"os"
	"github.com/kulkarnisameer7/saigo/exercise-001-corpus/src/corpus"
)

func main() {
	files := os.Args[1:]
	fmt.Println("File path: ", files)
	for i := 0; i < len(files); i++ {
		filename := string(files[i])
		corpus.WordCount(filename)
	}
}
