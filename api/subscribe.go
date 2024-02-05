package api

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Form struct {
	Name  string `form:"name" binding:"required"`
	Email string `form:"email" binding:"required"`
}

func Subscribe(c *gin.Context) {
	db := c.MustGet("db").(*pgxpool.Pool)

	var form Form
	_ = c.Bind(&form)

	id, err := uuid.NewRandom()
	if err != nil {
		log.Println(err)
		return
	}
	_, err = db.Exec(context.Background(), "INSERT INTO subscriptions (id, email, name, subscribed_at) VALUES ($1, $2, $3, $4)", id, form.Email, form.Name, time.Now().UTC())
	log.Println(err)
}
