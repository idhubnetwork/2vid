package db_mysql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB_mysql *sql.DB

type Credential struct {
	credential string `json:"credential"`
}

func init() {
	DB_mysql, err := sql.Open("mysql",
		"user:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal(err)
	}
	_ = DB_mysql
}

// args is []string{iss, sub, aud, jti}
func GetCredential(args ...string) (credential Credential, err error) {
	row := DB_mysql.QueryRow("select credential from credentials where iss = ?, sub = ?, aud = ?, jti = ?", args)
	err = row.Scan(&credential.credential)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(credential)
	return
}

// args is []string{iss, sub, aud}
func GetCredentials(args ...string) (credentials []Credential, err error) {
	rows, err := DB_mysql.Query("select credential from credentials where iss = ?, sub = ?, aud = ?", args)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err = rows.Scan(&credentials[i].credential)
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

func GetStatus(args ...string) (jwt_id int, status int, err error) {
	var row *sql.Row
	if len(args) == 3 {
		row = DB_mysql.QueryRow("select jwt_id, status from credentials where iss = ?, sub = ?, aud = ?", args)
	} else {
		row = DB_mysql.QueryRow("select jwt_id, status from credentials where iss = ?, sub = ?, aud = ?, jti = ?", args)
	}
	err = row.Scan(&jwt_id, &status)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(status)
	return
}
