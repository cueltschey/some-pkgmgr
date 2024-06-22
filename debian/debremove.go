package debian

import (
  "database/sql"
  "fmt"
  "os"
  "log"

	_ "github.com/mattn/go-sqlite3"
)

const getInstalledSQL = `
SELECT * FROM installed WHERE PackageName = ?;
`


func RemovePackage(DbPath string, PackageName string){
  db, err := sql.Open("sqlite3", DbPath)
    if err != nil {
      log.Fatalf("Failed to open database: %v", err)
  }

  rows, err := db.Query(getInstalledSQL, PackageName)
  if err != nil {
    log.Fatalf("Failed to Query database: %v" ,err)
  }

  for rows.Next() {
		var fileName string
		var pkg string
		err = rows.Scan(&pkg, &fileName)
		if err != nil {
			log.Fatal(err)
		}
    fileInfo, err := os.Stat(fileName)
    if !os.IsNotExist(err) && !fileInfo.IsDir() {
      fmt.Printf("Package: %s, File: %s\n", pkg, fileName)
      os.Remove(fileName)
    } 
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
