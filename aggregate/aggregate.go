package aggregate

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gyuho/goraph/graph"
	"github.com/lycis/aggregated/configuration"
	"github.com/lycis/aggregated/extraction"
	"github.com/lycis/aggregated/operation"
	"strings"
)

var loadedAggregates map[string]*Aggregate

// An Aggregate contains all the information necessary to
// evaluate its value. This element represents one aggregate
// entry in the config file.
//
// Example:
//   my-aggregate:
//     name: My Aggregate
//     type: http
//     args:
//       url: http://example.com/foobar
//
type Aggregate struct {
	Id           string
	Name         string
	Type         string
	Args         interface{}
	OperationId  string
	Extractor    extraction.Extraction
	dependencies []string
	
	// Complete definition from YAML - for fields that are
	// not "default"
	Definition   map[interface{}]interface{}
}

func (a Aggregate) Dependencies() []string {
	if a.Type == "auto" {
		log.Debug("Resolving dependencies for 'auto' extraction")
		var d []string
		for k, _ := range loadedAggregates {
			log.WithFields(log.Fields{"aggregate-id": a.Id, "loaded-id": k}).Debug("Checking previously loaded aggregate")
			if strings.HasPrefix(k, a.Id) && k != a.Id {
				d = append(d, k)
			}
		}
		return d
	}

	return a.dependencies
}

// Sets Extractor function
func (a *Aggregate) UpdateExtractor() {
	constructor := extraction.Get(a.Type)
	if constructor == nil {
		panic(&AggregateDefinitionError{"unsupported type"})
	}

	extraction := constructor(a.Id, a.Args)
	a.Extractor = extraction
	for _, d := range extraction.Dependencies() {
		a.dependencies = append(a.dependencies, d)
	}
}

// Executes the defined aggregation and returns the
// aggregated (duh!) value.
func (a Aggregate) Value() (value extraction.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(error)
			if ok {
				value = extraction.SingleValue{""}
				err = e
			} else {
				panic(r)
			}
		}
	}()

	// get order of calls
	dependencyGraph, err := a.calculateDependencyGraph()
	if err != nil {
		return extraction.SingleValue{""}, nil
	}

	order, ok := graph.TopologicalSort(dependencyGraph)
	log.WithField("ok", ok).Debug("Toplogical sort of dependency graph performed")
	if !ok {
		log.WithField("aggregate-id", a.Id).Error("Loop in dependencies detected.")
		return extraction.SingleValue{""}, AggregateEvaluationError{"loop in dependencies"}
	}

	valueCache, err := a.resolveDependencyGraph(order)
	if err != nil {
		return extraction.SingleValue{""}, err
	}

	value = a.Extractor.Extract(valueCache)
	log.WithFields(log.Fields{"aggregate-id": a.Id, "value": value.String()}).Info("Evaluated own value")

	//value = a.executeOperation(extractedValue.String())
	return value, nil
}

func (a Aggregate) calculateDependencyGraph() (graph.Graph, error) {
	dependencyGraph := graph.NewDefaultGraph()
	if err := addDependenciesToGraph(a, dependencyGraph); err != nil {
		return nil, err
	}

	return dependencyGraph, nil
}

func (a Aggregate) resolveDependencyGraph(order []string) (map[string]extraction.Value, error) {
	// resolve dependency values
	valueCache := make(map[string]extraction.Value)
	for _, id := range order {
		dependency := GetAggregate(id)
		if dependency == nil {
			return nil, AggregateEvaluationError{fmt.Sprintf("dependency '%s' not defined")}
		}

		v, err := dependency.Value()
		if err != nil {
			return nil, err
		}
		log.WithFields(log.Fields{"aggregate-id": a.Id, "dependency-id": id, "value": v}).Info("Evaluated dependency value")

		valueCache[id] = v
	}
	return valueCache, nil
}

func (a Aggregate) executeOperation(value string) string {
	if a.OperationId == "" {
		log.WithFields(log.Fields{"aggregate-id": a.Id}).Info("No operation.")
		return value
	}

	log.WithFields(log.Fields{"aggregate-id": a.Id, "operation-id": a.OperationId, "value": value}).Info("Executing operation.")
	op := operation.Get(a.OperationId)
	if op == nil {
		log.WithFields(log.Fields{
			"aggregate-id": a.Id,
			"operation-id": a.OperationId,
		}).Error("Undefined operation")
		panicEvalError(a.Id, "undefined operation %s", a.OperationId)
	}

	newVal := op.Execute(value)
	log.WithFields(log.Fields{
		"aggregate-id": a.Id,
		"operation-id": a.OperationId,
		"value-old":    value,
		"value-new":    newVal,
	}).Info("Operation executed")
	return newVal
}

func addDependenciesToGraph(a Aggregate, g graph.Graph) error {
	log.WithFields(log.Fields{"aggregate-id": a.Id, "depdencies": len(a.Dependencies())}).Debug("calculating dependencies")
	for _, dependId := range a.Dependencies() {
		log.WithFields(log.Fields{"aggregate-id": a.Id, "dependency-id": dependId}).Debug("dependency found")
		if !g.FindVertex(dependId) {
			log.WithField("vertex", dependId).Debug("vertex added")
			g.AddVertex(dependId)
		}

		// add connection between dependencies
		g.AddEdge(dependId, a.Id, 0)
		dependencyAggregate := GetAggregate(dependId)
		addDependenciesToGraph(*dependencyAggregate, g)
	}

	return nil
}

// Parse the configuration for Aggregate definitions
// returns the number of loaded aggregates
func LoadAggregates(y configuration.YamlContent) int {
	aggregates := make(map[string]*Aggregate)

	for name, def := range y {
		defer func() {
			if r := recover(); r != nil {
				log.WithError(r.(error)).WithField("id", name).Error("Aggregate definition error")
			}
		}()

		if name != "aggregated" {
			log.WithField("id", name).Debug("Processing aggregate definition")
			aggregate := BuildAggregateFromDefinition(name, def)
			aggregate.Id = name
			log.WithField("id", name).Debug("Definition processed.")

			log.WithField("id", name).Debug("Preparing aggregate")
			aggregate.UpdateExtractor()
			log.WithFields(log.Fields{"id": name, "dependencies": len(aggregate.Dependencies())}).Debug("Aggregate prepared")

			log.WithFields(log.Fields{
				"id":   name,
				"name": aggregate.Name,
				"type": aggregate.Type,
			}).Info("Discovered aggregate")
			aggregates[name] = &aggregate
		}
	}

	loadedAggregates = aggregates
	return len(loadedAggregates)
}

// Returns the aggregate with the given id if it can be found.
// If not nil will be returned
func GetAggregate(id string) *Aggregate {
	a, ok := loadedAggregates[id]
	if !ok {
		return nil
	}
	return a
}

// Returns a list of all loaded aggregates
func Aggregates() []*Aggregate {
	list := make([]*Aggregate, 0)
	for _, a := range loadedAggregates {
		list = append(list, a)
	}
	return list
}
