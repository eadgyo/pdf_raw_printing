package generator

const moduloID = 36

var accessLetter = []string{
	"0",
	"1",
	"2",
	"3",
	"4",
	"5",
	"6",
	"7",
	"8",
	"9",
	"A",
	"B",
	"C",
	"D",
	"E",
	"F",
	"G",
	"H",
	"I",
	"J",
	"K",
	"L",
	"M",
	"N",
	"O",
	"P",
	"Q",
	"R",
	"S",
	"T",
	"U",
	"V",
	"W",
	"X",
	"Y",
	"Z",
}

var lastIndexPerSuffix = map[string]int{}
var alreadyKnown = map[string]bool{}

func Register(s string) {
	alreadyKnown[s] = true
}

func gen(i int) string {
	s := ""
	power := 1

	for i != 0 {
		intermediate := i / power
		rest := intermediate % moduloID
		i -= power * rest
		power *= moduloID

		s = accessLetter[rest] + s
	}

	return s
}

func GetSize() int {
	return len(alreadyKnown)
}

func Generate(prefix string) string {
	var i int = 0
	var ok bool
	if i, ok = lastIndexPerSuffix[prefix]; !ok {
		lastIndexPerSuffix[prefix] = i
	}

	i++
	lastIndexPerSuffix[prefix] = i
	s := prefix + gen(i)
	if _, ok := alreadyKnown[s]; ok {
		return Generate(prefix)
	}
	alreadyKnown[s] = true
	return s
}
