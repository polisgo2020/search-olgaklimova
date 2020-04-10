package index

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

var index = make(map[int]string)
var indexfile = make(map[int]string)
var numberfile = make(map[int]int)
var mu = &sync.Mutex{}
var records int
var newrecords int
var n int

func TextAnalysis(alltext chan string, nfile int, wg *sync.WaitGroup) *sync.WaitGroup {
	defer wg.Done()
LOOP2:
	for {
		select {
		case filetext := <-alltext:
			lowFiletext := strings.ToLower(filetext) //приведение к нижнему регистру

			f := func(c rune) bool {
				return !unicode.IsLetter(c)
			}
			allwords := strings.FieldsFunc(lowFiletext, f)
			for k := 0; k < len(allwords); k++ {
				СreateIndex(nfile, allwords[k])
				runtime.Gosched()
			}
		default:
			break LOOP2
		}
		break LOOP2
	}
	return wg
}

func СreateIndex(nfile int, allword string) string {
	mu.Lock()
	var check int
	var indfi string
	for k := 0; k <= records; k++ {
		if allword == index[k] {
			if k < records {
				indfi = indexfile[k]
				nf := strconv.Itoa(nfile)
				if strings.Contains(indexfile[k], nf) != true {
					indexfile[k] = (indfi + nf + " ")
				}
			}
			check = 1
		}
	}
	if check == 0 {
		index[records] = allword
		nf := strconv.Itoa(nfile)
		indexfile[records] = (" " + nf + " ")
		records++
	}
	check = 0
	mu.Unlock()
	rec := records - 1
	return (index[rec])
}

func WriteIndex() string {
	fileWithInd := os.Args[2]
	newfile, err := os.Create(fileWithInd)
	// fileWithInd := ("C:\\Users\\asus\\source\\repos\\hello\\invertindex.txt")
	// newfile, err := os.Create(fileWithInd)
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}

	for k := 0; k < records; k++ {
		mu.Lock()
		newfile.WriteString(filepath.Join(index[k] + " {" + indexfile[k] + "}\n"))
		mu.Unlock()
	}
	defer newfile.Close()
	return ("\nФайл с инвертированным индексом записан")
}

//coincidencesIndexes: в allIndexes []byte лежат все номера файлов в байтовом представлении
//где встречаются введенные слова
//в firstnumber кладем первый элемент массива
//идем по массиву,если находим тот же элемент allIndexes[count] == firstnumber,
//то в массиве он абнуляется allIndexes[count] = 0 и n++, где n - счетчик совпадений номеров файлов
//если колличество совпадений равно колличеству слов numberWords,
// то inputstr = (inputstr + " " + string(firstnumber))
//иначе кладем в firstnumber новое значение и начинаем цикл сначала
func CoincidencesIndexes(allIndexes []byte, countBytes int, numberWords int) string {
	var n, w int
	var inputstr string
	firstnumber := allIndexes[0]
	allIndexes[0] = 0
	n = 1
	for w < countBytes {
		for count := 0; count < countBytes; count++ {
			if allIndexes[count] == 0 {
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
			if allIndexes[count] == 0 {
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

func IndexSearch(stringFile string, sentence string) string {
	var outputstr string
	var coincidencesWords, countBytes int
	allIndexes := make([]byte, 100)
	index := make([]byte, 100)

	isLetter := func(c rune) bool {
		return !unicode.IsLetter(c)
	}
	fileWords := strings.FieldsFunc(stringFile, isLetter)

	stringFile = strings.Replace(stringFile, " ", "", -1)
	isNumber := func(c rune) bool {
		return !unicode.IsNumber(c)
	}
	indexFileNumbers := strings.FieldsFunc(stringFile, isNumber)

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
		outputstr = ("Введенные слова не встречаются вместе")
	} else {
		outputstr = CoincidencesIndexes(allIndexes, countBytes, len(searchWords))
		if outputstr != "" {
			outputstr = ("Введенные слова были найдены в текстах:" + outputstr + "\n")
		} else {
			outputstr = ("Введенные слова не встречаются вместе")
		}
	}
	return outputstr
}
