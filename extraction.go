package main

import (
	log "github.com/Sirupsen/logrus"
)

type extraction interface {
	Extract() string
}

type HttpExtraction struct {
	Url string
}

func (e HttpExtraction) Extract() string {
	log.Debug("http extraction")
	return ""
}