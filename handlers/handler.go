package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idhubnetwork/jsontokens"
)

const (
	// Read action need audience did authorization.
	READ_ERROR = "Only credential audience can read!"
)

// Handler Error Json
type ActionErr struct {
	ActionError string `json:"FaliedAction"`
}

// Handler 200 OK JSON
type ActionSuccess struct {
	Action string `json:"Action"`
}

// Distribute the request to the corresponding handler.
func HandleCredential(c *gin.Context) {
	tmp, ok := c.Get("DIDJsonToken")
	if !ok {
		c.JSON(http.StatusForbidden, ActionErr{"non DID Json Token"})
		return
	}
	jt := tmp.(*jsontokens.JsonToken)

	action, ok := jt.Get("action").(string)
	if !ok {
		c.JSON(http.StatusForbidden, ActionErr{"jsontoken non action"})
		return
	}

	switch action {
	case "READ":
		readCredential(c, jt)
	case "CREATE":
		createCredential(c, jt)
	case "UPDATE":
		updateCredential(c, jt)
	case "DELETE":
		deleteCredential(c, jt)
	default:
		c.JSON(http.StatusForbidden, ActionErr{"jsontoken invalid action"})
		return
	}
}
