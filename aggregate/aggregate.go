package aggregate

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gyuho/goraph/graph"
	"github.com/lycis/aggregated/configuration"
	"github.com/lycis/aggregated/extraction"
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
}

func (a Aggregate) Dependencies() []string {
	return a.dependencies
}

// Sets Extractor function
func (aggregate *Aggregate) UpdateExtractor() {
	switch aggregate.Type {
	case "http":
		aggregate.applyHttpExtractor()
	case "aggregate":
		aggregate.applyAggregateExtractor()
	case "auto":
		aggregate.applyAutoExtractor()
	default:
		panic(&AggregateDefinitionError{"unsupported type"})
	}
}

// applies an HTTP Extractor to the given aggregate
func (aggregate *Aggregate) applyHttpExtractor() {
	parameters, ok := aggregate.Args.(map[interface{}]interface{})
	if !ok {
		panicParameterError(aggregate.Id, "http")
	}

	extractor := extraction.HttpExtraction{
		Url: parameters["url"].(string),
	}

	aggregate.Extractor = extractor
}

func (a *Aggregate) applyAggregateExtractor() {
	parameters, ok := a.Args.([]interface{})
	if !ok {
		panicParameterError(a.Id, "aggregate")
	}

	var extractor extraction.AggregateExtraction
	for _, sub := range parameters {
		typedSub, ok := sub.(string)
		if !ok {
			panicParameterError(a.Id, "aggregate")
		}
		extractor.Ids = append(extractor.Ids, typedSub)
		a.dependencies = append(a.dependencies, typedSub)
	}

	a.Extractor = extractor
	log.WithField("count", len(a.dependencies)).Debug("Dependencies added.")
}

func (a *Aggregate) applyAutoExtractor() {
	extractor := extraction.AutoExtraction{a.Id}
	a.Extractor = extractor
}

// Executes the defined aggregation and returns the
// aggregated (duh!) value.
func (a Aggregate) Value() (string, error) {

	// get order of calls
	dependencyGraph, err := a.calculateDependencyGraph()
	if err != nil {
		return "", nil
	}

	order, ok := graph.TopologicalSort(dependencyGraph)
	log.WithField("ok", ok).Debug("Toplogical sort of dependency graph performed")
	if !ok {
		log.WithField("aggregate-id", a.Id).Error("Loop in dependencies detected.")
		return "", AggregateEvaluationError{"loop in dependencies"}
	}

	// resolve dependency values
	valueCache := make(map[string]string)
	for _, id := range order {
		dependency := GetAggregate(id)
		if dependency == nil {
			return "", AggregateEvaluationError{fmt.Sprintf("dependency '%s' not defined")}
		}

		v, err := dependency.Value()
		if err != nil {
			return "", err
		}
		log.WithFields(log.Fields{"aggregate-id": a.Id, "dependency-id": id, "value": v}).Info("Evaluated dependency value")

		valueCache[id] = v
	}

	value := a.Extractor.Extract()
	log.WithFields(log.Fields{"aggregate-id": a.Id, "value": value}).Info("Evaluated own value")
	
	// TODO apply operation
	return value, nil
}

func (a Aggregate) calculateDependencyGraph() (graph.Graph, error) {
	dependencyGraph := graph.NewDefaultGraph()
	if err := addDependenciesToGraph(a, dependencyGraph); err != nil {
		return nil, err
	}

	return dependencyGraph, nil
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
