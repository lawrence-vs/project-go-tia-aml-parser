package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

// Structs to represent the XML hierarchy
type CAEXFile struct {
	XMLName           xml.Name          `xml:"CAEXFile"`
	InstanceHierarchy InstanceHierarchy `xml:"InstanceHierarchy"`
}

type InstanceHierarchy struct {
	Name            string            `xml:"Name,attr"`
	InternalElement []InternalElement `xml:"InternalElement"`
}

type InternalElement struct {
	ID              string            `xml:"ID,attr"`
	Name            string            `xml:"Name,attr"`
	Attribute       []Attribute       `xml:"Attribute"`
	InternalElement []InternalElement `xml:"InternalElement"`
}

type Attribute struct {
	Name  string `xml:"Name,attr"`
	Type  string `xml:"AttributeDataType,attr,omitempty"`
	Value string `xml:"Value"`
}

// Helper function to recursively process InternalElements
func processInternalElements(elements []InternalElement, depth int) {
	for _, elem := range elements {
		indent := ""
		for i := 0; i < depth; i++ {
			indent += "  "
		}
		fmt.Printf("%sElement ID: %s, Name: %s\n", indent, elem.ID, elem.Name)
		for _, attr := range elem.Attribute {
			fmt.Printf("%s  Attribute - Name: %s, Type: %s, Value: %s\n", indent, attr.Name, attr.Type, attr.Value)
		}
		processInternalElements(elem.InternalElement, depth+1)
	}
}

func main() {
	// Open the XML file
	xmlFile, err := os.Open("../example.xml")
	if err != nil {
		fmt.Println("Error opening XML file:", err)
		return
	}
	defer xmlFile.Close()

	// Parse the XML file
	var caexFile CAEXFile
	decoder := xml.NewDecoder(xmlFile)
	err = decoder.Decode(&caexFile)
	if err != nil {
		fmt.Println("Error parsing XML file:", err)
		return
	}

	// Display the parsed data
	fmt.Println("Automation IO List:")
	fmt.Println("====================")
	processInternalElements(caexFile.InstanceHierarchy.InternalElement, 0)
}
