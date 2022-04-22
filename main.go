package main

import (
	"flag"

	"github.com/ebookmanagement/ebookmanager"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var rootPath string
	var dbPath string
	var Cmd string
	flag.StringVar(&rootPath, "path", ".", "root path")
	flag.StringVar(&dbPath, "db", "./ebookinfo.db", "database path")
	flag.StringVar(&Cmd, "cmd", "r", "r-remove duplicated, c-count files")
	flag.Parse()

	m := ebookmanager.NewEBookManager(rootPath, dbPath)
	m.Start()
	switch Cmd {
	case "c":
		m.CountFiles()
	default:
		m.RemoveDuplicated()
	}
	m.Done()
}
