package analysis

import "errors"

type Layout struct {
	keymap    [3][10]string
	fingermap [3][10]int
}

func verifyKeymap(arr [3][10]string) bool {
	seen := make(map[string]bool)
	for _, row := range arr {
		for _, val := range row {
			if len(val) != 1 || seen[val] {
				return false
			}
			seen[val] = true
		}
	}
	return true
}

func NewLayout(keymap [3][10]string, fingermap [3][10]int) (*Layout, error) {
	if verifyKeymap(keymap) {
		return &(Layout{keymap: keymap, fingermap: fingermap}), nil
	}
	return nil, errors.New("invalid keymap: either a repeated key or more than one character key detected")
}

func (l *Layout) FingerKeymap() map[int][]string {
	fingermap := make(map[int][]string)
	for i := 0; i < 3; i++ {
		for j := 0; j < 10; j++ {
			finger := l.fingermap[i][j]
			key := l.keymap[i][j]
			if _, ok := fingermap[finger]; !ok {
				fingermap[finger] = []string{key}
			} else {
				fingermap[finger] = append(fingermap[finger], key)
			}
		}
	}
	return fingermap
}

func (l *Layout) Left() [3][5]string {
	var left [3][5]string
	for i, row := range l.keymap {
		half := len(row) / 2
		leftRow := make([]string, half)
		copy(leftRow, row[:half])
		copy(left[i][:half], leftRow)
	}
	return left
}

func (l *Layout) Right() [3][5]string {
	var right [3][5]string
	for i, row := range l.keymap {
		half := len(row) / 2
		rightRow := make([]string, half)
		copy(rightRow, row[half:])
		copy(right[i][:half], rightRow)
	}
	return right
}

func (l *Layout) Inrolls() []string {
	var inrolls []string
	leftHandInrolls := make([]string, 0)
	for _, row := range l.Left() {
		for i := 0; i < len(row)-1; i++ {
			leftHandInrolls = append(leftHandInrolls, row[i]+row[i+1])
		}
	}
	rightHandInrolls := make([]string, 0)
	for _, row := range l.Right() {
		for i := 0; i < len(row)-1; i++ {
			rightHandInrolls = append(rightHandInrolls, row[i+1]+row[i])
		}
	}
	inrolls = append(leftHandInrolls, rightHandInrolls...)
	return inrolls
}

func (l *Layout) Outrolls() []string {
	inrolls := l.Inrolls()
	outrolls := []string{}
	for _, inroll := range inrolls {
		outroll := []byte(inroll)
		outroll[0], outroll[1] = outroll[1], outroll[0]
		outrolls = append(outrolls, string(outroll))
	}
	return outrolls
}

func (l *Layout) Sfbs() []string {
	sfbs := []string{}
	for _, keys := range l.FingerKeymap() {
		for _, key1 := range keys {
			for _, key2 := range keys {
				if key1 != key2 {
					sfbs = append(sfbs, key1+key2)
				}
			}
		}
	}
	return sfbs
}
