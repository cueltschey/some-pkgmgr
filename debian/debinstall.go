package debian

import (
  "database/sql"
  "log"
  "os"
  "fmt"
  "io"
  "path/filepath"
  "strings"

  "some-pkgmgr/common"

  _ "github.com/mattn/go-sqlite3"
)

func InstallPackage(uriBase string, tmpDir string, dbPath string, packageName string){

  // Step 1: Get package info from the database
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

  // Step 2: Download .deb file
  url := uriBase + packageURL

  err = util.DownloadFile(url, tmpDir)
  if err != nil {
    log.Fatalf("Failed to Download %s:", packageName)
  }


  // Step 3: Extract .deb file 
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

  // Step 4: Install binary to root
  
  recursiveInstall(tmpDir + "/bin", "/")
}


func recursiveInstall(curDir string, destDir string) error {
  dir, err := os.Open(curDir)
  if err != nil {
    return fmt.Errorf("Failed to open %s: ", curDir)
  }
  files, err := dir.Readdir(-1)
	if err != nil {
    return fmt.Errorf("Failed to list directory %s:", curDir)
	}

	for _, file := range files {

    fileInfo, err := os.Stat(filepath.Join(curDir, file.Name()))
    if err != nil {
      if os.IsNotExist(err) {
        return fmt.Errorf("The path %s does not exist.\n", file.Name())
      } else {
        return fmt.Errorf("Error Installing %s: ", file.Name())
      }
    }
    if fileInfo.IsDir() {
      recursiveInstall(filepath.Join(curDir, file.Name()), filepath.Join(destDir, file.Name()))
    } else {
      fmt.Printf("Installing %s to %s\n", file.Name(), destDir)
      srcFile, err := os.Open(filepath.Join(curDir,file.Name()))
      if err != nil {
        return fmt.Errorf("Failed to open file %s: ", file.Name())
      }
      defer srcFile.Close()

      destFile, err := os.Create(filepath.Join(destDir, file.Name()))
      if err != nil {
        return fmt.Errorf("Error Creating %s: ", file.Name())
      }
      defer destFile.Close()

      _, err = io.Copy(destFile, srcFile)
      if err != nil {
        return fmt.Errorf("Failed to copy %s to %s: ", srcFile, destFile)
      }
      
      if strings.Contains(curDir, "usr/bin") || strings.Contains(curDir,"usr/sbin") {
        err = os.Chmod(filepath.Join(destDir, file.Name()), 0777)
      }

      err = destFile.Sync()
      if err != nil {
        return fmt.Errorf("Failed to sync file %s: ", destFile)
      }
  
    }

  }
  return nil
}
  
