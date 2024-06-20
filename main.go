package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "strings"

    "some-pkgmgr/debian"

    "gopkg.in/yaml.v3"

)

// Define the configuration structure
type Config struct {
    Debian DebConfig `yaml:"debian"`
}

// Config for debian packages
type DebConfig struct {
    DebUri string `yaml:"deb-uri"`
    Keyring string `yaml:"keyring"`
    TmpDir string `yaml:"tmpdir"`
    DbPath string `yaml:"dbpath"`
}

// Define root commands

func main() {
    // read configuration
    filename := "config.yml"
    configFile, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatalf("Error Reading config.yml: %v", err)
    }

    var config Config
    err = yaml.Unmarshal(configFile, &config)
    if err != nil {
        log.Fatalf("Error Parsing config.yml: %v", err)
    }
    // Read in CLI arguments to execute commands
   

    argc := len(os.Args)
    argv := os.Args[1:]
    if argc < 2 {
      os.Exit(1)
    }

    action := ""
    packageName := ""
    for i := 0; i < argc - 1; i++ {
        if argv[i] == "-d" || strings.HasPrefix(argv[i], "--debian="){
            if argv[i] == "-d" {
                action = argv[i+1]
            } else {
                parts := strings.SplitN(argv[i], "=", 2)
                if len(parts) == 2 {
                    action = parts[1]
                }
            }
            if i + 3 < argc{
              packageName = argv[i+2]
            }
            break
        }
    }

    switch action {
    case "install":
        fmt.Println("Installing Debian packages...")
        debian.InstallPackage(config.Debian.DebUri, config.Debian.TmpDir, config.Debian.DbPath, packageName)
    case "update":
        fmt.Printf("Updating Packages in Database: %s\n", config.Debian.DbPath )
        debian.UpdatePackages(config.Debian.DebUri, config.Debian.TmpDir, config.Debian.DbPath)
    default:
        fmt.Println("Invalid or missing action. Use -d <install|update> or --action=<install|update>")
    }


}
