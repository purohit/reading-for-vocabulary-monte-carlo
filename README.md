# reading-for-vocabulary-monte-carlo

Code to run Monte Carlo simulations for my blog post: https://puroh.it/reading-for-a-fine-vocabulary/.

Usage:
```go
go build
./reading-for-vocabulary-monte-carlo
# Use --help for options
```

Running parallel jobs (to get statistics for 20000 to 40000 vocabulary sizes,
for example):

```
parallel -j1 --tag --eta './reading-for-vocabulary-monte-carlo -words-to-learn {} -trials 100' ::: (seq 20000 1000 40000) > out.txt
```
