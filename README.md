# 2vid: To Verifiable Identity
2vid is a storage and management server for Verifiable Credentails.

## A request to response

`request` --> `did_authentication` --> `actions_distribute` --> `action_authority` --> `credential_permission` --> `database_crud` --> `response`

### did authentication and authorization
Server authrnticate and did authorize according to following json format.
Use `jsontokens` to generate json string:
```json
{
	"did":"did:idhub:0x1234567890exampleEthereumAddress",
	"action":"CREATE || READ || UPDATE || DELETE",
	"destination":"server handler router url",
	"expiration":"Unix timestamp indicates token expiration",
	"jwt_iss":"did:idhub:0x1234567890exampleEthereumAddress",
	"jwt_sub":"credential subject defined by did application",
	"jwt_aud":"did:idhub:0x1234567890exampleEthereumAddress",
	"jwt_jti":"credential unique number from did application but optional"
}
```
Use `jsontokens` sign the json message, get the request token:
```go
type JsonToken struct {
	ClaimJson string                 `json:"msg"`
	Signature string                 `json:"sig"`
}
```
Package the token into the request's header field `Authentication` in the following format:
```
DIDJsonToken {{TokenString}}
```
Or use `form`, `json`, `xml` to pass in the request parameters as follows:
```go
type Token struct {
	Token string `form:"token" json:"token" xml:"token" binding:"required"`
}
```

### credential permission status
Credential `status` field indicates CRUD permission in database.

|     |Reserved|Reserved|DELETED|UPDATED|ISS_UPDATE|ISS_DELETE|AUD_UPDATE|AUD_DELETE|
|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|:--------:|
|1|||allowed but no yet deleted|allowed but no yet updated|iss could update|iss could delete OR confirm deleted|update need aud agree|aud could delete OR confirm deleted|
|0|||forbid to delete OR have been deleted and to be confirmed|have been updated and to be confirmed|iss can not update|iss can not delete OR have been deleted|update don't need aud agree|aud can not delete OR have been deleted||


Status for DELETE in database:
* `0011 1010` Both iss and aud can not delete.
* `0001 1011` Deleted by iss and to be confirmed by aud.
* `0001 1110` Deleted by aud and to be confirmed by iss.
* `0011 1011` Can deleted by aud directly.
* `0011 1110` Can deleted by iss directly.

Status for UPDATE in database:
* `0011 0101` Both iss and aud can not update.
* `0011 1101` Can updated by iss directly.
* `0010 0111` Updated by iss and  to be confirmed by aud.

### response and error
Reference to [GoDoc](https://godoc.org/github.com/idhubnetwork/2vid/handlers)

## Logger

JSON

### log monitoring

Elasticsearch

### log level

DEBUG[only test or develop], INFO, WARN, ERROR, FATAL

### log out file

* 2vid.out [INFO]
* 2vid_verify_failed.out [WARN]
* 2vid_database_error.out [ERROR]
* 2vid_panic.out [FATAL]

### log message

* http request and response [INFO] : statusCode|latency|clientIP|method|path|comment
* authentication middleware error [WARN] : statusCode|latency|clientIP|method|path|error|jsontoken
* database error [ERROR] : statusCode|latency|clientIP|method|path|comment|action|jsontoken|credential
* panic recover [FATAL] : statusCode|latency|clientIP|method|path|comment

## Redis Design

### Credential CRUD

`authentication` --> `redis hash`/`mysql select` --> `status`/`jwt_id`
hash key:
`sha3(jwt_iss + jwt_sub + jwt_aud)`
hash value:
```json
{
	"status"     : "permission int",
	"jwt_id"     : "mysql id",
	"credential" : "json web token"
}
```

`authentication` --> `redis pub/sub` --> `redis hash delete` --> `response` --> `mysql operation` --> `redis hash add`
5 channles:
* update_tbd [jwt_id, status, credential]
* update     [jwt_id]
* delete_tbd [jwt_id, status]
* delete     [jwt_id]
* create     [credential]