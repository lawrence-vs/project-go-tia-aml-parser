# Siemens TIA AML Parser

## Context 
An AML file is an XML-based data format for storing and exchanging project planning data. 
The TIA Selection Tool (TST) can be used to configure siemens automation harware from an order list for required components 

## Scope
This project aim to create an automation the IO list creation using an exported AutomationML file.

## Getting Started
1. install go dependencies
    ```bash
    cd app/src
    go mod init tia-aml-parser
    go install github.com/xuri/excelize/v2@latest
    go mod tidy
    go build -o tia-aml-parser
    cd ../../
    ```
1. copy your aml file to the app directory.
2. then run the sciprt to export the aml file data to an excel IO list.
- manually
```bash
cd app/src
go run main.go
```

- using powershell

```
./run.ps1
```