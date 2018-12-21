package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idhubnetwork/jsontokens"
)

const (
	DELETE_ISSUER_ERROR = "Credential issuer can delete but no authorization!"

	DELETE_AUDIENCE_ERROR = "Credential audience can delete but no authorization!"

	UPDATE_ERROR = "Only credential issuer can update!"

	UPDATE_NEED_ERROR = "Need credential audience agree update!"

	// Read action need audience did authorization.
	READ_ERROR = "Only credential audience can read!"

	// 0011 0000
	DEFAULT_STATUS = 0x30

	// 0001 1010
	UPDATE_ISSUER_OP = 0x1a

	// 1110 0111
	UPDATE_ISSUER_OP_TBD = 0xe7

	// 0001 1000
	UPDATE_AUDIENCE_OP = 0x18

	// 0001 1010 & 0011 0101 = 0001 0000 0x10
	IF_CAN_NOT_UPDATE = 0x1a

	// 0001 0000
	CAN_NOT_UPDATE = 0x10

	// 0000 0001
	DELETE_ISSUER_OP = 0x01

	// 0000 0100
	DELETE_AUDIENCE_OP = 0x04

	// 1101 1011
	DELETE_ISSUER_OP_TBD = 0xdb

	// 1101 1110
	DELETE_AUDIENCE_OP_TBD = 0xde

	// 0010 0101 & 0011 1010 = 0010 0000 0x20
	IF_CAN_NOT_DELETE = 0x25

	// 0010 0000
	CAN_NOT_DELETE = 0x20
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
