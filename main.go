package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	lenToWord, err := getWordSet()
	if err != nil {
		panic(err)
	}
	br := bufio.NewReader(os.Stdin)
	currString := ""
	for {
		r, _, err := br.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if !unicode.IsLetter(r) {
			if len(currString) == 0 {
				_, err = os.Stdout.WriteString(string(r))
				if err != nil {
					panic(err)
				}
				currString = ""
				continue
			}
			var dogWord string
			words, ok := lenToWord[len(currString)]
			if !ok {
				dogWord = combineMultipleLengths(lenToWord, len(currString))
			} else {
				dogWord = pickRandom(words)
			}
			dogWord = matchCase(currString, dogWord)
			_, err = os.Stdout.WriteString(dogWord + string(r))
			if err != nil {
				panic(err)
			}
			currString = ""
			continue
		}
		currString += string(r)
	}
}

func matchCase(model string, s string) string {
	if len(model) == 0 || len(s) == 0 {
		return s
	}
	newStr := ""
	for i, v := range []rune(model) {
		if len(s) > i {
			if unicode.IsUpper(v) {
				newStr += string(unicode.ToUpper([]rune(s)[i]))
			} else {
				newStr += string(s[i])
			}
		} else {
			break
		}
	}
	return newStr
}

func getWordSet() (map[int][]string, error) {
	data, err := ioutil.ReadFile("wordset.txt")
	if err != nil {
		return nil, err
	}
	wordset := strings.Fields(string(data))

	lenToWord := make(map[int][]string, len(wordset))
	for i := range lenToWord {
		lenToWord[i] = make([]string, 0, 5)
	}

	for _, v := range wordset {
		lenToWord[len(v)] = append(lenToWord[len(v)], v)
	}
	return lenToWord, nil
}

func pickRandom(s []string) string {
	return s[rand.Intn(len(s))]
}

func combineMultipleLengths(wordset map[int][]string, i int) string {
	max := 0
	for k := range wordset {
		if k > max && k <= i { // get max k such that it's under i
			max = k
		}
	}
	currInt := max
	result := ""
	for {
		if len(result) == i {
			break
		}
		if i-len(result) < currInt {
			currInt = i - len(result)
		}
		if v, ok := wordset[currInt]; ok {
			result += pickRandom(v)
			for i-len(result) >= currInt {
				result += pickRandom(v)
			}
		} else {
			currInt--
		}
	}
	return result
}
