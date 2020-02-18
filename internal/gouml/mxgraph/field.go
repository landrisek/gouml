package mxgraph

import (
	"bytes"
	"go/types"
)

type field struct {
	st *types.Struct
}

func (f field) size() int {
	if f.st == nil {
		return 0
	}
	return f.st.NumFields()
}

func (f field) WriteTo(buf *bytes.Buffer, depth int) {
	if f.st == nil {
		return
	}
	for i := 0; i < f.st.NumFields(); i++ {
		newline(buf, depth)
		v := f.st.Field(i)
		buf.WriteString(exportedIcon(v.Exported()))
		buf.WriteString(v.Name())
		buf.WriteString(": ")
		f.typeString(buf, v.Type())
	}
}

func (f field) typeString(buf *bytes.Buffer, typ types.Type) {
	switch typ := typ.(type) {
	case *types.Struct:
		buf.WriteString("struct{")
		for i := 0; i < typ.NumFields(); i++ {
			if i > 0 {
				buf.WriteString("; ")
			}
			v := typ.Field(i)
			buf.WriteString(v.Name())
			buf.WriteString(": ")
			f.typeString(buf, v.Type())
		}
		buf.WriteString("}")
		return
	}
	buf.WriteString(extractName(typ.String()))
}
