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
	flag.StringVar(&rootPath, "path", ".", "root path")
	flag.StringVar(&dbPath, "db", "./ebookinfo.db", "database path")
	flag.StringVar(&cmd, "cmd", "r", "r-remove duplicated, c-count files, s-search files")
	flag.StringVar(&data, "d", "", "data for cmd if needed")
	flag.Parse()

	m := ebookmanager.NewEBookManager(rootPath, dbPath)
	m.Start()
	switch cmd {
	case "c":
		m.CountFiles()
	case "s":
		m.Search(data)
	default:
		m.RemoveDuplicated()
	}
	m.Done()
}
