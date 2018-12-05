# 2vid: To Verifiable Identity
2vid is a storage and management server for Verifiable Credentails.

## A request to response

`request` --> `did_authentication` --> `actions_distribute` --> `action_authority` --> `credential_permission` --> `database_crud` --> `response`

### did authentication and authorization
Server authrnticate and did authorize according to following json format.
Use `jsontokens` to generate json string:
```json
{
did:"did:idhub:0x1234567890exampleEthereumAddress",
action:"CREATE || READ || UPDATE || DELETE",
destination:"server handler router url",
expiration:"Unix timestamp indicates token expiration",
jwt_iss:"did:idhub:0x1234567890exampleEthereumAddress",
jwt_sub:"credential subject defined by did application",
jwt_aud:"did:idhub:0x1234567890exampleEthereumAddress",
jwt_jti:"credential unique number from did application but optional"
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
