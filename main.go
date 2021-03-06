package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/LindsayBradford/go-dbf/godbf"
)

func exit_error(reason string, code int) {
	fmt.Fprintf(os.Stderr, "%v\n", reason)
	os.Exit(code)
}

func showInformation(dbfTable *godbf.DbfTable, tableName string) {
	fmt.Printf("-- Generated by dbf2sql on %s\n"+
		"-- table: %v\n"+
		"-- %v fields\n"+
		"-- %v records\n\n",
		time.Now().Format(time.RFC822), tableName, len(dbfTable.Fields()), dbfTable.NumberOfRecords())
}

func createSQLTable(fields []godbf.FieldDescriptor, tableName string) {
	fmt.Printf("CREATE TABLE `%s` (\n", tableName)
	for i, fd := range fields {
		var fieldType string
		switch fd.FieldType() {
		case godbf.Numeric:
			fieldType = fmt.Sprintf("DECIMAL(%v,%v)", fd.Length(), fd.DecimalCount())
		case godbf.Float:
			fieldType = fmt.Sprintf("FLOAT(%v,%v)", fd.Length(), fd.DecimalCount())
		case godbf.Character:
			fieldType = fmt.Sprintf("CHAR(%v)", fd.Length())
		case godbf.Date:
			fieldType = "DATE"
		case godbf.Logical:
			fieldType = "BOOLEAN"
		}

		if i > 0 {
			fmt.Printf(",\n")
		}
		fmt.Printf("\t`%s` %s", fd.Name(), fieldType)
	}
	fmt.Printf("\n);\n\n")
}

func fillSQLTable(dbfTable *godbf.DbfTable, tableName string) {
	fmt.Printf("INSERT INTO `%s` VALUES\n", tableName)
	for i := 0; i < dbfTable.NumberOfRecords(); i++ {
		if i > 0 {
			fmt.Printf(",\n")
		}
		fmt.Printf("\t(")
		for j := 0; j < len(dbfTable.FieldNames()); j++ {
			if j > 0 {
				fmt.Printf(",")
			}
			fmt.Printf("'%s'", strings.Replace(dbfTable.FieldValue(i, j), "'", "''", -1))
		}
		fmt.Printf(")")
	}
	fmt.Printf(";\n\n")
}

func main() {
	dbf := flag.String("f", "", "dBase DBF file (required)")
	sql := flag.String("o", "-", "Output SQL file ('-' for stdout)")
	table := flag.String("t", "", "Table name (required)")
	encoding := flag.String("e", "UTF8", "Encoding")
	flag.Parse()

	if *dbf == "" || *table == "" {
		flag.CommandLine.Usage()
		os.Exit(1)
	}

	dbfTable, err := godbf.NewFromFile(*dbf, *encoding)
	if err != nil {
		exit_error(err.Error(), 1)
	}

	if *sql != "-" {
		if f, err := os.OpenFile(*sql, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644); err == nil {
			os.Stdout.Close()
			os.Stdout = f
		} else {
			exit_error(err.Error(), 1)
		}
	}

	showInformation(dbfTable, *table)
	createSQLTable(dbfTable.Fields(), *table)
	fillSQLTable(dbfTable, *table)
}
