package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/idhubnetwork/jsontokens"
)

const (
	// 0011 0000
	DEFAULT_STATUS = 0x30
)

func CreateCredential(c *gin.Context, jt *jsontokens.JsonToken) {}
