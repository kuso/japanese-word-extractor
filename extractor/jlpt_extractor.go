package extractor

import (
	"bytes"
	"compress/gzip"
	"errors"
	"github.com/gojp/nihongo/lib/dictionary"
	"github.com/ikawaha/kagome/tokenizer"
	"github.com/kuso/japanese-word-extractor/helper"
	"golang.org/x/text/width"
	"log"
	"os"
	"strconv"
	"strings"
)

func NewJLPTDictionary() (*JLPTDictionary, error) {
	dict := JLPTDictionary{}
	dict.WordLevelMap = make(map[string]int, 10000)
	dict.EnglishMeanings = make(map[string]string, 10000)

	workingDir := helper.Basepath
	dataDir := workingDir + "/data"
	n1Lines, err := helper.ReadLines(dataDir, "N1-Vocab.csv")
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

	n2Lines, err := helper.ReadLines(dataDir, "N2-Vocab.csv")
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

	n3Lines, err := helper.ReadLines(dataDir, "N3-Vocab.csv")
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

	n4Lines, err := helper.ReadLines(dataDir, "N4-Vocab.csv")
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

	n5Lines, err := helper.ReadLines(dataDir, "N5-Vocab.csv")
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

func NewEdict2Dictionary() (*dictionary.Dictionary, error) {
	file, err := os.Open(helper.Basepath + "/data/edict2.json.gz")
	if err != nil {
		log.Fatal("Could not load edict2.json.gz: ", err)
		return nil, err
	}
	defer file.Close()

	reader, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal("Could not create reader: ", err)
		return nil, err
	}

	dict, err := dictionary.Load(reader)
	if err != nil {
		log.Fatal("Could not load dictionary: ", err)
		return nil, err
	}
	return &dict, nil
}

func GetDefinition(entry *dictionary.Entry) string {
	var defs []string
	for _, g := range entry.Glosses {
		defs = append(defs, g.English)
	}
	if len(defs) > 0 {
		return strings.Join(defs, "; ")
	}
	return ""
}

func GetDictForm(t tokenizer.Tokenizer, token tokenizer.Token) (string, string, error) {

	// 0: word class
	// 6: dictionary form with kanji
	// 7: katagana
	features := token.Features()
	if len(features) >= 9 {
		wordClass := features[0]
		formClass := features[5]

		// need to find actual dictionary form
		if wordClass == "動詞" && formClass != "基本形" {
			tmpTokens := t.Tokenize(strings.TrimSpace(features[6]))
			if len(tmpTokens) != 3 {
				return "", "", errors.New("error getting dictionary form, " + features[6])
			}
			tmpToken := tmpTokens[1]
			actualDictForm := tmpToken.Features()[6]
			actualDictFormHiragana := helper.Kata2hira(tmpToken.Features()[7])
			return actualDictForm, actualDictFormHiragana, nil
		}
	}
	return "", "", nil
}

func GetTokens(input string, dict *JLPTDictionary, dict2 *dictionary.Dictionary) []*JLPTToken {
	t := tokenizer.New()
	//tokens := t.Tokenize("寿司が食べたい。") // t.Analyze("寿司が食べたい。", tokenizer.Normal)
	tokens := t.Tokenize(strings.TrimSpace(input)) // t.Analyze("寿司が食べたい。", tokenizer.Normal)

	jTokens := make([]*JLPTToken, 0)
	for _, token := range tokens {
		if token.Class == tokenizer.DUMMY {
			// BOS: Begin Of Sentence, EOS: End Of Sentence.
			continue
		}

		s := bytes.Buffer{}
		for _, item := range token.Features() {
			s.WriteString(item)
			s.WriteString(", ")
		}

		jToken := JLPTToken{}
		jToken.Text = token.Surface
		jToken.DictForm = jToken.Text
		jTokens = append(jTokens, &jToken)

		// convert fullwidth to halfwidth digit and skip if text is number
		narrowWidth := width.Narrow.String(strings.TrimSpace(jToken.DictForm))
		if _, err := strconv.Atoi(narrowWidth); err == nil {
			jToken.Text = narrowWidth
			jToken.DictForm = narrowWidth
			continue
		}

		// 6: dictionary form with kanji
		features := token.Features()
		if len(features) >= 9 {
			jToken.Class = features[0]
			if jToken.Class == "動詞" && features[5] != "基本形" {
				actualDictForm, actualDictFormHiragana, err := GetDictForm(t, token)
				if err != nil {
					log.Println("error")
					log.Println(err)
					continue
				}
				jToken.DictForm = actualDictForm
				jToken.DictFormHiragana = actualDictFormHiragana
			} else {
				jToken.DictForm = features[6]
				jToken.DictFormHiragana = helper.Kata2hira(features[7])
			}
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

		if !ok && dict2 != nil {
			results := dict2.Search(jToken.DictForm, 10)
			for _, result := range results {
				defStr := GetDefinition(&result)
				jToken.Meaning = defStr
				break
			}
		}
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
