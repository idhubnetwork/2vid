package main

import (
	"github.com/gin-gonic/gin"
)

func handleCredential(c *gin.Context) {
	jt, ok := c.Get(jt)
	if !ok {}
	did, ok := jt.Get("did").(string)
	if !ok || len(did) != 32 {}
	action, ok := jt.Get("action").(string)
	if !ok {}
	jwt-iss, ok := jt.Get("jwt-iss").(string)
	if !ok || len(jwt-iss) != 32 {}
	jwt-aud, ok := jt.Get("jwt-aud").(string)
	if !ok || len(jwt-aud) != 32 {}
	jwt-sub, ok := jt.Get("jwt-sub").(string)
	if !ok {}
	switch action {
		case "READ": readCredential(c, jt)
		case "CREATE": createCredential(c, jt)
		case "UPDATE": updateCredential(c, jt)
		case "DELETE": deleteCredential(c, jt)
		default : return
	}
}
