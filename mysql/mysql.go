package db_mysql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB_mysql *sql.DB

func init() {
	DB_mysql, err := sql.Open("mysql",
		"user:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal(err)
	}
}

// args is []string{iss, sub, aud, jti}
func GetCredential(args ...string) (credential string, err error) {
	row = DB_mysql.QueryRow("select credential from credentials where iss = ?, sub = ?, aud = ?, jti = ?", args)
	err = row.Scan(&credential)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(credential)
	return
}

// args is []string{iss, sub, aud}
func GetCredentials(args ...string) (credentials []string, err error) {
	rows, err = DB_mysql.Query("select credential from credentials where iss = ?, sub = ?, aud = ?, jti = ?", args)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err = rows.Scan(&credentials[i])
		if err != nil {
			log.Fatal(err)
		}
		i++
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return
}
