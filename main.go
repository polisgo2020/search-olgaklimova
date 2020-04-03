//os.Args[1] - путь к папке
//os.Args[2] - путь к файлу для записи индекса
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

var index = make(map[int]string)
var indexfile = make(map[int]string)
var numberfile = make(map[int]int)
var mu = &sync.Mutex{}
var records int
var newrecords int

func openfile(chanFilename chan string, nfile int) {
	alltext := make(chan string)

LOOP:
	for {
		select {
		case filename := <-chanFilename:
			openfile, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}
			defer openfile.Close()

			byteFile, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Fatal(err)
			}

			stringFile := string(byteFile)

			go textAnalysis(alltext, nfile)
			alltext <- stringFile
		default:
			break LOOP
		}
		close(alltext)
		break LOOP
	}
}

func textAnalysis(alltext chan string, nfile int) {
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
				createIndex(nfile, allwords[k])
			}
		default:
			break LOOP2
		}
		break LOOP2
	}
}

func createIndex(nfile int, allword string) {
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
}

func main() {
	//Нахождение папки с файлами
	papka := os.Args[1]
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

	chanFilename := make(chan string)
	for _, fi := range fileInfos {
		fmt.Println(fi.Name())
		nfile++
		var filename string = (papka + "/" + fi.Name())
		go openfile(chanFilename, nfile)
		chanFilename <- filename
	}
	close(chanFilename)

	time.Sleep(time.Millisecond * 100500)

	//Создание файла и запись полученного индекса
	newfile, err := os.Create(os.Args[2])
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
