package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"unicode/utf8"
)
//import "io/ioutil"

import "bufio"

type TOKEN int

const (
	INSERT = iota
	FROM
	OPEN_BRACKET
	CLOSE_BRACKET
	COMMA
	COLUMN_NAME
	VALUES
	COLUMN

)

func main() {
	//args := os.Args
	//pwd, _:= os.Getwd()
	//fmt.Println(args, pwd)
	path := "/home/stefanpopescu/p/go/insert-fmt/test.sql"
	file, _ :=os.Open(path)
	defer file.Close()
	sqlScanner := bufio.NewScanner(file)
	sqlScanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		start := 0
		for width := 0; start < len(data); start += width {
			var r rune
			r, width = utf8.DecodeRune(data[start:])

		}
	})

	file, err := ioutil.ReadFile("/home/stefanpopescu/p/go/insert-fmt/test.sql")
	if err == nil {
		fmt.Println("Hello ARAM")
		fmt.Println(file)
		return
	}
	fmt.Println(err)
}
