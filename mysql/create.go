package db_mysql

const (
	// 0011 0000
	DEFAULT_STATUS = 0x30
)

// Credential CREATE in mysql.
// Insert credential to mysql, param is a pointer reference Credential Struct.
func CreateCredential(credential *Credential) error {
	status := DEFAULT_STATUS | credential.Status

	result, err := DB_mysql.Exec(`insert into credentials(iss,
	sub,aud,exp,nbf,iat,jti,net,ipfs,context,credential,status) 
	values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
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
		status)

	if err != nil {
		return err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}
