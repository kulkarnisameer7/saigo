package corpus

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

//Create a struct to protect our lovely counts map

type CountStore struct {
	sync.Mutex
	counts map[string]int
}

//Override New to allow creation of a new count store

func New() *CountStore {
	return &CountStore{
		counts: make(map[string]int),
	}
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func WordCount(file_name string) {

	//Create a work queue
	workQueue := make(chan string)

	//To check if everyone is done
	complete := make(chan bool)

	//define value for concurrent processing
	var concurrency = 5
	
	// Counts: [Word] => Count
	//counts := make(map[string]int)

	//Create a count store here and pass it on to each thread
	countStore := New()
	
	go func() {
		// Open File
		file, err := os.Open(file_name)
		Check(err)

		defer file.Close()
		// Punctuation-Removal
		var delimiters = strings.NewReplacer("\"", "", ".", "", ",", "", "?", "")

		// Read Lines
		scanner := bufio.NewScanner(file)

		// For Each Line
		for scanner.Scan() {

			// Line (Remove Delimiters)
			line := scanner.Text()
			line = delimiters.Replace(line)

			//put the line in the queue
			workQueue <- line

		}

		//close the queue since we have gathered all the lines
		close(workQueue)
	}()
	
	for i:=0; i < concurrency; i++ {

		//Work in parallel on all the lines
		go startLineProcessing(workQueue, complete, countStore)
	}

	for i:=0; i < concurrency; i++ {

		//mark that all work is done, join the threads
		<-complete
	}

	// Sort Words By Count
	sorted := sortByWordCount(countStore.counts)

	// Display
	for _, v := range sorted {
		fmt.Printf("%-10s %5d\n", v.word, v.freq)
	}
}

func startLineProcessing(workQueue chan string,  complete chan bool, cs *CountStore ) {

	for line := range workQueue {
		
		// Extract Words
		words := strings.Fields(line)

		// Increment Word-Count
		for _, word := range words {

			if word == "" {
				fmt.Println("Got empty word when splitting " + line)
			}

			//instead of locking the whole function, just lock the data structure
			cs.Lock()

			if _, ok := cs.counts[word]; ok {

				// Word Exists
				cs.counts[word]++

			} else {

				// New Word!
				cs.counts[word] = 1
			}

			cs.Unlock()
		}
	}

	complete <- true
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
