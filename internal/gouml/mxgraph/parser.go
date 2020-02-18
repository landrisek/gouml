package mxgraph

import (
	"bytes"
	"go/types"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// NewParser ...
func NewParser(logger log.Logger) *parser {
	return &parser{
		logger: log.With(logger, "component", "parser"),
		models: Models{},
	}
}

type parser struct {
	logger log.Logger
	models Models
}

func (p *parser) Build(pkgs []*types.Package) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		level.Debug(p.logger).Log("msg", "built uml", "ms", elapsed.Truncate(time.Millisecond))
	}()

	objects := []types.Object{}
	for _, pkg := range pkgs {
		scope := pkg.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			objects = append(objects, obj)
		}
	}

	for _, obj := range objects {
		switch obj := obj.(type) {

		// declared type
		case *types.TypeName:
			p.models.append(obj)

		}
	}
}

func (p parser) WriteTo(buf *bytes.Buffer) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		level.Debug(p.logger).Log("msg", "write to file", "ms", elapsed.Truncate(time.Millisecond))
	}()

	p.models.WriteTo(buf)
	newline(buf, 0)
	newline(buf, 0)
}
