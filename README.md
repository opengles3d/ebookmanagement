# ebookmanagement

A tool written in go to manage ebooks

1. Remove duplicated ebooks with same sha256, and keep the record in ebookinfo.db with sqlite3
2. count files in directory
3. change name by the epub information
4. remove no use content in file name

Usage of ./ebookmanagement:
  -db string
    	database path (default "./ebookinfo.db")
  -path string
    	root path (default ".")
  -cmd string
    	r-remove duplicated, c-count files, s-search files, n-change name, p-remove specified string in file name (default "r")

  -d string
    	parameters of cmd
        s contents to search for
        n file or directory(all included files) to change name
        t pattern for cmd is p


Ex:./ebookmanagement -cmd p -d /home/shaocq/temp/ttt/ -t '【公众号：书单严选】' 
