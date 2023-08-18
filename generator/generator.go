package generator

var keys = [30]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", ",", ".", "'", "/"}

func GenColumns() [][3]string {
	var result [][3]string
	for i := 0; i < 30; i++ {
		for j := i + 1; j < 30; j++ {
			for k := j + 1; k < 30; k++ {
				if i != j && j != k && i != k {
					result = append(result, [3]string{keys[i], keys[j], keys[k]})
				}
			}
		}
	}
	return result
}
