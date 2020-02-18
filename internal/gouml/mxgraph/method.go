package mxgraph

import (
	"bytes"
	"go/types"
)

type methods []method

func (ms methods) WriteTo(buf *bytes.Buffer, depth int) {
	for _, m := range ms {
		m.WriteTo(buf, depth)
	}
}

type method struct {
	f *types.Func
}

func (m method) WriteTo(buf *bytes.Buffer, depth int) {
	if m.f == nil {
		return
	}

	newline(buf, depth)
	buf.WriteString(exportedIcon(m.f.Exported()))
	// Name
	buf.WriteString(m.f.Name())

	// Signature
	sig, _ := m.f.Type().(*types.Signature)

	// parameters
	param := sig.Params()
	buf.WriteString("(")
	for i := 0; i < param.Len(); i++ {
		if i > 0 {
			buf.WriteString(", ")
		}
		v := param.At(i)
		name, typ := v.Name(), extractName(v.Type().String())
		buf.WriteString(name)
		buf.WriteString(": ")
		buf.WriteString(typ)
	}
	buf.WriteString(")")

	// results
	res := sig.Results()
	if res.Len() > 0 {
		buf.WriteString(": ")
	}
	if res.Len() > 1 {
		buf.WriteString("(")
	}
	for i := 0; i < res.Len(); i++ {
		if i > 0 {
			buf.WriteString(", ")
		}
		v := res.At(i)
		name, typ := v.Name(), extractName(v.Type().String())
		if name != "" {
			buf.WriteString(name)
			buf.WriteString(": ")
		}
		buf.WriteString(typ)
	}
	if res.Len() > 1 {
		buf.WriteString(")")
	}
}