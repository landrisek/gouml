package gouml

import (
	"github.com/go-kit/kit/log"
	"github.com/kazukousen/gouml/internal/gouml/plantuml"
	"github.com/kazukousen/gouml/internal/gouml/mxgraph"
)

// PlantUMLParser ...
func PlantUMLParser(logger log.Logger) Parser {
	return plantuml.NewParser(logger)
}

// PlantXMLParser ...
func PlantXMLParser(logger log.Logger) Parser {
	return mxgraph.NewParser(logger)
}
