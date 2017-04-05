package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"time"
)

var path = "./doc/world.txt"
var newFilePath = "./created/"

func echoString(w http.ResponseWriter, r *http.Request) {

	fmt.Println(w, "Hello", html.EscapeString(r.URL.Path))
	http.ServeFile(w, r, r.URL.Path[1:])
}

func getFileLength() int64 {

	// re-open file
	var file, err = os.Stat(path)
	checkError(err)
	fmt.Println("file size :", file.Size())

	return file.Size()

}

func readFile(fileLength int64) {
	// re-open file
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	checkError(err)
	defer file.Close()

	readFileChunk := fileLength / 4

	time.Sleep(10 * time.Millisecond)
	go routine1(readFileChunk)
	time.Sleep(10 * time.Millisecond)
	go routine2(readFileChunk)
	time.Sleep(10 * time.Millisecond)
	go routine3(readFileChunk)
	time.Sleep(10 * time.Millisecond)
	go routine4(readFileChunk)
	time.Sleep(10 * time.Millisecond)

}

func routine1(readFileChunk int64) {
	fmt.Println("in routine1::")
	f, err := os.Open(path)
	checkError(err)
	bytes := make([]byte, readFileChunk)
	_, err = f.Read(bytes)
	checkError(err)
	fmt.Println("routine1 bytes is : ", string(bytes))
	createFile(string(bytes), "file1.txt")

}

func routine2(readFileChunk int64) {
	fmt.Println("in routine2::")
	f, err := os.Open(path)
	checkError(err)
	_, err = f.Seek(readFileChunk*1, 0)
	checkError(err)
	bytes := make([]byte, readFileChunk)
	_, err = f.Read(bytes)
	checkError(err)
	fmt.Println("routine2 bytes is : ", string(bytes))
	createFile(string(bytes), "file2.txt")

}

func routine3(readFileChunk int64) {
	f, err := os.Open(path)
	checkError(err)
	_, err = f.Seek(readFileChunk*2, 0)
	checkError(err)
	bytes := make([]byte, readFileChunk)
	_, err = f.Read(bytes)
	checkError(err)
	fmt.Println("routine3 bytes is : ", string(bytes))
	createFile(string(bytes), "file3.txt")
}

func routine4(readFileChunk int64) {

	f, err := os.Open(path)
	checkError(err)
	_, err = f.Seek(readFileChunk*3, 0)
	checkError(err)
	bytes := make([]byte, readFileChunk)
	_, err = f.Read(bytes)
	checkError(err)
	fmt.Println("routine4 bytes is : ", string(bytes))
	createFile(string(bytes), "file4.txt")

}

func createFile(data string, filename string) {
	var path = newFilePath + "/" + filename

	// detect if file exists
	if _, err := os.Stat(path); os.IsExist(err) {
		var err = os.Remove(path)
		checkError(err)
	}

	// create file if not exists

	var newFile, err = os.Create(path)
	checkError(err)
	defer newFile.Close()

	// open file using READ & WRITE permission
	var file, errToOpen = os.OpenFile(path, os.O_RDWR, 0644)
	checkError(errToOpen)
	defer file.Close()

	// write some text to file
	_, err = file.WriteString(data)
	checkError(err)

	// save changes
	err = file.Sync()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

func main() {

	http.HandleFunc("/", echoString)

	fileLength := getFileLength()

	readFile(fileLength)

	fmt.Println(fileLength)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
