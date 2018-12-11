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

// Struct Credential is an jwt object for storage and validation.
//
// Iss, Aud only idhub-did and length-32 string
// Sub, Jti is Unique Identification for the two same entities
// Exp is credential expiration for validate effectiveness
// Net is blockchain network id, only supported ethereum
// Credential is jwt form credential
type Credential struct {
	Iss        string `json:"-" validate:"required,len=52,contains=did:idhub:0x"`
	Aud        string `json:"-" validate:"required,len=52,contains=did:idhub:0x"`
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

// mysql init, close in package main
func init() {
	var err error
	DB_mysql, err = sql.Open("mysql",
		"root:@tcp(127.0.0.1:3306)/2vid_test")
	if err != nil {
		log.Fatal(err)
	}
}

// select credential.status and credential.jwt_id from mysql
// params args is []string{iss, aud, sub, jwt}, jti is optional
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

// VerifyWritedData is a validator for credential when update and create
//  in mysql.
// Only issuer can create and update credential in mysql.
func VerifyWritedData(did string, jwt string) (*Credential, error) {
	var credential = new(Credential)
	tmp := jsontokens.NewJWT()
	err := tmp.SetJWT(jwt)
	if err != nil {
		return nil, errors.New("invalid jwt to init")
	}

	err = tmp.Verify()
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("invalid jwt signature")
	}
	fmt.Println(tmp)

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
	expiration, ok := tmp.Get("exp").(float64)
	if !ok {
		return nil, errors.New("credential must have valid expiration")
	}
	credential.Exp = int(expiration)
	credential.Net, ok = tmp.Get("net").(string)
	if !ok {
		return nil, errors.New("credential must have valid blockchain net id")
	}
	status, ok := tmp.Get("status").(float64)
	if !ok {
		return nil, errors.New("credential must have valid permission status")
	}
	credential.Status = int(status)

	credential.Nbf, ok = tmp.Get("nbf").(int)
	credential.Iat, ok = tmp.Get("iat").(int)
	credential.Jti, ok = tmp.Get("jti").(string)
	credential.IPFS, ok = tmp.Get("ipfs").(string)
	credential.Context, ok = tmp.Get("context").(string)
	credential.Credential = jwt

	validate := validator.New()
	err = validate.Struct(credential)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return nil, errors.New("validate failed")
		}

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
			fmt.Println(err.StructField())     // by passing alt name to ReportError like below
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
		return nil, errors.New("validate failed")
	}
	fmt.Println("validate success\n")
	return credential, nil
}
