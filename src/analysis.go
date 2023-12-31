package analysis

import (
	"chaitanyabsprip/layout_analysis/util"
	"sort"
	"strings"
)

type Config struct {
	Keymap    [3][10]string
	EffortMap [3][10]float64
	FingerMap [3][10]int
}

type Analyzer struct {
	config    Config
	frequency Frequency
	layout    Layout
}

type Analysis struct {
	Layout        Layout
	EffortRating  float64
	FingerEffort  map[int]float64
	InrollRating  float64
	OutrollRating float64
	RepeatEffort  float64
	SfbRating     float64
	TopSfbs       map[string]float64
}

func NewAnalysis(config Config, frequency Frequency) *Analysis {
	a := NewAnalyzer(config, frequency)
	fe := a.FingerEffort()
	return &(Analysis{
		Layout:        a.layout,
		EffortRating:  a.EffortRating(fe),
		FingerEffort:  fe,
		InrollRating:  a.InrollRating(),
		OutrollRating: a.OutrollRating(),
		RepeatEffort:  a.RepeatEffort(),
		SfbRating:     a.SfbRating(),
		TopSfbs:       a.TopSfbs(),
	})
}

func NewAnalyzer(config Config, frequency Frequency) *Analyzer {
	layout, err := NewLayout(config.Keymap, config.FingerMap)
	if err != nil {
		panic(err)
	}
	return &Analyzer{
		config:    config,
		frequency: frequency,
		layout:    *layout,
	}
}

func (a *Analyzer) EffortRating(fingerEffortMap map[int]float64) float64 {
	rating := 0.0
	for _, v := range fingerEffortMap {
		rating += v
	}
	return rating
}

func (a *Analyzer) FingerEffort() map[int]float64 {
	fingerEffortMap := make(map[int]float64)
	ngramNorm := a.frequency.UnigramNormalised
	for i := 0; i < 3; i++ {
		for j := 0; j < 10; j++ {
			finger := a.config.FingerMap[i][j]
			key := a.layout.keymap[i][j]
			fingerEffortMap[finger] += a.config.EffortMap[i][j] * (ngramNorm[key])
		}
	}
	return fingerEffortMap
}

func (a *Analyzer) InrollRating() float64 {
	ngramNormalised := a.frequency.BigramNormalised
	rating := 0.0
	for _, bigram := range a.layout.Inrolls() {
		rating += ngramNormalised[bigram]
	}
	return rating * 100
}

func (a *Analyzer) OutrollRating() float64 {
	ngramNormalised := a.frequency.BigramNormalised
	rating := 0.0
	for _, bigram := range a.layout.Outrolls() {
		rating += ngramNormalised[bigram]
	}
	return rating * 100
}

func (a *Analyzer) RepeatEffort() float64 {
	ngramNormalised := a.frequency.BigramNormalised
	rating := 0.0
	for _, row := range a.layout.keymap {
		for _, key := range row {
			rating += ngramNormalised[key+key]
		}
	}
	return rating * 100
}

func (a *Analyzer) SfbRating() float64 {
	ngramNormalised := a.frequency.BigramNormalised
	rating := 0.0
	for _, sfb := range a.layout.Sfbs() {
		rating += ngramNormalised[sfb]
	}
	return rating * 100
}

func (a *Analyzer) FingerBigramFrequency(keys [3]string) map[string]float64 {
	freq := make(map[string]float64)
	for _, sfb := range a.layout.SingleSfbs(keys[:]) {
		freq[sfb] = a.frequency.BigramNormalised[sfb] * 100
	}
	return freq
}

func (a *Analyzer) TopSfbs() map[string]float64 {
	ngramNormalised := a.frequency.BigramNormalised
	temp := make(map[string]float64)
	topsfbs := make(map[string]float64)
	for _, sfb := range a.layout.Sfbs() {
		temp[sfb] = ngramNormalised[sfb]
	}
	bigrams := make([]string, len(temp))
	for bigram := range temp {
		bigrams = append(bigrams, bigram)
	}
	sort.SliceStable(bigrams, func(i, j int) bool { return temp[bigrams[i]] > temp[bigrams[j]] })
	for _, bigram := range bigrams[:10] {
		topsfbs[bigram] = temp[bigram] * 100
	}
	return topsfbs
}

type Column struct {
	ID        int
	Frequency float64
}

func sumMapValues(val map[string]float64) float64 {
	sum := 0.0
	for _, v := range val {
		sum += v
	}
	return sum
}

func GetFingerBigramFrequency(analyzer Analyzer, cols map[string]int) map[string]Column {
	freqs := make(map[string]float64)
	colmap := make(map[string]Column)
	for keys := range cols {
		keysArr := [3]string{}
		copy(keysArr[:], strings.Split(keys, ""))
		freqs[keys] = sumMapValues(analyzer.FingerBigramFrequency(keysArr))
	}
	freqs = util.SortMapValues(freqs)
	for k, v := range freqs {
		colmap[k] = Column{ID: cols[k], Frequency: v}
	}
	return colmap
}
