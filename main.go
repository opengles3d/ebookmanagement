package main

import (
	"flag"

	"github.com/ebookmanagement/ebookmanager"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var rootPath string
	var dbPath string
	flag.StringVar(&rootPath, "path", ".", "root path")
	flag.StringVar(&dbPath, "db", "./ebookinfo.db", "database path")
	flag.Parse()

	m := ebookmanager.NewEBookManager(rootPath, dbPath)
	m.Start()
	m.RemoveDuplicated()
	m.Done()
}
