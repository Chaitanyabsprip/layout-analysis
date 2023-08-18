package analysis

import (
	"encoding/json"
	"os"
	"strings"
)

func GetCorpus() string {
	content, _ := os.ReadFile("/Users/chaitanyasharma/projects/languages/go/layout_analysis/data/corpus.txt")
	return strings.ToLower(string(content))
}

type Frequency struct {
	UnigramCount      map[string]int
	BigramCount       map[string]int
	UnigramNormalised map[string]float64
	BigramNormalised  map[string]float64
}

func NewFrequency() *Frequency {
	c := corpus{corpus: GetCorpus()}
	uniC, uniN := c.ngramFrequency(1)
	biC, biN := c.ngramFrequency(2)
	return &Frequency{
		UnigramCount:      uniC,
		BigramCount:       biC,
		UnigramNormalised: uniN,
		BigramNormalised:  biN,
	}
}

type corpus struct {
	corpus string
}

type nGramType int

const (
	invalid nGramType = iota - 1
	uniG
	biG
	uniGN
	biGN
)

func (n nGramType) Filename() string {
	prefix := "/Users/chaitanyasharma/projects/languages/go/layout_analysis/generated/"
	var name string
	switch n {
	case uniG:
		name = "unigram.json"
	case biG:
		name = "bigram.json"
	case uniGN:
		name = "unigram_norm.json"
	case biGN:
		name = "bigram_norm.json"
	}
	if name == "" {
		return "Unknown data type"
	}
	return prefix + name
}

func copyMap[VAL int | float64](target map[string]VAL, source map[string]VAL) {
	for k, v := range source {
		(target)[k] = v
	}
}

func getNgram(n int, isCount bool) nGramType {
	if isCount {
		switch n {
		case 1:
			return uniG
		case 2:
			return biG
		}
	}
	switch n {
	case 1:
		return uniGN
	case 2:
		return biGN
	}
	return invalid
}

func saveToFile[T int | float64](val map[string]T, ngram nGramType) {
	filename := ngram.Filename()
	json, err := json.Marshal(val)
	if err == nil {
		_ = os.WriteFile(filename, json, 0644)
	}
}

func readFromFile(ngramtype nGramType) []byte {
	content, err := os.ReadFile(ngramtype.Filename())
	if err != nil {
		return nil
	}
	return content
}

func getCache[T int | float64](ngramtype nGramType) map[string]T {
	jsonStr := readFromFile(ngramtype)
	val := make(map[string]T)
	err := json.Unmarshal(jsonStr, &val)
	if err == nil {
		return val
	}
	return nil
}

func (f *corpus) ngramCount(n int) map[string]int {
	ngramtype := getNgram(n, true)
	cached := getCache[int](ngramtype)
	if cached != nil {
		return cached
	}
	ngrams := make(map[string]int)
	chars := strings.Split(strings.TrimSpace(f.corpus), "")
	for i := 0; i < len(chars)-n+1; i++ {
		ngram := strings.Join(chars[i:i+n], "")
		if strings.ContainsAny(ngram, " \n\t") {
			continue
		}
		ngrams[ngram]++
	}
	saveToFile(ngrams, ngramtype)
	return ngrams
}

func sum(slice map[string]int) int {
	var sum int
	for _, v := range slice {
		sum += v
	}
	return sum
}

func normalise(slice map[string]int, total int) map[string]float64 {
	normalised := make(map[string]float64)
	for k, v := range slice {
		normalised[k] = float64(v) / float64(total)
	}
	return normalised
}

func (f *corpus) ngramFrequency(n int) (map[string]int, map[string]float64) {
	ngramtype := getNgram(n, false)
	count := f.ngramCount(n)
	cached := getCache[float64](ngramtype)
	if cached != nil {
		return count, cached
	}
	total := sum(count)
	normalised := normalise(count, total)
	saveToFile(normalised, ngramtype)
	return count, normalised
}
