package components

import (
	_ "embed"
	"text/template"
)

type VolumeConfig struct{
	Filename string
	Capacity int
}

//go:embed templates/volume.xml
var VolumeXMLTemplate string
var VolumeXML = template.Must(template.New("volume").Parse(VolumeXMLTemplate))