package deblist

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
  "runtime"
  "bufio"
  "database/sql"
  "strings"
  "time"

  _ "github.com/mattn/go-sqlite3"
)

type Package struct {
  Name              string
  Source            string
  Version           string
  InstalledSize     int
  Maintainer        string
  Architecture      string
  Depends           string
  Recommends        string
  Enhances          string
  Description       string
  Homepage          string
  DescriptionMD5    string
  Section           string
  Priority          int
  Filename          string
  Size              int
  MD5sum            string
  SHA256            string
}

const createTableSQL = `
CREATE TABLE IF NOT EXISTS packages (
    id INTEGER PRIMARY KEY,
    Name TEXT,
    Source TEXT,
    Version TEXT,
    InstalledSize INTEGER,
    Maintainer TEXT,
    Architecture TEXT,
    Depends TEXT,
    Recommends TEXT,
    Enhances TEXT,
    Description TEXT,
    Homepage TEXT,
    DescriptionMD5 TEXT,
    Section TEXT,
    Priority INTEGER,
    Filename TEXT,
    Size INTEGER,
    MD5sum TEXT,
    SHA256 TEXT
);
`

const insertPackageSQL = `
    INSERT OR REPLACE INTO packages (
        Name, Source, Version, InstalledSize, Maintainer,
        Architecture, Depends, Recommends, Enhances,
        Description, Homepage, DescriptionMD5, Section,
        Priority, Filename, Size, MD5sum, SHA256
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
`
func UpdatePackages(uriBase string, outputDir string, DbPath string) {

  archToURL := map[string]string{
		"amd64":   "binary-amd64",
		"arm64":   "binary-arm64",
		"arm":     "binary-armhf",
		"386":     "binary-i386",
		"ppc64":   "binary-ppc64",
		"s390x":   "binary-s390x",
	}

  // Step one: install Packages file into TmpDir
	url := uriBase + "dists/" + "buster/" + "main/" + archToURL[runtime.GOARCH] + "/Packages.gz"

  downloadFinished := make(chan bool)
  downloadSpinnerFinished := make(chan bool)
  go spinnerLoading("Step (1/2): Download Package File", downloadFinished, downloadSpinnerFinished)
  go func(){
    err := downloadFile(url, outputDir)
    if err != nil {
      fmt.Errorf("Error downloading Packages.gz: %v\n", err)
    }
    downloadFinished <- true
  }()

  <-downloadSpinnerFinished


  gzFileName := filepath.Join(outputDir, filepath.Base(url))
	targetFileName := gzFileName[:len(gzFileName)-len(".gz")]

  err := gunzipFile(gzFileName, targetFileName)
	if err != nil {
		fmt.Errorf("failed to decompress Packages file: %v", err)
	}

  // Step two: Update the database with the new entries

  freshPackages, err := parsePackagesFile(targetFileName)
  if err != nil{
    fmt.Errorf("failed to parse package file: %v", err)
  }
  db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
    fmt.Errorf("failed to open DB: %v", err)
	}

  _, err = db.Exec(createTableSQL)
    if err != nil {
    fmt.Errorf("failed to create packages table: %v", err)
  }

  databaseFinished := make(chan bool)
  databaseSpinnerFinished := make(chan bool)

  go spinnerLoading("Step (2/2): Update Database", databaseFinished, databaseSpinnerFinished)
  go func(){
    err = addPackagesToDatabase(freshPackages, db)
    if err != nil{
      fmt.Errorf("failed to add packages to DB: %v", err)
    }
    databaseFinished <- true
  }()

  <-databaseSpinnerFinished

	defer db.Close()
  }

func downloadFile(url, outputDir string) error {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
      return fmt.Errorf("failed to create output directory: %v", err)
		}
	}
	fileName := filepath.Base(url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get file: %v", err)
	}
	defer resp.Body.Close()

	outFilePath := filepath.Join(outputDir, fileName)
	outFile, err := os.Create(outFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	outFile.Close()

	return nil
}

func gunzipFile(filePath string, targetFileName string) error {
	gzipFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open gzip file: %v", err)
	}
	defer gzipFile.Close()

	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzipReader.Close()

	outFile, err := os.Create(targetFileName)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, gzipReader)
	if err != nil {
		return fmt.Errorf("failed to write to output file: %v", err)
	}
  return nil
}


func parsePackagesFile(filename string) ([]Package, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var packages []Package
    var currentPackage Package

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()

        if line == "" { // Blank line indicates end of package entry
            if currentPackage.Name != "" {
                packages = append(packages, currentPackage)
                currentPackage = Package{}
            }
        } else {
            parts := strings.SplitN(line, ": ", 2)
            if len(parts) == 2 {
                switch parts[0] {
                case "Package":
                    currentPackage.Name = parts[1]
                case "Source":
                    currentPackage.Source = parts[1]
                case "Version":
                    currentPackage.Version = parts[1]
                case "Installed-Size":
                    fmt.Sscanf(parts[1], "%d", &currentPackage.InstalledSize)
                case "Maintainer":
                    currentPackage.Maintainer = parts[1]
                case "Architecture":
                    currentPackage.Architecture = parts[1]
                case "Depends":
                    currentPackage.Depends = parts[1]
                case "Recommends":
                    currentPackage.Recommends = parts[1]
                case "Enhances":
                    currentPackage.Enhances = parts[1]
                case "Description":
                    currentPackage.Description = parts[1]
                case "Homepage":
                    currentPackage.Homepage = parts[1]
                case "Description-md5":
                    currentPackage.DescriptionMD5 = parts[1]
                case "Section":
                    currentPackage.Section = parts[1]
                case "Priority":
                    fmt.Sscanf(parts[1], "%d", &currentPackage.Priority)
                case "Filename":
                    currentPackage.Filename = parts[1]
                case "Size":
                    fmt.Sscanf(parts[1], "%d", &currentPackage.Size)
                case "MD5sum":
                    currentPackage.MD5sum = parts[1]
                case "SHA256":
                    currentPackage.SHA256 = parts[1]
                }
            }
        }
    }
    if currentPackage.Name != "" {
        packages = append(packages, currentPackage)
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    return packages, nil
}

func addPackagesToDatabase(freshPackages []Package, db *sql.DB) (error){
  tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        } else {
            tx.Commit()
        }
    }()

    stmt, err := tx.Prepare(insertPackageSQL)
    if err != nil {
        return err
    }
  defer stmt.Close()

  for _, pkg := range freshPackages {
    _, err = stmt.Exec(
      pkg.Name, pkg.Source, pkg.Version, pkg.InstalledSize, pkg.Maintainer,
      pkg.Architecture, pkg.Depends, pkg.Recommends, pkg.Enhances,
      pkg.Description, pkg.Homepage, pkg.DescriptionMD5, pkg.Section,
      pkg.Priority, pkg.Filename, pkg.Size, pkg.MD5sum, pkg.SHA256,
    )
      if err != nil {
        return err
      }
    }

  return nil
}


func spinnerLoading(message string, done chan bool, spinnerDone chan bool) {
	// Define a set of spinner characters
	spinChars := []rune{'|', '/', '-', '\\'}
	i := 0
	for {
		select {
		case <-done:
      fmt.Printf(" \u2713 \n")
      spinnerDone <- true
			return
		default:
			fmt.Printf("\r%c %s", spinChars[i], message)
			i = (i + 1) % len(spinChars)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
