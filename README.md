# ebookmanagement

A tool written in go to manage ebooks

1. Remove duplicated ebooks with same sha256, and keep the record in ebookinfo.db with sqlite3


Usage of ./ebookmanagement:
  -db string
    	database path (default "./ebookinfo.db")
  -path string
    	root path (default ".")
  -cmd string
    	r-remove duplicated, c-count files (default "r")

