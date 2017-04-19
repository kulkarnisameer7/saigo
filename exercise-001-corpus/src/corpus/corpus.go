package corpus

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func WordCount(file_name string) {

	// Counts: [Word] => Count
	counts := make(map[string]int)

	// Open File
	file, err := os.Open(file_name)
	Check(err)

	// Punctuation-Removal
	var delimiters = strings.NewReplacer("\"", "", ".", "", ",", "", "?", "")

	// Read Lines
	scanner := bufio.NewScanner(file)

	// For Each Line
	for scanner.Scan() {

		// Line (Remove Delimiters)
		line := scanner.Text()
		line = delimiters.Replace(line)

		// Extract Words
		words := strings.Fields(line)

		// Increment Word-Count
		for _, word := range words {

			if word == "" {
				fmt.Println("Got empty word when splitting " + line)
			}
			
			if _, ok := counts[word]; ok {

				// Word Exists
				counts[word]++

			} else {

				// New Word!
				counts[word] = 1
			}
		}
	}

	// Sort Words By Count
	sorted := sortByWordCount(counts)

	// Deisplay
	for _, v := range sorted {
		fmt.Printf("%-10s %5d\n", v.word, v.freq)
	}
}

type Count struct {
	word string
	freq int
}

type Ordered []Count

func (o Ordered) Len() int {
	return len(o)
}

func (o Ordered) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o Ordered) Less(i, j int) bool {
	p := o[i]
	q := o[j]

	if p.freq != q.freq {
		return p.freq > q.freq
	}

	return p.word < q.word
}

func sortByWordCount(wordCounts map[string]int) []Count {

	// Build Counts List
	counts := make([]Count, 0)
	
	for w, f := range wordCounts {

		if f == 0 {
			fmt.Println("Zero: ", w)
		}
		count := Count{word: w, freq: f}
		counts = append(counts, count)
	}

	sort.Sort(Ordered(counts))
	return counts
}
