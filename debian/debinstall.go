package debian

import (
  "database/sql"
  "log"
  //"fmt"
  "path/filepath"

  "some-pkgmgr/common"

  _ "github.com/mattn/go-sqlite3"
)

func InstallPackage(uriBase string, tmpDir string, dbPath string, packageName string){
  db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

  var packageURL string
  query := "SELECT Filename FROM packages WHERE Name = ?"
	err = db.QueryRow(query, packageName).Scan(&packageURL)
  if err != nil {
		if err == sql.ErrNoRows {
      log.Printf("No packages found for: %s\n", packageName)
		} else {
			log.Fatalf("Failed to query database: %v", err)
		}
	}
  url := uriBase + packageURL

  err = util.DownloadFile(url, tmpDir)
  if err != nil {
    log.Fatalf("Failed to Download %s:", packageName)
  }

  arFileName := filepath.Join(tmpDir, filepath.Base(url))

  _, err = util.ExecuteCommand(tmpDir, "ar", "x", arFileName)
  if err != nil {
    log.Fatalf("Failed to Extract %s:", arFileName)
  }
  err = util.TarUnzipFile(tmpDir + "/data.tar.xz", tmpDir + "/bin")
  if err != nil {
    log.Fatalf("Failed to Extract data.tar.xz: %v", err)
  }
  err = util.TarUnzipFile(tmpDir + "/control.tar.xz", tmpDir)
  if err != nil {
    log.Fatalf("Failed to Extract control.tar.xz:")
  }
}
