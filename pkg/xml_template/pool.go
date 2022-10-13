package xml_template

import (
	_ "embed"
	"text/template"
)

type PoolConfig struct{
	Path string
	Name string
}

//go:embed templates/pool.xml
var PoolXMLTemplate string
var PoolXML = template.Must(template.New("pool").Parse(PoolXMLTemplate))