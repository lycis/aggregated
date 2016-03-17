package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/aggregated/aggregate"
	"github.com/lycis/aggregated/configuration"
	"os"
)

// service configuration
var config ServiceDefinition

// all defined aggregates
// id -> Aggregate
var aggregates map[string]aggregate.AggregateF

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

	if configFile != "" {
		log.WithFields(log.Fields{
			"file": configFile,
		}).Info("Loading configuration from file")
		config = configuration.FromFile(configFile)
	} else if configDir != "" {
		log.WithField("directory", configDir).Infof("Loading configuration from directory")
		config = configuration.FromDirectory(configDir)
	} else {
		log.Fatal("No configuration source provided.")
		os.Exit(1)
	}
}
