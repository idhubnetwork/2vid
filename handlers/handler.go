package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idhubnetwork/jsontokens"
)

const (
	UPDATE_ERROR = "Only credential issuer can update!"

	UPDATE_NEED_ERROR = "Need credential audience agree update!"

	// Read action need audience did authorization.
	READ_ERROR = "Only credential audience can read!"

	// 0011 0000
	DEFAULT_STATUS = 0x30

	// 0000 0010
	UPDATE_ISSUER_OP = 0x02

	// 1110 0111
	UPDATE_ISSUER_OP_TBD = 0xe7

	// 0001 1000
	UPDATE_AUDIENCE_OP = 0x18

	// 0001 1010 & 0011 0101 = 0001 0000 0x10
	IF_CAN_NOT_UPDATE = 0x1a

	// 0001 0000
	CAN_NOT_UPDATE = 0x10
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
