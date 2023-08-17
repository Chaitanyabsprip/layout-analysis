package analysis

import (
	"fmt"
	"sort"
	"strings"
)

var (
	line    string = strings.Repeat("─", 10*3)
	topline string = "\t┌" + line + "┐"
	btmline string = "\t└" + line + "┘"
)

func keymapStr(l Layout) string {
	var result strings.Builder
	for i, row := range l.keymap {
		emptyrow := ""
		if i < 2 {
			emptyrow = fmt.Sprintf("\t│ %s  │\n", strings.Repeat(" ", 3*9))
		}
		rowstr := fmt.Sprintf("\t│ %s │\n%s", strings.Join(row[:], "  "), emptyrow)
		result.WriteString(rowstr)
	}
	return fmt.Sprintf("\n %s\n%s%s\n", topline, result.String(), btmline)
}

func fingerEffortStr(fingerEffort map[int]float64) string {
	fingers := []int{}
	for k := range fingerEffort {
		fingers = append(fingers, k)
	}
	sort.SliceStable(fingers, func(i, j int) bool {
		return fingers[i] < fingers[j]
	})
	var result strings.Builder
	for _, finger := range fingers {
		result.WriteString(fmt.Sprintf("finger %d     %.3f\n", finger, fingerEffort[finger]))
	}
	return result.String()
}

func topsfbStr(topsfbs map[string]float64) string {
	bigrams := []string{}
	for k := range topsfbs {
		bigrams = append(bigrams, k)
	}
	sort.SliceStable(bigrams, func(i, j int) bool {
		return topsfbs[bigrams[i]] > topsfbs[bigrams[j]]
	})
	var result strings.Builder
	for _, bigram := range bigrams {
		result.WriteString(fmt.Sprintf("%s     %.3f %%\n", bigram, topsfbs[bigram]))
	}
	return "\n" + result.String()
}

func underline(input string) string {
	under := "\x1b[4m"
	reset := "\x1b[0m"
	return under + input + reset
}

func bold(input string) string {
	bld := "\x1b[1m"
	reset := "\x1b[0m"
	return bld + input + reset
}

func Print(a Analysis) {
	var out strings.Builder
	str := out.WriteString
	line := func(ln string) { out.WriteString(ln + "\n") }
	line(keymapStr(a.Layout))
	line(underline("Finger Effort:"))
	line(fingerEffortStr(a.FingerEffort))
	str(underline("Total Effort Rating:") + " ")
	line(fmt.Sprintf(bold("%.3f"), a.EffortRating))
	line("")
	str(underline("Top Sfbs:"))
	line(topsfbStr(a.TopSfbs))
	str(underline("Sfb Rating:"))
	line(fmt.Sprintf(bold(" %.3f %%"), a.SfbRating))
	line("")
	str(underline("Inroll Rating:"))
	line(fmt.Sprintf(" %.3f", a.InrollRating))
	str(underline("Outroll Rating:"))
	line(fmt.Sprintf(" %.3f", a.OutrollRating))
	str(underline("Repeat Effort:"))
	line(fmt.Sprintf(" %.3f", a.RepeatEffort))
	fmt.Print(out.String())
}
