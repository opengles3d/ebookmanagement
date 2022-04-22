package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type EbookItem struct {
	Name string
	Hash string
}

func listRoot(path string) {
	listAllFiles(path)
}

func listAllFiles(path string) {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("ReadDir failed,error:%v\n", err)
		return
	}
	for _, info := range fileInfos {
		if info.IsDir() {
			listAllFiles(path + "/" + info.Name())
		} else {
			fullpath := path + "/" + info.Name()
			hash, err := calculateHash(fullpath)
			if err == nil {
				path, ok := findIt(hash)
				if ok {
					if path != fullpath {
						//fmt.Printf("Find duplicated:%s with %s, hash:%s\n", fullpath, path, hash)
						var tempName string
						if strings.Contains(path, "(1)") || strings.Contains(path, "副本") {
							tempName = path
						} else {
							tempName = fullpath
						}
						err = os.Remove(tempName)
						if err != nil {
							fmt.Printf("remove:%s failed\n", tempName)
						} else {
							fmt.Printf("remove:%s successfully\n", tempName)
						}
					}
					//remove it
				} else {
					item := EbookItem{
						Name: fullpath,
						Hash: hash,
					}
					items = append(items, item)
					addToDB(item)
				}
			} else {
				fmt.Printf("calculateHash failed,error:%v\n", err)
			}
		}
	}
}

func calculateHash(path string) (hash string, err error) {
	fp, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer fp.Close()
	buff := make([]byte, 8*1024)
	m := sha256.New()
	for {
		lens, err := fp.Read(buff)
		if err == io.EOF || lens < 0 {
			break
		}
		m.Write(buff)
	}
	res := hex.EncodeToString(m.Sum(nil))
	return res, nil
}

func findIt(hash string) (string, bool) {
	for _, item := range items {
		if item.Hash == hash {
			return item.Name, true
		}
	}
	return "", false
}

func addToDB(item EbookItem) {
	stmt, err := db.Prepare("INSERT INTO ebookinfo(path, hash) values(?,?)")
	if err != nil {
		return
	}
	_, err = stmt.Exec(item.Name, item.Hash)
	if err != nil {
		return
	}
}

func readDB() {
	//查询数据
	rows, err := db.Query("SELECT * FROM ebookinfo")
	if err != nil {
		return
	}

	for rows.Next() {
		var id string
		var path string
		var hash string
		err = rows.Scan(&id, &path, &hash)
		if err != nil {
			break
		}
		item := EbookItem{
			Name: path,
			Hash: hash,
		}
		items = append(items, item)
	}

}

var items = make([]EbookItem, 0)
var db *sql.DB

func createDBIfNotExist() {
	db_dql := `
	  CREATE TABLE "ebookinfo" (
		'id' INTEGER PRIMARY KEY AUTOINCREMENT,
		'path' VARCHAR(1024) NOT NULL,
		'hash' VARCHAR(128) NOT NULL
	);
	`

	db.Exec(db_dql)
}

func main() {
	var rootPath string
	var dbPath string
	flag.StringVar(&rootPath, "path", ".", "root path")
	flag.StringVar(&dbPath, "db", "./ebookinfo.db", "database path")
	flag.Parse()
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Printf("Open db failed,error:%v\n", err)
		return
	}
	createDBIfNotExist()
	readDB()
	listRoot(rootPath)
}
