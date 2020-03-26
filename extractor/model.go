package extractor

type JLPTDictionary struct {
	WordLevelMap    map[string]int
	EnglishMeanings map[string]string
}

type JLPTToken struct {
	Text             string
	DictForm         string
	DictFormHiragana string
	Level            int
	Meaning          string
	Class            string
}
