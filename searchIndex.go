//Открытие файла с инвертированным индексом осуществляется через os.Args[2]
//Слова для поиска сичтываются с консоли

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"
)

func CoincidencesIndexes(allIndexes []byte, countBytes int, numberWords int) string {
	var n, w int
	var inputstr string
	firstnumber := allIndexes[0]
	allIndexes[0] = 0
	n = 1

	for w < countBytes {
		for count := 0; count < countBytes; count++ {
			if allIndexes[count] == 32 || allIndexes[count] == 0 {
				allIndexes[count] = 0
			} else if allIndexes[count] == firstnumber {
				n++
				allIndexes[count] = 0
			}
		}
		if n != numberWords {
		} else {
			inputstr = (inputstr + " " + string(firstnumber))
		}
		n = 0
		for count := 0; count < countBytes; count++ {
			if allIndexes[count] == 32 || allIndexes[count] == 0 {
			} else {
				firstnumber = allIndexes[count]
				allIndexes[count] = 0
				n = 1
				break
			}
		}
		w++
	}
	return inputstr
}

func indexSearch() {
	fileWithInd := os.Args[2]
	var sentence, inputstr string
	var coincidencesWords, countBytes int
	allIndexes := make([]byte, 100)
	index := make([]byte, 100)

	openfile, err := os.Open(fileWithInd)
	if err != nil {
		log.Fatal(err)
	}
	defer openfile.Close()

	byteFile, err := ioutil.ReadFile(fileWithInd)
	if err != nil {
		log.Fatal(err)
	}
	stringFile := string(byteFile)
	isLetter := func(c rune) bool {
		return !unicode.IsLetter(c)
	}
	fileWords := strings.FieldsFunc(stringFile, isLetter)

	stringFile = strings.Replace(stringFile, " ", "", -1)
	isNumber := func(c rune) bool {
		return !unicode.IsNumber(c)
	}
	indexFileNumbers := strings.FieldsFunc(stringFile, isNumber)

	fmt.Println("\nВведите сочетание слов для поиска: ")
	in := bufio.NewReader(os.Stdin)

	sentence, err = in.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка ввода: ", err)
	}

	lowSenttext := strings.ToLower(sentence)
	searchWords := strings.FieldsFunc(lowSenttext, isLetter)

	for n := 0; n < len(searchWords); n++ {
		for k := 0; k < len(fileWords); k++ {
			if fileWords[k] == searchWords[n] {
				index = []byte(indexFileNumbers[k])
				for m := 0; m < len(indexFileNumbers[k]); m++ {
					allIndexes[countBytes] = index[m]
					countBytes++
				}
				coincidencesWords++
			}
		}
	}

	if coincidencesWords < len(searchWords) {
		inputstr = ("Введенные слова не встречаются вместе")
	} else {
		inputstr = CoincidencesIndexes(allIndexes, countBytes, len(searchWords))
		if inputstr != "" {
			inputstr = ("Введенные слова были найдены в текстах:" + inputstr + "\n")
		} else {
			inputstr = ("Введенные слова не встречаются вместе")
		}
	}
	fmt.Println(inputstr)
}
