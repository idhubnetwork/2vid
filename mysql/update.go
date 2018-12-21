package db_mysql

import (
	"2vid/logger"
)

// Set Credential status is TO BE UPDATED and storage NEW Credential in
//  mysql and wait for agree from audience.
func UpdateCredential_TBD(jwt_id, status int, credential *Credential) error {
	new_status := DEFAULT_STATUS | credential.Status

	tx, err := DB_mysql.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	logger.Log.Debug(new_status)

	result, err := tx.Exec(`insert into updated_credentials(iss,
	sub,aud,exp,nbf,iat,jti,net,ipfs,context,credential,status,jwt_id) 
	values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		credential.Iss,
		credential.Sub,
		credential.Aud,
		credential.Exp,
		credential.Nbf,
		credential.Iat,
		credential.Jti,
		credential.Net,
		credential.IPFS,
		credential.Context,
		credential.Credential,
		new_status,
		jwt_id)

	if err != nil {
		return err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	result, err = tx.Exec("update credentials set status = ? where jwt_id = ?",
		status,
		jwt_id)

	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	logger.Log.Debug(id)
	return nil
}

// Credential UPDATE in mysql
func UpdateCredential(jwt_id int) error {
	tx, err := DB_mysql.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`update credentials inner join (select 
	iss,sub,aud,exp,nbf,iat,jti,net,ipfs,context,credential,status,jwt_id from 
	updated_credentials where jwt_id = ?) tmp on credentials.jwt_id = 
	tmp.jwt_id set credentials.iss = tmp.iss, credentials.sub = tmp.sub, 
	credentials.aud = tmp.aud, credentials.exp = tmp.exp, credentials.nbf = tmp.nbf, 
	credentials.iat = tmp.iat, credentials.jti = tmp.jti, credentials.net = tmp.net, 
	credentials.ipfs = tmp.ipfs, credentials.context = tmp.context, 
	credentials.credential = tmp.credential, credentials.status = tmp.status`,
		jwt_id)

	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	result, err = tx.Exec("delete from updated_credentials where jwt_id = ?",
		jwt_id)

	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
