package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/aggregated/aggregate"
	"net/http"
)

type AggregateValue struct {
	Name  string
	Value string
}

// This function is invoked for every HTTP request to get the value
// of an aggregate. It evaluates and returns the value of the aggregate.
func HandleGetAggregateValue(response http.ResponseWriter, request *http.Request) {
	id := request.URL.Path[1:]
	if len(id) < 1 {
		log.WithField("host", request.Host).Error("http: No aggregate id provided")
		http.Error(response, "no aggregate id provided", http.StatusBadRequest)
		return
	}

	a := aggregate.GetAggregate(id)
	if a == nil {
		log.WithFields(log.Fields{
			"host":         request.Host,
			"aggregate-id": id,
		}).Info("http: Requested undefined aggregate")
		http.NotFound(response, request)
		return
	}

	log.WithFields(log.Fields{
		"host":         request.Host,
		"aggregate-id": a.Id,
	}).Info("http: Aggregate value requested")

	value, err := a.Value()
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	v := AggregateValue{
		Name:  a.Name,
		Value: value,
	}

	asJson, err := json.MarshalIndent(&v, "", "  ")
	if err != nil {
		log.WithFields(log.Fields{
			"host":         request.Host,
			"aggregate-id": a.Id,
			"error":        err.Error(),
		}).Error("Failed to marshal JSON response")
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	log.WithFields(log.Fields{
		"host":         request.Host,
		"aggregate-id": a.Id,
		"response":     string(asJson),
	}).Info("JSON response served")
	fmt.Fprint(response, string(asJson))
}
