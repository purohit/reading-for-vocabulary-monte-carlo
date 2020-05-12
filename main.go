package main

import (
	"flag"
	"log"
	"math/rand"
	"runtime"
	"time"
)

const (
	// Zipf parameters
	zipf_s = 1.00001 // s ~= 1 is for word distributions
	zipf_v = 1.0     // v = 1

	wordsPerBook = 85000  // approximate number of words per book
	wordsToLearn = 35000  // words an educated English speaker knows; http://testyourvocab.com/blog/2013-05-08-Native-speakers-in-greater-detail
	vocabSize    = 171146 // reasonable total English vocabulary size to cap Zipf outliers; https://wordcounter.io/blog/how-many-words-are-in-the-english-language/
)

var (
	// see flags for documentation
	trials     int // on my 8-core laptop, 1000 trials take ~3 minutes
	experiment string
	threshold  uint64
)

// generates a Zipf distribution for a constant stream of words
func getZipf() *rand.Zipf {
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	return rand.NewZipf(r, zipf_s, zipf_v, vocabSize-1)
}

// srsTrial is the same as simpleTrial, but takes into
// account learning at exponentially increasing gaps (based on the theories of SRS, or spaced-repetition systems).
// here, a word only counts as "seen another time" when the last time we saw it was more than twice as long ago as the previous time
func srsTrial() uint64 {
	type lti struct {
		n        uint64 // seen this word n total times
		k        uint64 // seen this word k "good" (i.e. with proper spacing) times
		lastSeen uint64
	}
	var (
		learnedWords uint64
		counter      uint64 // monotonically increasing clock
	)
	seen := make([]*lti, wordsToLearn)
	for i := range seen {
		seen[i] = new(lti)
	}
	zipf := getZipf()
	for ; learnedWords < wordsToLearn; counter++ {
		word := zipf.Uint64()
		if word >= wordsToLearn { // beyond the vocabulary we're interested in
			continue
		}
		seen[word].n++
		interval := counter - seen[word].lastSeen
		if interval > 2<<seen[word].k {
			seen[word].k++
			seen[word].lastSeen = counter
			if seen[word].k == threshold {
				learnedWords++
			}
		}
	}
	return counter
}

// simpleTrial makes a lot of assumptions, and doesn't take into account
// time between seeing words, as spaced repetition would
func simpleTrial() uint64 {
	seen := make([]uint64, wordsToLearn)
	var (
		learnedWords uint64
		wordsSeen    uint64
	)
	// Assume a constant stream of words
	zipf := getZipf()
	for ; learnedWords < wordsToLearn; wordsSeen++ {
		// word is the ordinal value in a frequency list (i.e. the 30th most common word)
		// assuming a vocabulary distributed per Zipf's law, capped at vocabSize
		word := zipf.Uint64()
		if word >= wordsToLearn { // beyond the vocabulary we're interested in
			continue
		}
		seen[word]++
		if seen[word] == threshold {
			learnedWords++
		}
	}
	return wordsSeen
}

func work(jobs <-chan int, results chan<- uint64) {
	for _ = range jobs {
		if experiment == "simple" {
			results <- simpleTrial()
		} else if experiment == "srs" {
			results <- srsTrial()
		}
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
	log.Printf("%.0f words read on average, or approx %.1f books, over %d %s trials", avg, float64(avg)/float64(wordsPerBook), trials, experiment)
}

func init() {
	flag.StringVar(&experiment, "experiment", "simple", "experiment to run (simple or srs)")
	flag.IntVar(&trials, "trials", 10, "number of trials to run")
	flag.Uint64Var(&threshold, "threshold", 12, "number of times to see a word before it's considered learned")
	flag.Parse()
}
