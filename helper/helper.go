package helper

import (
	"bufio"
	"os"
	"strings"
)

func hira2kata (hira rune) rune {
	if (hira >= 'ぁ' && hira <= 'ゖ') || (hira >= 'ゝ' && hira <= 'ゞ') {
		return hira + 0x60
	}
	return hira
}

func Hira2kata(hira string) string {
	return strings.Map(hira2kata, hira)
}

func kata2hira (kata rune) rune {
	if (kata >= 'ァ' && kata <= 'ヶ') || (kata >= 'ヽ' && kata <= 'ヾ') {
		return kata - 0x60
	}
	return kata
}

func Kata2hira(kata string) string {
	return strings.Map(kata2hira, kata)
}

func ReadLines(baseDir string, fileName string) ([]string, error) {
	f, err := os.Open(baseDir + "/" + fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return lines, nil
}
