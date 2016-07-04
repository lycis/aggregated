package main

// This file contains all the plugins that will be loaded.
//
// That may be operations, aggregate types or extractions. To
// load a plugin just add it to the import below.

import (
	_ "github.com/lycis/aggregated/extraction/aggregate"
	_ "github.com/lycis/aggregated/extraction/auto"
	_ "github.com/lycis/aggregated/extraction/http"
	_ "github.com/lycis/aggregated/extraction/static"
	_ "github.com/lycis/aggregated/operation/math"
	_ "github.com/lycis/aggregated/operation/strings"
)
