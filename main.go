//os.Args[1] - путь к папке
//os.Args[2] - путь к файлу для записи индекса
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/polisgo2020/search-olgaklimova/index"
)

func openfile(filename string, nfile int, wg *sync.WaitGroup) {
	defer wg.Done()
	alltext := make(chan string)

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

	wg.Add(1)
	go index.TextAnalysis(alltext, nfile, wg)
	alltext <- stringFile
}

func main() {
	papka := os.Args[1]
	wg := &sync.WaitGroup{}

	if len(os.Args) > 2 {
		papka = os.Args[1]
	} else {
		fmt.Println("Ошибка ввода папки с файлами или файла для записи")
		os.Exit(1)
	}
	fmt.Println("\nВ папке были найдены файлы:")
	dir, err := os.Open(papka)
	if err != nil {
		return
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return
	}

	for i, fi := range fileInfos {
		fmt.Println(fi.Name())
		i++
		var filename string = (papka + "/" + fi.Name())
		wg.Add(1)
		go openfile(filename, i, wg)
	}

	wg.Wait()

	fmt.Println(index.WriteIndex())
}
