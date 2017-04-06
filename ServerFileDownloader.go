package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var clientAddress string = "http://localhost/"
var fileName string = "world.txt"
var newFileAt = "./created"
var address string = clientAddress + fileName

func main() {

	//get file size
	chunkSize := getFileSize()

	//create the goRoutines
	var i uint64
	var wg sync.WaitGroup
	wg.Add(4)

	for i = 1; i < 5; i++ {

		start := chunkSize * (i - 1)
		end := chunkSize * i
		byterange := fmt.Sprint("bytes=", start, "-", end-1)
		go routineToFetchData(byterange, i, &wg)
	}

	wg.Wait()

	//merge downloaded files
	mergeChunkDataFiles()

	fmt.Println("Done")

}

func getFileSize() uint64 {
	req, err := http.NewRequest("HEAD", address, nil)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fileLength, _ := strconv.ParseUint(resp.Header.Get("Content-Length"), 10, 64)
	fmt.Println("len is :", fileLength)
	chunkSize := fileLength / 4

	return chunkSize
}

func routineToFetchData(byterange string, routineNo uint64, wg *sync.WaitGroup) {

	req, err := http.NewRequest("GET", address, nil)
	req.Header.Add("Range", byterange)
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("routine", routineNo, " data : ", os.Stdout, string(data))
	writeFile(string(data), routineNo)
	wg.Done()
}

func writeFile(data string, filename uint64) {
	var path = newFileAt + "/" + strconv.FormatUint(filename, 10) + ".txt"

	fmt.Println(path)

	// detect if file exists
	if _, err := os.Stat(path); os.IsExist(err) {
		var err = os.Remove(path)
		printError(err)
	}

	// create file if not exists
	var newFile, err = os.Create(path)
	printError(err)
	defer newFile.Close()

	// open file using READ & WRITE permission
	var file, errToOpen = os.OpenFile(path, os.O_RDWR, 0644)
	printError(errToOpen)
	defer file.Close()

	// write some text to file
	_, err = file.WriteString(data)
	printError(err)

	// save changes
	err = file.Sync()
	printError(err)
}

func printError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

func mergeChunkDataFiles() {

	fileHandle, err := os.Create(newFileAt + "/" + "download.txt")
	if err != nil {
		log.Fatal(err)
	}

	for i := 1; i < 5; i++ {
		// Open file for reading
		file, err := os.Open(newFileAt + "/" + strconv.Itoa(i) + ".txt")
		if err != nil {
			log.Fatal(err)
		}

		//Read data
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}

		//write data
		writer := bufio.NewWriter(fileHandle)
		defer fileHandle.Close()

		writer.Write(data)
		writer.Flush()

	}

}
