//os.Args[1] - путь к папке
//os.Args[2] - путь к файлу для записи индекса
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var index = make(map[int]string)
var indexfile = make(map[int]string)
var numberfile = make(map[int]int)
var mu = &sync.Mutex{}
var records int
var newrecords int

func openfile(fname chan string, nfile int) {
	var alltext chan []byte = make(chan []byte)
	filetext := make([]byte, 100000000)
LOOP:
	for {
		select {
		case str := <-fname:
			openfile, err := os.Open(str)
			if err != nil {
				log.Fatal(err)
			}
			defer openfile.Close()
		LOOP1:
			for {
				_, err := openfile.Read(filetext) //Считываем всю информацию из файла
				if err == io.EOF {
					break LOOP1
				}
			}
			go textAnalysis(alltext, nfile)
			alltext <- filetext
		default:
			break LOOP
		}
		close(alltext)
		break LOOP
	}
}

func textAnalysis(alltext chan []byte, n int) {
	var elemcount, wordlatters int
	var allword string
	word := make([]byte, 30)
LOOP2:
	for {
		select {
		case text := <-alltext:
			for text[elemcount] != 0 {
				for count := 0; text[elemcount] != ' ' && text[elemcount] != 0 && text[elemcount] != '"' && text[elemcount] != '\n' && text[elemcount] != ':' && text[elemcount] != '.' && text[elemcount] != 39 && text[elemcount] != 45 && text[elemcount] != 150 && text[elemcount] != 151; elemcount++ {
					if (text[elemcount] < 65) || text[elemcount] > 122 || (text[elemcount] < 97 && text[elemcount] > 90) {
					} else {
						if text[elemcount] > 64 && text[elemcount] < 91 {
							text[elemcount] = text[elemcount] + 32
						}
						word[count] = text[elemcount]
						count++
						wordlatters++
					}
				}
				elemcount++

				for j := 0; j < wordlatters; j++ {
					allword = (allword + string(word[j]))
					word[j] = 0
				}
				go createIndex(n, allword)
				wordlatters = 0
				allword = ""
			}
			elemcount = 0
		default:
			break LOOP2
		}
		break LOOP2
	}
}

func createIndex(n int, allword string) {
	mu.Lock()
	var check int
	var indfi string
	for k := 0; k <= records; k++ {
		if allword == index[k] {
			if k < records {
				indfi = indexfile[k]
				nf := strconv.Itoa(n)
				if strings.Contains(indexfile[k], nf) != true {
					indexfile[k] = (indfi + nf + " ")
				}
			}
			check = 1
		}
	}
	if check == 0 {
		index[records] = allword
		nf := strconv.Itoa(n)
		indexfile[records] = (" " + nf + " ")
		check = 0
		records++
	}
	check = 0
	mu.Unlock()
}

func main() {
	//Нахождение папки с файлами
	papka := os.Args[1]
	// var papka string
	// papka = ("C:\\Users\\asus\\source\\repos\\hello\\files")

	fmt.Println("\nВ папке были найдены файлы:")
	var nfile int
	dir, err := os.Open(papka)
	if err != nil {
		return
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return
	}

	var fname chan string = make(chan string)
	for _, fi := range fileInfos {
		fmt.Println(fi.Name())
		nfile++
		var str string = (papka + "/" + fi.Name())
		go openfile(fname, nfile)
		fname <- str
	}
	close(fname)

	time.Sleep(time.Millisecond * 100)

	//Создание файла и запись полученного индекса
	newfile, err := os.Create("invertindex.txt")
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
	fmt.Println("\nФайл с инвертированным индексом записан")
	indexSearch()
}
