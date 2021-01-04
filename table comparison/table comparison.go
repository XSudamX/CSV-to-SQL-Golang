package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := sql.Open("mysql", "sudam@localhost:sudam@tcp(192.168.148.129:3306)/test")
	if err != nil {
		fmt.Println("Failed connection to database", err)
		panic(err.Error())
	}
	defer db.Close()

	result, err := db.Exec("SELECT batch_transaction_id, batch_id, invoice_id, amount_gross FROM batch_transaction_test WHERE NOT EXISTS ( SELECT * FROM csvtable WHERE batch_transaction_test.invoice_id=csvtable.invoice_id )")
	fmt.Println("Succeded; ", result)
	if err != nil {
		fmt.Println("Error in SQL Execution : ", err)
	}
}
SELECT
    *
FROM
    batch_transaction_test batch_transaction_id,
    batch_id,
    invoice_id,
    amount_gross
LEFT JOIN csvtable * ON
    batch_transaction_test.invoice_id = csvtable.invoice_id