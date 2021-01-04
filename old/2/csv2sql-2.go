package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
)

var csvFileName string
var tableName string

func main() {
	//-------------------------------------------------------------------------
	// sort out the command line arguments
	//-------------------------------------------------------------------------
	// get the command line args passed to the program

	csvFileName = "Book1.csv"
	tableName = "ExampleTable2"

	// check we have a table name and csv file to work with - otherwise abort
	if csvFileName == "" || tableName == "" {
		fmt.Println("ERROR: please provide both a table name and a CSV File Name")
		//fmt.Println("Usage:",flag.Usage,"Command Line:",flag.CommandLine)
		os.Exit(-2)
	}

	//-------------------------------------------------------------------------
	// open and prepare the CSV input file
	//-------------------------------------------------------------------------

	file, err := os.Open(csvFileName)
	if err != nil {
		fmt.Printf("Error Opening your CSV File: %s\n", err)
		os.Exit(-3)
	}
	defer file.Close()

	// connect a CSV reader to the file handle - which is the actual opened
	// CSV file
	// TODO : is there an error from this to check?
	reader := csv.NewReader(file)

	//-------------------------------------------------------------------------
	// open and prepare the SQL output file
	//-------------------------------------------------------------------------
	// get a new filename to write the SQl converted data into - call our
	// function SQLFileName() to obtain a suitable string for the new filename
	// TODO : add option to output to stdout instead of a file only
	sqlOutFile := SQLFileName()
	if debugSwitch {
		fmt.Println("Opening the SQL output file:", sqlOutFile)
	}

	// open the new file using the name we obtained above - handle 'filesql'
	filesql, err := os.Create(sqlOutFile)
	// error - if we have one when trying open & create the new file
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// now new file is open - defer the close of the file handle until we return
	defer filesql.Close()
	// attach the opened new sql file handle to a buffered file writer
	// the buffered file writer has the handle 'sqlFileBuffer'
	sqlFileBuffer := bufio.NewWriter(filesql)

	//--------------------------------------------------------------------------------
	// prepare to read the each line of the CSV file - and write out to the SQL
	//--------------------------------------------------------------------------------

	// track the number of lines in the csv file	// track number of fields in csv file
	lineCount := 0
	csvFields := 0

	// grab time now - so can calculate how long it takes to process the file
	start := time.Now()

	// create a buffer to hold each line of the SQL file as we build it
	// handle to this buffer is called 'strbuffer'
	var strbuffer bytes.Buffer

	// START - processing of each line in the CSV input file
	//-------------------------------------------------------------------------
	// loop through the csv file until EOF - or until we hit an error in parsing it.
	// Data is read in for each line of the csv file and held in the variable
	// 'record'.  Build a string for each line - wrapped with the SQL and
	// then output to the SQL file writer in its completed new form
	//-------------------------------------------------------------------------
	for {
		record, err := reader.Read()

		// if we hit end of file (EOF) or another unexpected error
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// get the number of fields in the CSV file on this line
		csvFields = len(record)

		// if we are processing the first line - use the record field contents
		// as the SQL table column names - add to the temp string 'strbuffer'
		// use the tablename provided by the user
		// TODO : - add option to skip this line if user is adding data to an existing table?
		if lineCount == 0 {
			strbuffer.WriteString("PRAGMA foreign_keys=OFF;\nBEGIN TRANSACTION;\nCREATE TABLE " + tableName + " (")
			//strbuffer.WriteString("CREATE TABLE " + tableName + " (")

		}

		// if any line except the first one :
		// print the start of the SQL insert statement for the record
		// and  - add to the temp string 'strbuffer'
		// use the tablename provided by the user
		if lineCount > 0 {
			strbuffer.WriteString("INSERT INTO " + tableName + " VALUES (")
		}
		// loop through each of the csv lines individual fields held in 'record'
		// len(record) tells us how many fields are on this line - so we loop right number of times
		for i := 0; i < len(record); i++ {
			// if we are processing the first line used for the table column name - update the
			// record field contents to remove the characters: space | - + @ # / \ : ( ) '
			// from the SQL table column names. Can be overridden on command line with '-k true'
			if (lineCount == 0) && (keepOrigCols == false) {
				// for debug - output info so we can see current field being processed
				if debugSwitch {
					fmt.Printf("Running header clean up for '%s' ", record[i])
				}
				// call the function cleanHeader to do clean up on this field
				record[i] = cleanHeader(record[i])
				// for debug - output info so we can see any changes now made
				if debugSwitch {
					fmt.Printf("changed to '%s'\n", record[i])
				}
			}
			// if a csv record field is empty or has the text "NULL" - replace it with actual NULL field in SQLite
			// otherwise just wrap the existing content with ''
			// TODO : make sure we don't try to create a 'NULL' table column name?
			if len(record[i]) == 0 || record[i] == "NULL" {
				strbuffer.WriteString("NULL")
			} else {
				strbuffer.WriteString("\"" + record[i] + "\"")
			}
			// if we have not reached the last record yet - add a coma also to the output
			if i < len(record)-1 {
				strbuffer.WriteString(",")
			}
		}
		// end of the line - so output SQL format required ');' and newline
		strbuffer.WriteString(");\n")
		// line of SQL is complete - so push out to the new SQL file
		bWritten, err := sqlFileBuffer.WriteString(strbuffer.String())
		// check it wrote data ok - otherwise report the error giving the line number affected
		if (err != nil) || (bWritten != len(strbuffer.Bytes())) {
			fmt.Printf("WARNING: Error writing to SQL file line %d: %s", lineCount, err)
			return
		}
		// reset the string buffer - so it is empty ready for the next line to build
		strbuffer.Reset()
		// for debug - show the line number we are processing from the CSV file
		if debugSwitch {
			fmt.Print("..", lineCount)
		}
		// increment the line count - and loop back around for next line of the CSV file
		lineCount++
	}
	// END - reached the end of processing each line of the input CSV file
	//
	if debugSwitch {
		fmt.Println("\ncsv file processing complete - outputted to the new SQL file: ", sqlOutFile)
	}
	// finished processing the csv input file lines - so close off the SQL statements
	strbuffer.WriteString("COMMIT;\n")
	// write out final line to the SQL file
	bWritten, err := sqlFileBuffer.WriteString(strbuffer.String())
	// check it wrote data ok - otherwise report the error giving the line number affected
	if (err != nil) || (bWritten != len(strbuffer.Bytes())) {
		fmt.Printf("WARNING: Error outputting final line of the SQL file: line %d: %s", lineCount, err)
		return
	}
	if debugSwitch {
		fmt.Println("SQL file write complete")
	}
	fmt.Println("\nDONE\n\tCSV file processing complete, and the new SQL file format was written to: ", sqlOutFile)
	// finished the SQl file data writing - flush any IO buffers
	// NB below flush required as the data was being lost otherwise - maybe a bug in go version 1.2 only?
	sqlFileBuffer.Flush()
	// reset the string buffer - so it is empty as it is no longer needed
	strbuffer.Reset()
	// stop the timer for the SQL file creation process
	end := time.Now()

	// print out some stats about the csv file processed
	fmt.Println("\nSTATS\n\tCSV file", csvFileName, "has", lineCount, "lines with", csvFields, "CSV fields per record")
	fmt.Println("\tThe conversion took", end.Sub(start), "to run.\n\nAll is well.\n")

	defer pushtodb()
}
