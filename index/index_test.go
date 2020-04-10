package index

import (
	"sync"
	"testing"
)

func TestCreateIndex(t *testing.T) {
	in1 := 1
	in2 := ("tucking")
	expected := ("tucking")
	actual := СreateIndex(in1, in2)
	if !testEq(expected, actual) {
		t.Errorf("%v is not eqal to expected %v", actual, expected)
	}
}
func TestWriteIndex(t *testing.T) {
	expected := ("\nФайл с инвертированным индексом записан")
	actual := WriteIndex()
	if !testEq(expected, actual) {
		t.Errorf("%v is not eqal to expected %v", actual, expected)
	}
}

func TestCoincidencesIndexes(t *testing.T) {
	in1 := []byte{49, 50, 51, 52, 49, 50, 51, 52, 49, 50, 51, 52, 49, 50, 51}
	in2 := 15
	in3 := 4
	expected := (" 1 2 3")
	actual := CoincidencesIndexes(in1, in2, in3)
	if !testEq(expected, actual) {
		t.Errorf("%v is not eqal to expected %v", actual, expected)
	}
}

func TestIndexSearch(t *testing.T) {
	in1 := ("suddenly { 1 4 2 3 } the { 1 4 3 2 } red { 1 2 3 4 } brown { 1 2 3 4 }")
	in2 := ("suddenly the red")
	expected := ("Введенные слова были найдены в текстах: 1 2 3 4\n")
	actual := IndexSearch(in1, in2)
	if !testEq(expected, actual) {
		t.Errorf("%v is not eqal to expected %v", actual, expected)
	}
}

func testEq(a, b string) bool {

	if len(a) != len(b) {
		return false
	}
	return true
}

func TestTextAnalysis(t *testing.T) {
	in1 := make(chan string)
	in2 := 1
	in3 := &sync.WaitGroup{}
	out := &sync.WaitGroup{}
	in3.Add(1)
	actual := TextAnalysis(in1, in2, in3)
	expected := out
	testAn(expected, actual)
}

func testAn(a, b *sync.WaitGroup) bool {
	a.Wait()
	b.Wait()
	return true
}

/*
	go test -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html
*/
