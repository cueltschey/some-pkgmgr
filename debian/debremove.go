package debian

import (
  "fmt"
  "os"

)

const getInstalledSQL = `
SELECT * FORM installed WHERE PackageName = ?;
`


func RemovePackage(DbPath string, PackageName string){
  db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
      log.Fatalf("Failed to open database: %v", err)
  }

  rows, err = db.Exec(getInstalledSQL, "cowsay") 
  for row, err := range rows {
    fmt.Printf(row)
  }

}
