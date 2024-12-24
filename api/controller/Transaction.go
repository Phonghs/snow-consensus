package controller

import (
	"SnowConsensus/p2p"
	"SnowConsensus/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (b *Base) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
		"id":      b.Node.GetID(),
	})
}

func (b *Base) CreateTransaction(c *gin.Context) {
	log.Println("CreateTransaction in node", b.Node.GetID())
	transaction, err := b.Node.CreateTransaction(util.GenerateRandomString(10))
	log.Println("CreateTransaction done in node", b.Node.GetID())
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": transaction,
		})
	}
}

func (b *Base) VerifyTransaction(c *gin.Context) {
	/*
		TODO: Implement rate limiting for each node to prevent overloading
	*/

	var message p2p.Message

	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, p2p.Message{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	query, err := b.Node.ReceiveQuery(message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, p2p.Message{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, query)
}
