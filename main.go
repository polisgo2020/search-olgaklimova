package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	//Нахождение папки с файлами

	papka := os.Args[1]

	dir, err := os.Open(papka)
	if err != nil {
		return
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return
	}

	fmt.Println("\nВ папке были найдены файлы:")
	var nfile, records, newrecords, wordlatters, check, elemcount int
	var allword, indfi string
	index := make(map[int]string)
	indexfile := make(map[int]string)
	numberfile := make(map[int]int)

	for _, fi := range fileInfos {

		fmt.Println(fi.Name())
		nfile++
		openfile, err := os.Open(papka + "/" + fi.Name())
		if err != nil {
			log.Fatal(err)
		}
		defer openfile.Close()

		filetext := make([]byte, 100000000)
		word := make([]byte, 100)
		for {
			_, err := openfile.Read(filetext) //Считываем всю информацию из файла
			if err == io.EOF {
				break
			}
		}

		//Проверка считанного текста по элементам
		for filetext[elemcount] != 0 {
			for count := 0; filetext[elemcount] != ' ' && filetext[elemcount] != 0 && filetext[elemcount] != '"' && filetext[elemcount] != '\n' && filetext[elemcount] != ':' && filetext[elemcount] != '.' && filetext[elemcount] != 39 && filetext[elemcount] != 45 && filetext[elemcount] != 150 && filetext[elemcount] != 151; elemcount++ {
				if (filetext[elemcount] < 65) || filetext[elemcount] > 122 || (filetext[elemcount] < 97 && filetext[elemcount] > 90) {
				} else {
					if filetext[elemcount] > 64 && filetext[elemcount] < 91 {
						filetext[elemcount] = filetext[elemcount] + 32
					}
					word[count] = filetext[elemcount]
					count++
					wordlatters++
				}
			}
			elemcount++

			//Запись каждого слова
			for j := 0; j < wordlatters; j++ {

				allword = (allword + string(word[j]))
				word[j] = 0
			}

			//Запись слова в инвертированный индекс
			for k := 0; k <= records; k++ {
				if allword == index[k] {
					if nfile == 1 {
						numberfile[k] = nfile
						indfi = indexfile[records]
					} else {
						if k < records-newrecords {
							indfi = indexfile[k]
							nf := strconv.Itoa(nfile)
							if numberfile[k] != nfile {
								indexfile[k] = (indfi + nf + " ")
							}
							numberfile[k] = nfile
						}
					}
					check = 1
				}
			}

			if check == 0 {
				for k := 0; k <= records; k++ {
					index[records] = allword
					check = 0
					nf := strconv.Itoa(nfile)
					indexfile[records] = (" " + nf + " ")
				}
				newrecords++
				records++
			}

			wordlatters = 0
			allword = ""
			check = 0
		}
		newrecords = 0
		elemcount = 0
	}

	//Создание файла и запись полученного индекса

	newfile, err := os.Create("invertindex.txt")
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	for k := 0; k < records; k++ {
		newfile.WriteString(index[k])
		newfile.WriteString(" {")
		newfile.WriteString(indexfile[k])
		newfile.WriteString("};\n")
	}

	defer newfile.Close()
	fmt.Println("\nФайл с инвертированным индексом записан")
}
