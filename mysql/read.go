package db_mysql

// Credential READ in mysql.
// args is []string{iss, sub, aud, jti}
func GetCredential(args ...string) (*Credential, error) {
	credential := new(Credential)
	row := DB_mysql.QueryRow("select credential from credentials where iss=? and sub=? and aud=? and jti=?",
		args[0], args[1], args[2], args[3])
	err := row.Scan(&credential.Credential)
	if err != nil {
		return nil, err
	}
	return credential, nil
}

// Credentials array READ in mysql.
// args is []string{iss, sub, aud}
func GetCredentials(args ...string) (credentials []*Credential, err error) {
	rows, err := DB_mysql.Query("select credential from credentials where iss=? and sub=? and aud=?",
		args[0], args[1], args[2])
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tmp string
	i := 0
	for rows.Next() {
		err = rows.Scan(&tmp)
		credential := new(Credential)
		credential.Credential = tmp
		credentials = append(credentials, credential)
		if err != nil {
			return nil, err
		}
		i++
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return
}
