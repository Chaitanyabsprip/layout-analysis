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
	Corpus string
}

type nGramType int

const (
	invalid nGramType = iota - 1
	uniG
	biG
	uniGN
	biGN
)

var (
	unigramCount = make(map[string]int)
	bigramCount  = make(map[string]int)
	unigramNorm  = make(map[string]float64)
	bigramNorm   = make(map[string]float64)
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

func saveToVar[T int | float64](val map[string]T, ngram nGramType) {
	switch ngram {
	case uniG:
		if source, ok := any(val).(map[string]int); ok {
			copyMap(unigramCount, source)
		}
	case biG:
		if source, ok := any(val).(map[string]int); ok {
			copyMap(bigramCount, source)
		}
	case uniGN:
		if source, ok := any(val).(map[string]float64); ok {
			copyMap(unigramNorm, source)
		}
	case biGN:
		if source, ok := any(val).(map[string]float64); ok {
			copyMap(bigramNorm, source)
		}
	}
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

func getCacheFromFile(ngramtype nGramType) any {
	jsonStr := readFromFile(ngramtype)
	val := make(map[string]int)
	err := json.Unmarshal(jsonStr, &val)
	if err == nil {
		saveToVar(val, ngramtype)
		return val
	}
	return nil
}

func getCache(ngramtype nGramType) any {
	switch ngramtype {
	case uniG:
		if len(unigramCount) != 0 {
			return unigramCount
		}
	case biG:
		if len(bigramCount) != 0 {
			return bigramCount
		}
	case uniGN:
		if len(unigramNorm) != 0 {
			return unigramNorm
		}
	case biGN:
		if len(bigramNorm) != 0 {
			return bigramNorm
		}
	}
	return getCacheFromFile(ngramtype)
}

func (f *Frequency) NgramCount(n int) map[string]int {
	ngramtype := getNgram(n, true)
	cached := getCache(ngramtype)
	if ng, ok := cached.(map[string]int); ok && ng != nil {
		return ng
	}
	ngrams := make(map[string]int)
	chars := strings.Split(strings.TrimSpace(f.Corpus), "")
	for i := 0; i < len(chars)-n+1; i++ {
		ngram := strings.Join(chars[i:i+n], "")
		if strings.ContainsAny(ngram, " \n\t") {
			continue
		}
		ngrams[ngram]++
	}
	saveToVar(ngrams, ngramtype)
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

func (f *Frequency) NgramNormalised(n int) map[string]float64 {
	ngramtype := getNgram(n, false)
	cached := getCache(ngramtype)
	if ng, ok := cached.(map[string]float64); ok && ng != nil {
		return ng
	}
	ngrams := f.NgramCount(n)
	total := sum(ngrams)
	normalised := normalise(ngrams, total)
	saveToVar(normalised, ngramtype)
	saveToFile(normalised, ngramtype)
	return normalised
}
