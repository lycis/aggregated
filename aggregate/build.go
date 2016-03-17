package aggregate

// Builds an aggregate from its definition out of the config file
// In the example the id would be "my aggregate" and the definition
// would be passed as map
func BuildAggregateFromDefinition(id string, i interface{}) Aggregate {
	var aggregate Aggregate

	def, ok := i.(map[interface{}]interface{})
	if !ok {
		panicAggregateError(id, "wrong definition")
	}

	name, ok := def["name"].(string)
	if !ok {
		panicAggregateError(id, "name is not a string")
	}
	aggregate.Name = name

	xtype, ok := def["type"].(string)
	if !ok {
		panicAggregateError(id, "type is not a string")
	}
	aggregate.Type = xtype

	args, ok := def["args"]
	if ok {
		aggregate.Args = args
	}

	untypedOp, ok := def["operation"]
	if ok {
		operation, ok := untypedOp.(string)
		if !ok {
			panicAggregateError(id, "operation is not a string")
		}
		aggregate.OperationId = operation
	}

	return aggregate
}
