package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"

    "some-pkgmgr/debian"

    "gopkg.in/yaml.v3"
    "github.com/spf13/cobra"

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
 var rootCmd = &cobra.Command{
        Use:   "some-pkgpgr",
        Short: "Some random package manager",
        Long: "A simple but fast package manager for debian, rpm, and AUR packages",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("<HELP message>")
        },
    }



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
   

    var debianCmd = &cobra.Command{
        Use:   "debian",
        Short: "Install and Update debian packages",
        Long:  "Installs debian packages on the system similarly to apt",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("HERE!")
            debupdate.UpdatePackages(config.Debian.DebUri, config.Debian.TmpDir, config.Debian.DbPath)
        },
    }

    debianCmd.PersistentFlags().StringP("action", "d", "INSTALL", "specify the debian operation <install | update>")
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
