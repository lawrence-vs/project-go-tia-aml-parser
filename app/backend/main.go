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

func exportAttributesByElement(instanceHierarchy InstanceHierarchy, filePath string) error {
	f := excelize.NewFile()
	sheetName := "GroupedAttributes"

	// Create a new sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("failed to create sheet: %v", err)
	}
	f.SetActiveSheet(index)

	// Step 1: Collect all unique attribute names
	uniqueAttributeNames := map[string]struct{}{}
	var collectAttributes func(elements []InternalElement)

	collectAttributes = func(elements []InternalElement) {
		for _, elem := range elements {
			for _, attr := range elem.Attribute {
				uniqueAttributeNames[attr.Name] = struct{}{}
			}
			collectAttributes(elem.InternalElement)
		}
	}
	collectAttributes(instanceHierarchy.InternalElement)

	// Convert map keys to a slice for ordered iteration
	attributeNames := make([]string, 0, len(uniqueAttributeNames))
	for name := range uniqueAttributeNames {
		attributeNames = append(attributeNames, name)
	}

	// Step 2: Write column headers
	f.SetCellValue(sheetName, "A1", "Element Name") // First column for element names
	for col, attrName := range attributeNames {
		cell := fmt.Sprintf("%s1", string('B'+col))
		f.SetCellValue(sheetName, cell, attrName)
	}

	// Step 3: Write rows grouped by element name
	var rowIndex int = 2
	var writeRows func(elements []InternalElement)

	writeRows = func(elements []InternalElement) {
		for _, elem := range elements {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), elem.Name) // Write element name

			// Fill attribute values in corresponding columns
			attrMap := map[string]string{}
			for _, attr := range elem.Attribute {
				attrMap[attr.Name] = attr.Value
			}

			for col, attrName := range attributeNames {
				cell := fmt.Sprintf("%s%d", string('B'+col), rowIndex)
				if value, exists := attrMap[attrName]; exists {
					f.SetCellValue(sheetName, cell, value)
				}
			}
			rowIndex++
			writeRows(elem.InternalElement) // Process nested elements
		}
	}
	writeRows(instanceHierarchy.InternalElement)

	// Save the Excel file
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("failed to save Excel file: %v", err)
	}

	return nil
}

func main() {
	// Ensure the file paths are passed as arguments (XML file and output path)
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <XML_FILE_PATH> <OUTPUT_FILE_PATH>")
		return
	}

	xmlFilePath := os.Args[1]
	excelFilePath := os.Args[2]

	// Open the XML file
	xmlFile, err := os.Open(xmlFilePath)
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

	// Export RAW information to Excel
	err = exportToExcel(caexFile.InstanceHierarchy, excelFilePath)
	if err != nil {
		fmt.Printf("Error exporting RAW Info to Excel: %v\n", err)
		return
	}
	fmt.Printf("Raw Information successfully exported to %s\n", excelFilePath)

	// Export attributes grouped by element to Excel
	err = exportAttributesByElement(caexFile.InstanceHierarchy, excelFilePath)
	if err != nil {
		fmt.Printf("Error exporting grouped attributes to Excel: %v\n", err)
		return
	}
	fmt.Printf("Grouped attributes successfully exported to %s\n", excelFilePath)
}
