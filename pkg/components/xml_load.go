package components

import (
	"errors"

	"github.com/beevik/etree"
)

type ComponentState interface{
	Load(map[string][]*etree.Element) error
	GetXPaths() map[string]string
}

func LoadXMLData[C ComponentState](c C, xml []byte) (C, error){
	// fetch network data
	doc := etree.NewDocument()
	err := doc.ReadFromBytes(xml)
	if (err != nil){
		return c, err
	}

	data := map[string][]*etree.Element{}
	for k, path := range c.GetXPaths(){
		data[k] = doc.FindElements(path)
	}

	err = c.Load(data)
	if (err != nil){
		return c, err
	}

	return c, nil
}

func EnsureXMLNode(data map[string][]*etree.Element, xpaths map[string]string, keys []string) error {
	for _, k := range keys{
		if (len(data[k]) == 0){
			return errors.New("Missing data in path: " + xpaths[k])
		}
	}

	return nil
}

func UnpackRequiredAttrs(keys []string, e *etree.Element) (map[string]string, error){
	attrsValues := map[string]string{}

	for _, k := range keys{
		if(e.SelectAttr(k) == nil){
			return attrsValues, errors.New("Missing required attribute: "+k)
		}

		attrsValues[k] = e.SelectAttrValue(k, "")
	}

	return attrsValues, nil 
}