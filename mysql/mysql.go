package db_mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/idhubnetwork/jsontokens"
	"gopkg.in/go-playground/validator.v9"
)

var DB_mysql *sql.DB

type Credential struct {
	Iss        string `json:"-" validate:"required,len=32,contains=did:idhub:0x"`
	Aud        string `json:"-" validate:"required,len=32,contains=did:idhub:0x"`
	Sub        string `json:"-" validate:"required"`
	Exp        int    `json:"-" validate:"required"`
	Nbf        int    `json:"-"`
	Iat        int    `json:"-"`
	Jti        string `json:"-"`
	Net        string `json:"-" validate:"required,contains=eth"`
	IPFS       string `json:"-"`
	Context    string `json:"-"`
	Status     int    `json:"-" validate:"required,lte=15,gte=0"`
	Credential string `json:"credential" validate:"required"`
}

func init() {
	DB_mysql, err := sql.Open("mysql",
		"user:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal(err)
	}
	_ = DB_mysql
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

func VerifyWritedData(did string, jwt string) (*Credential, error) {
	var credential *Credential
	tmp := jsontokens.NewJWT()
	err := tmp.SetJWT(jwt)
	if err != nil {
		return nil, errors.New("invalid jwt to init")
	}

	err = tmp.Verify()
	if err != nil {
		return nil, errors.New("invalid jwt signature")
	}

	if did != tmp.Get("iss").(string) {
		return nil, errors.New("only jwt issuer have opration permission")
	}

	var ok bool
	credential.Iss, ok = tmp.Get("iss").(string)
	if !ok {
		return nil, errors.New("credential must have valid issuer")
	}
	credential.Aud, ok = tmp.Get("aud").(string)
	if !ok {
		return nil, errors.New("credential must have valid audience")
	}
	credential.Sub, ok = tmp.Get("sub").(string)
	if !ok {
		return nil, errors.New("credential must have valid subject")
	}
	credential.Exp, ok = tmp.Get("exp").(int)
	if !ok {
		return nil, errors.New("credential must have valid expiration")
	}
	credential.Net, ok = tmp.Get("net").(string)
	if !ok {
		return nil, errors.New("credential must have valid blockchain net id")
	}
	credential.Status, ok = tmp.Get("status").(int)
	if !ok {
		return nil, errors.New("credential must have valid permission status")
	}

	credential.Nbf, ok = tmp.Get("nbf").(int)
	credential.Iat, ok = tmp.Get("iat").(int)
	credential.Jti, ok = tmp.Get("jti").(string)
	credential.IPFS, ok = tmp.Get("ipfs").(string)
	credential.Context, ok = tmp.Get("context").(string)
	credential.Credential = jwt

	var validate *validator.Validate
	validationErr := validate.Struct(credential)
	if validationErr != nil {
		return nil, validationErr
	}
	return credential, nil
}
