package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	//Нахождение папки с файлами
	files, err := ioutil.ReadDir(".")

	if err != nil {
		fmt.Print("\nПапка с файлами не найдена\n")
		log.Fatal(err)
	}

	var records, newrecords, tworeplay, threereplay, check, numfile, wordlatters, tr int
	index := make(map[int]string)
	twoelem := make([]int, 20)
	threlem := make([]int, 20)
	twthrelem := make([]int, 20)

	//При нахождении папки, открываем ее и находим файлы
	for _, file := range files {
		if file.IsDir() == true {
			files2, err := ioutil.ReadDir("./" + file.Name())
			tr++
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print("\nВ папке были найдены файлы:\n")
			for _, file2 := range files2 {
				var namefolder string = file.Name()
				var namefile string = file2.Name()

				fmt.Print(namefolder)
				fmt.Print("/")
				fmt.Print(namefile)
				fmt.Print(";\n")
				//Открываем найденные файлы
				openfile, er := os.Open(namefolder + "/" + namefile)
				if er != nil {
					log.Fatal(er)
				}
				//Закрываем файлы
				defer openfile.Close()

				var elemcount int
				var allword string
				filetext := make([]byte, 100)
				word := make([]byte, 100)

				for {
					_, err := openfile.Read(filetext) //Считываем всю информацию из файла
					if err == io.EOF {
						break
					}

					//Записываем каждое слово через ' '
					for filetext[elemcount] != 0 {
						for count := 0; filetext[elemcount] != ' ' && filetext[elemcount] != 0; elemcount++ {
							word[count] = filetext[elemcount]
							count++
							wordlatters++
						}
						elemcount++

						for j := 0; j < wordlatters; j++ {
							allword = (allword + string(word[j]))
							word[j] = 0
						}
						wordlatters = 0
						//numfile обозначает нумерацию файла(за основу взято наличие трех)
						if numfile == 0 {
							//Составляем список слов
							for k := 0; k < records; k++ {
								if allword != index[k] {
									check = 1 //проверка на повторения
								}
							}
							if check == 1 {

								index[records] = allword
								check = 0
								records++
							}
							if records == 0 {
								index[records] = allword
								records++
							}

						}

						if numfile > 0 {
							for k := 0; k < records; k++ {
								newrecords++
								if allword != index[k] {
									check = 1
								} else {
									if numfile == 1 {
										twoelem[tworeplay] = k //Запоминаем места слов, которые есть и в 1-м и в 2-м файле
										tworeplay++
										check = 0
										break
									}
									if numfile == 2 {
										if k < records-newrecords && records != 0 {
											threlem[threereplay] = k //Запоминаем места слов, которые есть и в 1-м и в 3-м файле
											threereplay++
											check = 0
											break
										} else {
											twthrelem[threereplay] = k //Запоминаем места слов, которые есть и в 2-м и в 3-м файле
											threereplay++
											check = 0
											break
										}

									}

								}

							}
							if check == 1 {
								index[records] = allword
								check = 0
								records++
							}
							newrecords = 0

						}
						allword = ""
					}
					numfile++
				}
			}

		} else if tr != 1 {
			fmt.Print("\nНевозможно открыть папку с файлами\n")
		}
	}

	//Создаем файл для записи инвертированного индекса
	newfile, err := os.Create("invertindex.txt")
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}

	defer newfile.Close()

	check = 0

	for i := 0; i < records; i++ {

		for k := 0; k < tworeplay; k++ {
			if twoelem[k] == i {
				for c := 0; c < threereplay; c++ {
					if threlem[c] == i {
						newfile.WriteString(index[i])
						newfile.WriteString("  { 1, 2, 3 }  \n") //Одно и тоже слово встречается в 1-м, 2-м И в 1-м,3-м
						check = 1
						break
					}
				}

				if check == 0 {
					newfile.WriteString(index[i])
					newfile.WriteString("  { 1, 2 }  \n") //Одно и тоже слово встречается в 1-м, 2-м
					check = 1
				}
			}
		}

		if check == 0 {
			for k := 0; k < threereplay; k++ {
				if threlem[k] == i {
					newfile.WriteString(index[i])
					newfile.WriteString("  { 1, 3 }  \n") //Одно и тоже слово встречается в 1-м, 3-м
					check = 1
					break
				}
			}
		}

		if check == 0 {
			for k := 0; k < threereplay; k++ {
				if twthrelem[k] == i {
					newfile.WriteString(index[i])
					newfile.WriteString("  { 2, 3 }  \n") //Одно и тоже слово встречается во 2-м, 3-м
					check = 1
					break
				}
			}
		}

		if check == 0 {
			newfile.WriteString(index[i])
			newfile.WriteString("  { 1 }  \n") //Слово встречается один раз
		}
		check = 0
	}
	fmt.Print("\nИнвертированный индекс записан в созданный файл\n")
}
