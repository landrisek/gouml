package mxgraph

import (
	"bytes"
	"go/token"
	"go/types"
	"encoding/xml"
)

// Models ...
type Models []model

func (ms *Models) append(obj *types.TypeName) {
	m := model{obj: obj}
	m.build()
	*ms = append(*ms, m)
}

// WriteTo ...
func (ms Models) WriteTo(buf *bytes.Buffer) {
	for _, m := range ms {
		m.writeClass(buf)
	}
}

type model struct {
	obj     *types.TypeName
	id      string
	kind    modelKind
	field   field
	methods methods
	wrap    *types.Named
}

func (m *model) build() {
	obj := m.obj
	// *types.TypeName represents ```type [typ] [underlying]```

	// get type
	typ := obj.Type()
	m.id = extractName(typ.String())
	// TODO: obj.IsAlias() is true

	// named type (means user-defined class in OOP)
	if named, _ := typ.(*types.Named); named != nil {

		// implemented methods
		for i := 0; i < named.NumMethods(); i++ {
			f := named.Method(i)
			if isCommand(f) {
				m.kind = modelKindEntity
			}
			m.methods = append(m.methods, method{f: f})
		}
	}

	// underlying
	switch un := typ.Underlying().(type) {
	// struct
	case *types.Struct:
		m.field = field{st: un}

	// interface
	case *types.Interface:
		m.kind = modelKindInterface
		for i := 0; i < un.NumMethods(); i++ {
			m.methods = append(m.methods, method{f: un.Method(i)})
		}

	// wrap
	case *types.Slice:
		if named, _ := un.Elem().(*types.Named); named != nil {
			m.wrap = named
		}
	case *types.Map:
		if named, _ := un.Elem().(*types.Named); named != nil {
			m.wrap = named
		}

	// first-class function
	case *types.Signature:
		f := types.NewFunc(token.NoPos, obj.Pkg(), obj.Name(), un)
		m.methods = append(m.methods, method{f: f})
	}

	if m.kind == "" {
		m.kind = modelKindValueObject
	}
}

func (m model) as() string {
	return m.id
}

func (m model) writeClass(buf *bytes.Buffer) {
	id := m.as()
	newline(buf, 0)
	// package
	buf.WriteString(`package "`)
	buf.WriteString(extractPkgName(id))
	buf.WriteString(`" {`)
	// class
	newline(buf, 1)
	buf.WriteString(m.kind.Printf(extractTypeName(id), id))
	if m.field.size() > 0 || len(m.methods) > 0 {
		buf.WriteString(` {`)
		// fields
		m.field.WriteTo(buf, 2)
		// methods
		m.methods.WriteTo(buf, 2)
		newline(buf, 1)
		buf.WriteString("}")
	}
	newline(buf, 0)
	buf.WriteString("}")
}

// Model is a mxgraph element.
type Model struct {
	XMLName xml.Name  `xml:"mxGraphModel"`
	Dx         string   `xml:"dx,attr"`
	Dy         string   `xml:"dy,attr"`
	Grid       string   `xml:"grid,attr,omitempty"`
	GridSize   string   `xml:"gridSize,attr,omitempty"`
	Guides     string   `xml:"guides,attr,omitempty"`
	Tooltips   string	`xml:"tooltips,attr,omitempty"`
	Connect    string   `xml:"connect,attr"`
	Arrows     string   `xml:"arrows,attr,omitempty"`
	Fold       string   `xml:"fold,attr,omitempty"`
	Page       string   `xml:"page,attr,omitempty"`
	PageScale  string   `xml:"pageScale,attr,omitempty"`
	PageWidth  string   `xml:"pageWidth,attr,omitempty"`
	PageHeight string   `xml:"pageHeight,attr,omitempty"`
	Root       Root     `xml:"root"`
	Math       string   `xml:"math,attr,omitempty"`
	Shadow     string   `xml:"shadow,,attr,omitempty"`
}

// Root is xml element.
type Root struct {
	XMLName xml.Name `xml:"root"`
	Root []Cell 
}

// Cell is a mxgraph element.
type Cell struct {
	XMLName  xml.Name `xml:"mxCell"`
	Edge     string   `xml:"edge,attr,omitempty"`
	Source   string   `xml:"source,attr,omitempty"`
	Target   string   `xml:"target,attr,omitempty"`
	ID       string   `xml:"id,attr"`
	Style    Style    `xml:"style,attr,omitempty"`
	ParentID string   `xml:"parent,attr,omitempty"`
	Value    string   `xml:"value,attr,omitempty"`
	Vertex   string   `xml:"vertex,attr,omitempty"`
	Geometry *Geometry
}

// Geometry is a mxgraph element.
type Geometry struct {
	XMLName xml.Name `xml:"mxGeometry"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	As      string   `xml:"as,attr"`
	Array   []Point  `xml:"array"`
	Relative string `xml:"relative,attr"`
}

// Point is a mxgraph element.
type Point struct {
	XMLName xml.Name `xml:"mxPoint"`
	As      string   `xml:"as,attr"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
}

// Style models the attribute value of a Cell.
// text;html=1;strokeColor=none;fillColor=none;align=center;verticalAlign=middle;whiteSpace=wrap;rounded=0;dashed=1;
type Style struct {
	Attributes map[string]string
}

// TextStyle return a new Style with text.
func TextStyle() Style {
	return Style{Attributes: map[string]string{"text": ""}}
}

// EmptyError is given if there are no attributes to render
type EmptyError struct {}

func (err *EmptyError) Error() string {
	return "No attributes added."
}

// MarshalXMLAttr implements xml.MarshalXMLAttr
func (a Style) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if 0 == len(a.Attributes) {
		return xml.Attr{}, nil
	}
	var value string
	for key, attribute := range a.Attributes {
		value += key + "=" + attribute + ";"
	}
	return xml.Attr{Name: xml.Name{Local: "style"}, Value: value}, nil
}
