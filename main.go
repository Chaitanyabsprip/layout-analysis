package main

import analysis "chaitanyabsprip/layout_analysis/src"

var (
	effortMap = [3][10]float64{
		{3.0, 2.0, 1.5, 1.7, 3.2, 3.2, 1.7, 1.5, 2.0, 3.0},
		{1.6, 1.3, 1.0, 1.0, 2.0, 2.0, 1.0, 1.0, 1.3, 1.6},
		{2.5, 1.8, 1.7, 1.3, 3.0, 3.0, 1.3, 1.7, 1.8, 2.5},
	}
	fingerMap = [3][10]int{
		{0, 1, 2, 3, 3, 6, 6, 7, 8, 9},
		{0, 1, 2, 3, 3, 6, 6, 7, 8, 9},
		{0, 1, 2, 3, 3, 6, 6, 7, 8, 9},
	}
)

func main() {
	corpus := analysis.GetCorpus()
	freq := analysis.Frequency{Corpus: corpus}
	keymap := [3][10]string{
		{"/", "v", "w", "m", "j", "z", "f", "u", "x", "q"},
		{"a", "s", "r", "t", "g", "p", "n", "e", "o", "i"},
		{".", "c", "l", "d", "k", "b", "h", "'", ",", "y"},
	}

	config := analysis.Config{Keymap: keymap, EffortMap: effortMap, FingerMap: fingerMap}
	anal := analysis.NewAnalysis(config, freq)
	analysis.Print(*anal)
}
