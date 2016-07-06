package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/aggregated/aggregate"
	"github.com/lycis/aggregated/extraction"
	"net/http"
)

type AggregateValue struct {
	Name   string
	Result interface{}
}

// This function is invoked for every HTTP request to get the value
// of an aggregate. It evaluates and returns the value of the aggregate.
func HandleGetAggregateValue(response http.ResponseWriter, request *http.Request) {
	id := request.URL.Path[1:]
	if len(id) < 1 {
		log.WithField("client", request.RemoteAddr).Error("http: No aggregate id provided")
		http.Error(response, "no aggregate id provided", http.StatusBadRequest)
		return
	}

	a := aggregate.GetAggregate(id)
	if a == nil {
		log.WithFields(log.Fields{
			"client":       request.RemoteAddr,
			"aggregate-id": id,
		}).Info("http: Requested undefined aggregate")
		http.NotFound(response, request)
		return
	}

	log.WithFields(log.Fields{
		"client":       request.RemoteAddr,
		"aggregate-id": a.Id,
	}).Info("http: Aggregate value requested")

	value, err := a.Value()
	if err != nil {
		log.WithFields(log.Fields{
			"client":       request.RemoteAddr,
			"aggregate-id": a.Id,
			"error":        err.Error(),
		}).Error("Served error response")
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	var v AggregateValue
	v.Name = a.Name

	_, ok := value.(extraction.SingleValue)
	if ok {
		v.Result = value.String()
		log.WithFields(log.Fields{
			"client":       request.RemoteAddr,
			"aggregate-id": a.Id,
			"value-count":  1}).Debug("Calculated values")
	} else {
		mval := value.(extraction.MultiValue)
		list := mval.NestedArray()
		v.Result = list
		log.WithFields(log.Fields{
			"client":       request.RemoteAddr,
			"aggregate-id": a.Id,
			"value-count":  len(list)}).Debug("Calculated values")
	}

	asJson, err := json.MarshalIndent(&v, "", "  ")
	if err != nil {
		log.WithFields(log.Fields{
			"client":       request.RemoteAddr,
			"aggregate-id": a.Id,
			"error":        err.Error(),
		}).Error("Failed to marshal JSON response")
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	log.WithFields(log.Fields{
		"client":       request.RemoteAddr,
		"aggregate-id": a.Id,
		"response":     string(asJson),
	}).Info("JSON response served")
	fmt.Fprint(response, string(asJson))
}
