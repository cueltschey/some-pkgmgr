package util

import (
  "fmt"
	"compress/gzip"
  "os"
  "os/exec"
  "io"
  "net/http"
	"path/filepath"
  "archive/tar"
  "strings"

  "github.com/ulikunitz/xz"
)

func GunzipFile(filePath string, targetFileName string) error {
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


func TarUnzipFile(tarFile, targetDir string) error {
	fmt.Println("Extracting", tarFile, "to", targetDir)

	// Open the .tar file
	file, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer file.Close()

	var reader io.Reader
	// Check if the file is gzipped (.tar.gz)
	if strings.HasSuffix(tarFile, ".tar.gz") {
		// Create a gzip reader
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer gzipReader.Close()
		reader = gzipReader
	} else if strings.HasSuffix(tarFile, ".tar.xz") {
		// Create an xz reader
		xzReader, err := xz.NewReader(file)
		if err != nil {
			return err
		}
		reader = xzReader
	} else {
		return fmt.Errorf("unsupported file format: %s", tarFile)
	}

	// Create a tar reader
	tarReader := tar.NewReader(reader)

	// Iterate through the files in the archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}

		// Prepare the file path
		targetFilePath := filepath.Join(targetDir, header.Name)

		// Check the type of file (regular file or directory)
		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory if not exists
			if err := os.MkdirAll(targetFilePath, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			// Create parent directory if not exists
			if err := os.MkdirAll(filepath.Dir(targetFilePath), 0755); err != nil {
				return err
			}

			// Create and write to the file
			file, err := os.Create(targetFilePath)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}
		}
	}

	return nil
}

func DownloadFile(url, outputDir string) error {
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

func ExecuteCommand(cmdPath string, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
  cmd.Dir = cmdPath
	output, err := cmd.CombinedOutput()
	return string(output), err
}

