package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// TXHandler is handler for working with transaction.
// This is wrapper function for commit and rollback.
func TXHandler(c *gin.Context, db *sqlx.DB, f func(*sqlx.Tx) error) {
	tx, err := db.Beginx()
	if err != nil {
		c.JSON(500, gin.H{"err": "start transaction failed"})
		c.Abort()
		return
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			log.Print("rollback operation.")
			return
		}
	}()
	if err := f(tx); err != nil {
		log.Printf("operation failed: %s", err)
		c.JSON(500, gin.H{"err": "operation failed"})
		c.Abort()
		return
	}
}
