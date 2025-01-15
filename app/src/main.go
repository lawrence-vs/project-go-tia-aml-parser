package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
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

// Export the parsed data to an Excel file
func exportToExcel(instanceHierarchy InstanceHierarchy, filePath string) error {
	f := excelize.NewFile()
	sheetName := "InstanceHierarchy"

	// Create a new sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("failed to create sheet: %v", err)
	}
	f.SetActiveSheet(index)

	// Add headers
	headers := []string{"Element ID", "Element Name", "Attribute Name", "Attribute Type", "Attribute Value"}
	for col, header := range headers {
		cell := fmt.Sprintf("%s1", string('A'+col))
		f.SetCellValue(sheetName, cell, header)
	}

	// Helper function to populate rows
	var rowIndex int = 2
	var writeRows func(elements []InternalElement, depth int)

	writeRows = func(elements []InternalElement, depth int) {
		for _, elem := range elements {
			for _, attr := range elem.Attribute {
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), elem.ID)
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIndex), elem.Name)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIndex), attr.Name)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIndex), attr.Type)
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIndex), attr.Value)
				rowIndex++
			}
			writeRows(elem.InternalElement, depth+1)
		}
	}

	// Write data to the Excel file
	writeRows(instanceHierarchy.InternalElement, 0)

	// Save the Excel file
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("failed to save Excel file: %v", err)
	}

	return nil
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

	// Export to Excel
	excelPath := "../test.xlsx"
	err = exportToExcel(caexFile.InstanceHierarchy, excelPath)
	if err != nil {
		fmt.Printf("Error exporting to Excel: %v\n", err)
		return
	}

	fmt.Printf("Data successfully exported to %s\n", excelPath)
}
