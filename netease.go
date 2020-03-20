package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func PossiableNeteaseWindows(){
	m := GetAllWindows()
	for k, v := range m {
		if strings.Contains(v.title, " - ") {
			fmt.Println(k, v)
		}
	}
}

func createFile(path string) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer f.Close()
}

func writeFile(path, description string) {
	err := ioutil.WriteFile(path, []byte(description), 0644)
	if err != nil {
		fmt.Println("Update failed with: ", err)
		return
	}
	fmt.Println("Update successfully.")
}

func main() {
	file, err := ioutil.ReadFile("./UI")
	if err != nil {
		fmt.Println("UI accidentally deleted. Please re-download the file!", err)
		return
	}
	fmt.Println(string(file))
	PossiableNeteaseWindows()

	// IO
	var handler int64
	_, err = fmt.Scan(&handler)
	if err != nil {
		fmt.Println("IO failed!", err)
		return
	}

	filepath := "./title.txt"
	createFile(filepath)

	fmt.Println("Process running......")
	for {
		go writeFile(filepath, GetAllWindows()[handler].title)
		time.Sleep(10 * time.Second)
	}
}
