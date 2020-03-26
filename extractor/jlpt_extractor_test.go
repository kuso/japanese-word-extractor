package extractor

import (
	"log"
	"testing"
)

func Test_JLPTDictionary_GetTokens(t *testing.T) {
	dict, err := NewJLPTDictionary()
	if err != nil {
		t.Fatal(err)
	}

	// first 8 is half width digit, second 8 is full widht
	testStr := "8８, 狙う, 狙っている"
	tokens := GetTokens(testStr, dict, nil)
	for _, token := range tokens {
		log.Println(token.Text + "," + token.DictForm)
	}
}

func Test_Edict2Dictionary_(t *testing.T) {
	dict, err := NewEdict2Dictionary()
	if err != nil {
		t.Fatal(err)
	}

	text := "新型"
	results := dict.Search(text, 10)
	for _, result := range results {
		defStr := GetDefinition(&result)
		log.Println(defStr)
		break
	}
}
