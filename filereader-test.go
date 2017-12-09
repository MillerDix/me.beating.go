package main

import (
	"fmt"
	"os"
	"bufio"
	"reflect"

	// "github.com/russross/blackfriday"
)

func main() {
	// data, err := ioutil.ReadFile("./流程.md")
	file, err := os.Open("./流程.md")
	if err != nil {
		fmt.Println("os open failed")
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(reflect.TypeOf(scanner.Text()))
		fmt.Println(scanner.Text())
	}
}