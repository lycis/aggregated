package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/aggregated/aggregate"
	"github.com/lycis/aggregated/configuration"
	"os"
)

// service configuration
var config *configuration.ServiceDefinition

// all defined aggregates
// id -> Aggregate
var aggregates map[string]aggregate.Aggregate

// Main routine
func main() {
	defer func() {
		if r := recover(); r != nil {
			log.WithField("error", r).Fatal("Uncaught error")
		}
	}()

	configFile, configDir := parseCliFlags()
	loadConfiguration(configFile, configDir)

	// bind aggregates here
}

// Parses the cli options and returns its parameters
func parseCliFlags() (cFile string, cDir string) {
	configFile := flag.String("config-file", "", "configuration file")
	configDir := flag.String("config-dir", "/etc/aggregated/conf.d", "configuration directory")
	flag.Parse()

	return *configFile, *configDir
}

// Loads the service configuration from the given option.
// That is either file or directory
func loadConfiguration(configFile string, configDir string) {
	if configFile != "" && configDir != "" {
		log.Fatal("Configuration file and directory provided. Please give only one of them.")
		os.Exit(1)
	}

	var content configuration.YamlContent
	if configFile != "" {
		log.WithFields(log.Fields{
			"file": configFile,
		}).Info("Loading configuration from file")
		content = configuration.FromFile(configFile)
	} else if configDir != "" {
		log.WithField("directory", configDir).Infof("Loading configuration from directory")
		content = configuration.FromDirectory(configDir)
	} else {
		log.Fatal("No configuration source provided.")
		os.Exit(1)
	}

	config = content.ServiceDefintion()
	aggregate.LoadAggregates(content)
}
