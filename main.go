package main

import (
    "fmt"
    "io/ioutil"
    "log"

    "some-pkgmgr/cli"
    "some-pkgmgr/debian"

    "gopkg.in/yaml.v3"
)

type Config struct {
    Debian DebConfig `yaml:"debian"`
}

type DebConfig struct {
    DebUri string `yaml:"deb-uri"`
    Keyring string `yaml:"keyring"`
    TmpDir string `yaml:"tmpdir"`
    DbPath string `yaml:"dbpath"`
}

func main() {
    filename := "config.yml"
    configFile, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatalf("Error Reading config.yml: %v", err)
    }

    var config Config
    err = yaml.Unmarshal(configFile, &config)
    if err != nil {
        log.Fatalf("Error parsing YAML file: %v", err)
    }

    // Print out the parsed config for verification
    fmt.Printf("Debian Config: %+v\n", config.Debian)
    cli.Run()
    deblist.UpdatePackages(config.Debian.DebUri, config.Debian.TmpDir, config.Debian.DbPath)
}
