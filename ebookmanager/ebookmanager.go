package ebookmanager

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type EBookItem struct {
	Name string
	Hash string
}

type EBookManager struct {
	db    *sql.DB
	items []EBookItem

	rootPath string
	dbPath   string
}

func NewEBookManager(rootPath string, dbPath string) *EBookManager {
	m := EBookManager{

		items:    make([]EBookItem, 0),
		rootPath: rootPath,
		dbPath:   dbPath,
	}

	return &m
}

func (m *EBookManager) Start() {
	m.openDB()
	m.createDBIfNotExist()
	m.readDB()
}

func (m *EBookManager) RemoveDuplicated() {
	m.processRoot()
}

func (m *EBookManager) CountFiles() {
	count := m.countFilesInDirs(m.rootPath)
	fmt.Printf("There are %d files in %s\n", count, m.rootPath)
}

func (m *EBookManager) Done() {
	m.db.Close()
}

func (m *EBookManager) processRoot() {
	m.processDirs(m.rootPath)
}

func (m *EBookManager) processDirs(path string) {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("ReadDir failed,error:%v\n", err)
		return
	}
	for _, info := range fileInfos {
		if info.IsDir() {
			m.processDirs(path + "/" + info.Name())
		} else {
			fullpath := path + "/" + info.Name()
			m.processFile(fullpath)
		}
	}
}

func (m *EBookManager) processFile(fullpath string) {
	hash, err := calculateHash(fullpath)
	if err != nil {
		fmt.Printf("calculateHash failed,error:%v\n", err)
		return
	}

	path, duplicated := m.findIt(hash)
	if duplicated {
		if path != fullpath {
			removeFile(path, fullpath)
		}
	} else {
		item := EBookItem{
			Name: fullpath,
			Hash: hash,
		}
		m.addToItems(item)
		m.addToDB(item)
	}

}

func (m *EBookManager) addToItems(item EBookItem) {
	m.items = append(m.items, item)
}

func (m *EBookManager) createDBIfNotExist() {
	db_dql := `
	  CREATE TABLE "ebookinfo" (
		'id' INTEGER PRIMARY KEY AUTOINCREMENT,
		'path' VARCHAR(1024) NOT NULL,
		'hash' VARCHAR(128) NOT NULL
	);
	`

	m.db.Exec(db_dql)
}

func (m *EBookManager) openDB() {
	var err error
	m.db, err = sql.Open("sqlite3", m.dbPath)
	if err != nil {
		fmt.Printf("Open db failed,error:%v\n", err)
		return
	}
}

func (m *EBookManager) readDB() {
	rows, err := m.db.Query("SELECT * FROM ebookinfo")
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
		item := EBookItem{
			Name: path,
			Hash: hash,
		}
		m.addToItems(item)
	}
}

func (m *EBookManager) addToDB(item EBookItem) {
	stmt, err := m.db.Prepare("INSERT INTO ebookinfo(path, hash) values(?,?)")
	if err != nil {
		return
	}
	_, err = stmt.Exec(item.Name, item.Hash)
	if err != nil {
		return
	}
}

func (m *EBookManager) findIt(hash string) (string, bool) {
	for _, item := range m.items {
		if item.Hash == hash {
			return item.Name, true
		}
	}
	return "", false
}

func (m *EBookManager) countFilesInDirs(path string) int {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("ReadDir failed,error:%v\n", err)
		return 0
	}
	count := 0
	for _, info := range fileInfos {
		if info.IsDir() {
			count += m.countFilesInDirs(path + "/" + info.Name())
		} else {
			count++
		}
	}
	return count
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

func removeFile(path string, fullpath string) {
	//fmt.Printf("Find duplicated:%s with %s, hash:%s\n", fullpath, path, hash)
	var tempName string
	if strings.Contains(path, "(1)") || strings.Contains(path, "副本") {
		tempName = path
	} else {
		tempName = fullpath
	}
	err := os.Remove(tempName)
	if err != nil {
		fmt.Printf("remove:%s failed\n", tempName)
	} else {
		fmt.Printf("remove:%s successfully\n", tempName)
	}
}
