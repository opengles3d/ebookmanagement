package main

import (
	"flag"

	"github.com/ebookmanagement/ebookmanager"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var rootPath string
	var dbPath string
	var cmd string
	var data string
	var pattern string
	flag.StringVar(&rootPath, "path", ".", "root path")
	flag.StringVar(&dbPath, "db", "./ebookinfo.db", "database path")
	flag.StringVar(&cmd, "cmd", "r", "r-remove duplicated, c-count files, s-search files, n-change name, p-remove specified string in file name")
	flag.StringVar(&data, "d", "", "data for cmd if needed")
	flag.StringVar(&pattern, "t", "", "pattern for cmd if needed")
	flag.Parse()

	m := ebookmanager.NewEBookManager(rootPath, dbPath)
	m.Start()
	switch cmd {
	case "c":
		m.CountFiles()
	case "s":
		m.Search(data)
	case "n":
		m.ChangeName(data)
	case "p":
		m.RemoveSpecifiedString(data, pattern)
	default:
		m.RemoveDuplicated()
	}
	m.Done()
}
