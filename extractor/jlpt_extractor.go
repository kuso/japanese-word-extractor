package extractor

import (
	"github.com/ikawaha/kagome/tokenizer"
	"log"
	"strconv"
	"strings"
	"github.com/kuso/japanese-word-extractor/helper"
)

func NewJLPTDictionary() (*JLPTDictionary, error) {
	dict := JLPTDictionary{}
	dict.WordLevelMap = make(map[string]int, 10000)
	dict.EnglishMeanings = make(map[string]string, 10000)


	n1Lines, err := helper.ReadLines("./data", "N1-Vocab.csv")
	if err != nil {
		return nil, err
	}
	for _, line := range n1Lines {
		parts := strings.Split(line, ",")
		var key string
		if parts[0] == "" {
			key = parts[1] + "-" + parts[1]
		} else {
			key = parts[1] + "-" + parts[0]
		}
		dict.WordLevelMap[key] = 1
		dict.EnglishMeanings[key] = parts[2]
	}

	n2Lines, err := helper.ReadLines("./data", "N2-Vocab.csv")
	if err != nil {
		return nil, err
	}
	for _, line := range n2Lines {
		parts := strings.Split(line, ",")
		var key string
		if parts[0] == "" {
			key = parts[1] + "-" + parts[1]
		} else {
			key = parts[1] + "-" + parts[0]
		}
		dict.WordLevelMap[key] = 2
		dict.EnglishMeanings[key] = parts[2]
	}

	n3Lines, err := helper.ReadLines("./data", "N3-Vocab.csv")
	if err != nil {
		return nil, err
	}
	for _, line := range n3Lines {
		parts := strings.Split(line, ",")
		var key string
		if parts[0] == "" {
			key = parts[1] + "-" + parts[1]
		} else {
			key = parts[1] + "-" + parts[0]
		}
		dict.WordLevelMap[key] = 3
		dict.EnglishMeanings[key] = parts[2]
	}

	n4Lines, err := helper.ReadLines("./data", "N4-Vocab.csv")
	if err != nil {
		return nil, err
	}
	for _, line := range n4Lines {
		parts := strings.Split(line, ",")
		var key string
		if parts[0] == "" {
			key = parts[1] + "-" + parts[1]
		} else {
			key = parts[1] + "-" + parts[0]
		}
		dict.WordLevelMap[key] = 4
		dict.EnglishMeanings[key] = parts[2]
	}

	n5Lines, err := helper.ReadLines("./data", "N5-Vocab.csv")
	if err != nil {
		return nil, err
	}
	for _, line := range n5Lines {
		parts := strings.Split(line, ",")
		var key string
		if parts[0] == "" {
			key = parts[1] + "-" + parts[1]
		} else {
			key = parts[1] + "-" + parts[0]
		}
		dict.WordLevelMap[key] = 5
	}

	return &dict, nil
}

func GetTokens(input string, dict *JLPTDictionary) []*JLPTToken {
	t := tokenizer.New()
	//tokens := t.Tokenize("寿司が食べたい。") // t.Analyze("寿司が食べたい。", tokenizer.Normal)
	tokens := t.Tokenize(strings.TrimSpace(input)) // t.Analyze("寿司が食べたい。", tokenizer.Normal)

	jTokens := make([]*JLPTToken, 0)
	for _, token := range tokens {
		if token.Class == tokenizer.DUMMY {
			// BOS: Begin Of Sentence, EOS: End Of Sentence.
			//fmt.Printf("%s\n", token.Surface)
			continue
		}

		jToken := JLPTToken{}
		jToken.Text = token.Surface
		jToken.DictForm = jToken.Text
		jTokens = append(jTokens, &jToken)

		features := token.Features()
		if len(features) >= 9 {
			jToken.DictForm = features[6]
			jToken.DictFormHiragana = helper.Kata2hira(features[7])
		}

		key := jToken.DictFormHiragana + "-" + jToken.DictForm
		level, ok := dict.WordLevelMap[key]
		if ok {
			jToken.Level = level
		}

		meaning, ok := dict.EnglishMeanings[key]
		if ok {
			jToken.Meaning = strings.Trim(meaning, "\"")
		}
		//features := strings.Join(token.Features(), ",")
		//fmt.Printf("%s\t%v\n", token.Surface, features)
	}
	return jTokens
}

func PrintTokens(tokens []*JLPTToken) {
	for _, token := range tokens {
		log.Println(token.Text, token.DictForm, token.Level)
	}
}

func OutputHTML(tokens []*JLPTToken) string {
	outStr := ""
	for _, token := range tokens {
		if token.Level > 0 {
			outStr = outStr + "<span class=\"jlptn" + strconv.Itoa(token.Level) + "\">" + token.Text + "</span>"
		} else {
			outStr = outStr + token.Text
		}
	}
	return outStr
}
