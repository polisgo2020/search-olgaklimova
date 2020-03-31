//Открытие файла с инвертированным индексом осуществляется через os.Args[2]
//Слова для поиска сичтываются с консоли

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func indexSearch() {

	file_with_ind := os.Args[2]
	// var file_with_ind string
	// file_with_ind = "C:\\Users\\asus\\source\\repos\\hello\\invertindex.txt"

	var allword, sentence, indstr, simind string
	var k, y, m, elemcount, elem, wordlatters, length, indlatters int
	search_words := make(map[int]string)
	indtext := make(map[int][]byte)
	senttext := make([]byte, 1000)
	word := make([]byte, 100)
	filetext := make([]byte, 1000000)
	indword := make([]byte, 100)
	indexes := make([]byte, 100)
	index := make([]byte, 100)
	allindexes := make(map[string]string)

	openfile, err := os.Open(file_with_ind)
	if err != nil {
		log.Fatal(err)
	}
	defer openfile.Close()

	fmt.Println("\nВведите сочетание слов для поиска: ")

	in := bufio.NewReader(os.Stdin)

	sentence, err = in.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка ввода: ", err)
	}

	senttext = []byte(sentence)
	length = len(senttext)
	for elemcount < length {
		for count := 0; elemcount < len(senttext) && senttext[elemcount] != ' ' && senttext[elemcount] != '"' && senttext[elemcount] != '\n' && senttext[elemcount] != ':' && senttext[elemcount] != '.' && senttext[elemcount] != 39 && senttext[elemcount] != 45 && senttext[elemcount] != 150 && senttext[elemcount] != 151; elemcount++ {
			if (senttext[elemcount] < 65) || senttext[elemcount] > 122 || (senttext[elemcount] < 97 && senttext[elemcount] > 90) {
			} else {
				if senttext[elemcount] > 64 && senttext[elemcount] < 91 {
					senttext[elemcount] = senttext[elemcount] + 32
				}
				word[count] = senttext[elemcount]
				count++
				wordlatters++
			}
		}
		elemcount++

		//Запись каждого слова, считанного с консоли
		for j := 0; j < wordlatters; j++ {
			allword = (allword + string(word[j]))
			word[j] = 0
		}

		search_words[k] = allword
		k++
		allword = ""
		wordlatters = 0
	}

	for {
		_, err := openfile.Read(filetext) //Считываем всю информацию из файла
		if err == io.EOF {
			break
		}
	}

	for filetext[elem] != 0 {
		for count := 0; filetext[elem] != ' ' && filetext[elem] != 0 && filetext[elem] != '\n'; elem++ {
			if filetext[elem] < 65 || filetext[elem] > 122 || (filetext[elem] < 97 && filetext[elem] > 90) {
			} else {
				indword[count] = filetext[elem]
				count++
				wordlatters++
			}
		}
		elem++

		//Запись каждого слова, считанного из файла
		for j := 0; j < wordlatters; j++ {
			simind = (simind + string(indword[j]))
			indword[j] = 0
		}

		for n := 0; n < k; n++ {
			if simind == search_words[n] {
				elem++
				for count := 0; filetext[elem] != '}' && filetext[elem] != 0; elem++ {
					index[count] = filetext[elem]
					count++
					indlatters++
				}
				elem++
				for j := 0; j < indlatters; j++ {
					indstr = (indstr + string(index[j]))
					//Запись всех индексов текстов в массив байт
					indexes[y] = index[j]
					y++
					index[j] = 0
				}
				//Создание map allindexes для хранения слова и номеров текстов, где оно встречается
				allindexes[simind] = indstr
				//Создание map allindexes для хранения слова и номеров(в байтовов представлении) текстов, где оно встречается
				indtext[m] = []byte(indstr)
				m++
			}
		}
		indstr = ""
		simind = ""
		wordlatters = 0
		indlatters = 0
	}

	var n, w int
	var firstnumber byte
	firstnumber = indexes[0]
	indexes[0] = 0
	n = 1
	//Цикл for, осуществляющий нахождение совпадений в массиве indexes[]
	for w < y {
		for count := 0; count < y; count++ {
			if indexes[count] == 32 || indexes[count] == 0 {
				indexes[count] = 0
			} else {
				if indexes[count] == firstnumber {
					n++
					indexes[count] = 0
				}
			}
		}
		if n != m {
		} else {
			indstr = (indstr + " " + string(firstnumber))
		}
		n = 0
		for count := 0; count < y; count++ {
			if indexes[count] == 32 || indexes[count] == 0 {
			} else {
				firstnumber = indexes[count]
				indexes[count] = 0
				n = 1
				break
			}
		}
		w++
	}

	if indstr != "" {
		indstr = ("Введенные слова были найдены в текстах:" + indstr + "\n")
	} else {
		indstr = ("Введенные слова не встречаются вместе")
	}

	fmt.Println(indstr)
}
