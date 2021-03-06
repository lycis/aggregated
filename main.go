package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/aggregated/aggregate"
	"github.com/lycis/aggregated/configuration"
	"net/http"
	"os"
)

// service configuration
var config *configuration.ServiceDefinition

// Main routine
func main() {
	defer func() {
		if r := recover(); r != nil {
			log.WithField("error", r).Fatal("Uncaught error")
		}
	}()

	configFile, configDir := parseCliFlags()
	loadConfiguration(configFile, configDir)
	bindForAggregates()
	startHttpServer()
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
	numLoaded := aggregate.LoadAggregates(content)
	log.WithField("count", numLoaded).Info("All aggregates loaded.")
}

func bindForAggregates() {
	for _, a := range aggregate.Aggregates() {
		resource := fmt.Sprintf("/%s", a.Id)
		log.WithField("resource", resource).Debug("Registering resource for aggregate")
		http.HandleFunc(resource, HandleGetAggregateValue)
	}
}

func startHttpServer() {
	log.WithField("bind", config.Bind).Info("Serving http://")
	err := http.ListenAndServe(config.Bind, nil)
	if err != nil {
		log.WithError(err).Fatal("Failed to serve")
	} else {
		log.Info("Stop serving http://")
	}
}
