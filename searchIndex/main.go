package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/polisgo2020/search-olgaklimova/index"
)

var goFormTmpl = []byte(`
<html>
	<body>
	<form action="/" method="post">
		Введите путь к файлу с индексом: <input type="file" name="fileWithInd">
		Введите сочетание слов для поиска: <input type="text" name="words">
		<input type="submit" value="Go">
	</form>
	</body>
</html>
`)

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write(goFormTmpl)
		return
	}
	fileWithInd := r.FormValue("fileWithInd")
	inputstr := r.FormValue("words")

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
	output := index.IndexSearch(stringFile, inputstr)

	fmt.Fprintln(w, output)
}

func main() {
	http.HandleFunc("/", mainPage)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
