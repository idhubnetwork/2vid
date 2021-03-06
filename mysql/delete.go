package db_mysql

// Credential DELETE in mysql
func DeleteCredential(jwt_id int) error {
	result, err := DB_mysql.Exec("delete from credentials where jwt_id = ?",
		jwt_id)

	if err != nil {
		return err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

// Set Credential status is TO BE DELETED in mysql
func DeleteCredential_TBD(jwt_id, status int) error {
	result, err := DB_mysql.Exec("update credentials set status = ? where jwt_id = ?",
		status,
		jwt_id)

	if err != nil {
		return err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}
