//os.Args[1] - путь к файлу с индексом
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/polisgo2020/search-olgaklimova/index"
)

func main() {
	fileWithInd := os.Args[1]
	if len(os.Args) > 1 {
		fileWithInd = os.Args[1]
	} else {
		fmt.Println("Ошибка ввода папки с индексом")
		os.Exit(1)
	}

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
	index.IndexSearch(stringFile)
}
