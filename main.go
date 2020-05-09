package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"
)

const (
	wordsPerBook = 100000 // approximate number of words per book; http://wordcounters.com/word-count/literary-books/
	threshold    = 12     // number of times you need to see a word passively to learn it; https://www.readingrockets.org/blogs/shanahan-literacy/how-many-times-should-students-copy-spelling-words
	wordsToKnow  = 35000  // words an educated English speaker knows; http://testyourvocab.com/blog/2013-05-08-Native-speakers-in-greater-detail
	maxVocabSize = 100000 // reasonable total English vocabulary size in books to not generate Zipf outliers

	trials = 1000 // run the Monte Carlo simulation n times; note on my 8-core laptop, this would take ~3 minutes
)

// simpleTrial makes a lot of assumptions, and doesn't take into account
// time between seeing words, as spaced repetition would
func simpleTrial() uint64 {
	seen := make([]uint64, wordsToKnow)
	var (
		knownWords uint64
		wordsSeen  uint64
	)
	// Assume a constant stream of words
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	zipf := rand.NewZipf(r, 1.01, 1.0, wordsToKnow-1)
	if zipf == nil {
		log.Fatal("Can't generate Zipf distribution")
	}
	for ; knownWords < wordsToKnow; wordsSeen++ {
		// Ord is the ordinal value in a frequency list (i.e. the 30th most common word)
		// assuming a vocabulary distributed per Zipf's law, capped at maxVocabSize
		ord := zipf.Uint64()
		seen[ord]++
		if seen[ord] == threshold {
			knownWords++
		}
	}
	return wordsSeen
}

func work(jobs <-chan int, results chan<- uint64) {
	for _ = range jobs {
		results <- simpleTrial()
	}
}

func main() {
	var total uint64
	// Run all simulations over max processors on the machine in parallel,
	// to speed up simulation time
	jobs := make(chan int, trials)
	results := make(chan uint64, trials)
	for w := 0; w < runtime.NumCPU(); w++ {
		go work(jobs, results)
	}
	for i := 0; i < trials; i++ {
		jobs <- i
	}
	close(jobs)
	for i := 0; i < trials; i++ {
		total += <-results
	}
	// Calculate and print the results
	avg := float64(total) / float64(trials)
	log.Printf("Avg %.0f words, or approx %.1f books, over %d trials", avg, float64(avg)/float64(wordsPerBook), trials)
}
