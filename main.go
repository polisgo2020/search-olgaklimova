package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	files, err := ioutil.ReadDir(".")

	if err != nil {
		log.Fatal(err)
	}

	var h int
	var newh int
	var i int
	var y int
	var t int
	var prov int
	var numfile int
	index := make(map[int]string)
	twoelem := make([]int, 20)
	threlem := make([]int, 20)
	twthrelem := make([]int, 20)

	for _, file := range files {
		if file.IsDir() == true {
			files2, err := ioutil.ReadDir("./" + file.Name())
			if err != nil {
				log.Fatal(err)
				fmt.Print("\nНевозможно открыть папку с файлами")
			}
			for _, file2 := range files2 {
				var b string = file.Name()
				var a string = file2.Name()
				fmt.Print("\n")
				fmt.Print("Файл ")
				fmt.Print(b)
				fmt.Print("/")
				fmt.Print(a)
				fmt.Print(": ")
				openfile, er := os.Open(b + "/" + a)
				if er != nil {
					log.Fatal(er)
					fmt.Print("Не получается открыть файл")
					continue
				}
				defer openfile.Close()

				var elemcount int
				var allword string
				data := make([]byte, 64)
				word := make([]byte, 64)

				for {
					_, err := openfile.Read(data)
					if err == io.EOF {
						break
					}

					for data[elemcount] != 0 {
						for count := 0; data[elemcount] != ' ' && data[elemcount] != 0; elemcount++ {
							word[count] = data[elemcount]
							count++
							i++
						}
						elemcount++

						for j := 0; j < i; j++ {
							allword = (allword + string(word[j]))
							word[j] = 0
						}

						i = 0
						if numfile == 0 {

							for k := 0; k < h; k++ {
								if allword != index[k] {
									prov = 1
								}
							}
							if prov == 1 {

								index[h] = allword
								prov = 0
								h++
							}
							if h == 0 {
								index[h] = allword
								h++
							}

						}

						if numfile > 0 {
							for k := 0; k < h; k++ {
								newh++
								if allword != index[k] {
									prov = 1
								} else {
									if numfile == 1 {
										twoelem[y] = k
										y++
										prov = 0
										break
									}
									if numfile == 2 {
										if k < h-newh {
											threlem[t] = k
											t++
											prov = 0
											break
										} else {
											twthrelem[t] = k
											t++
											prov = 0
											break
										}

									}

								}

							}
							if prov == 1 {
								index[h] = allword
								prov = 0
								h++
							}
							newh = 0

						}
						allword = ""
						// fmt.Print(index[h])
						// fmt.Println("\n")
					}
					numfile++

					// fmt.Print(string(data[:n]))
					// fmt.Print("\n")
				}
			}

		}
	}

	newfile, err := os.Create("invertindex.txt")
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}

	defer newfile.Close()

	prov = 0

	for i := 0; i < h; i++ {

		for k := 0; k < y; k++ {
			if twoelem[k] == i {
				for c := 0; c < t; c++ {
					if threlem[c] == i {
						newfile.WriteString(index[i])
						newfile.WriteString("  { 1, 2, 3 }  \n")
						prov = 1
						break
					}
				}

				if prov == 0 {
					newfile.WriteString(index[i])
					newfile.WriteString("  { 1, 2 }  \n")
					prov = 1
				}
			}
		}

		if prov == 0 {
			for k := 0; k < t; k++ {
				if threlem[k] == i {
					newfile.WriteString(index[i])
					newfile.WriteString("  { 1, 3 }  \n")
					prov = 1
					break
				}
			}
		}

		if prov == 0 {
			for k := 0; k < t; k++ {
				if twthrelem[k] == i {
					newfile.WriteString(index[i])
					newfile.WriteString("  { 2, 3 }  \n")
					prov = 1
					break
				}
			}
		}

		if prov == 0 {
			newfile.WriteString(index[i])
			newfile.WriteString("  { 1 }  \n")
		}

		prov = 0

	}

}
