package db_mysql

import (
	"fmt"
	"log"
)

// Credential READ in mysql.
// args is []string{iss, sub, aud, jti}
func GetCredential(args ...string) (credential Credential, err error) {
	row := DB_mysql.QueryRow("select credential from credentials where iss = ?, sub = ?, aud = ?, jti = ?", args)
	err = row.Scan(&credential.Credential)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(credential)
	return
}

// Credentials array READ in mysql.
// args is []string{iss, sub, aud}
func GetCredentials(args ...string) (credentials []Credential, err error) {
	rows, err := DB_mysql.Query("select credential from credentials where iss = ?, sub = ?, aud = ?", args)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err = rows.Scan(&credentials[i].Credential)
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
