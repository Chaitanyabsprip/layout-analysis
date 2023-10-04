package gen

var keys = [30]string{
	"a",
	"b",
	"c",
	"d",
	"e",
	"f",
	"g",
	"h",
	"i",
	"j",
	"k",
	"l",
	"m",
	"n",
	"o",
	"p",
	"q",
	"r",
	"s",
	"t",
	"u",
	"v",
	"w",
	"x",
	"y",
	"z",
	",",
	".",
	"'",
	"/",
}

func GenColumns() map[string]int {
	result := make(map[string]int)
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			for k := j + 1; k < len(keys); k++ {
				combination := keys[i] + keys[j] + keys[k]
				result[combination] = (1 << i) | (1 << j) | (1 << k)
			}
		}
	}
	return result
}

func GenLayouts() [][3][10]string {
	layouts := [][3][10]string{}
	return layouts
}
