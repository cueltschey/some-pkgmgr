// main.go

package main

import (
    "fmt"
    "io/ioutil"
    "log"

    "gopkg.in/yaml.v3"
)

type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
    Port int `yaml:"port"`
    Host string `yaml:"host"`
}

type DatabaseConfig struct {
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
}

func main() {
    // Specify the path to your YAML config file
    filename := "config.yml"

    // Read YAML file
    yamlFile, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatalf("Error reading YAML file: %v", err)
    }

    // Parse YAML file
    var config Config
    err = yaml.Unmarshal(yamlFile, &config)
    if err != nil {
        log.Fatalf("Error parsing YAML file: %v", err)
    }

    // Print out the parsed config for verification
    fmt.Printf("Server Config: %+v\n", config.Server)
    fmt.Printf("Database Config: %+v\n", config.Database)
}
