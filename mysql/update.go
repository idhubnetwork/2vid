package db_mysql

func UpdateCredential_TBD(jwt_id int, credential *Credential) error {
	status := DEFAULT_STATUS | credential.Status

	result, err := DB_mysql.Exec(`insert into updated_credentials(iss,
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
		status,
		jwt_id)

	if err != nil {
		return err
	}
	_, err := result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func UpdateCredential(jwt_id int) error {
	result, err := DB_mysql.Exec(`update credentials inner join (select 
	iss,sub,aud,exp,nbf,iat,jti,net,ipfs,context,credential,status from 
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
	_, err := result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
