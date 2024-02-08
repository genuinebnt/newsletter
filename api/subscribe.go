package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rs/zerolog/log"
)

type Form struct {
	Name  string `form:"name" binding:"required"`
	Email string `form:"email" binding:"required"`
}

func Subscribe(c *gin.Context) {
	db := c.MustGet("db").(*pgxpool.Pool)

	var form Form

	err := c.Bind(&form)
	if err != nil {
		log.Error().Msgf("Failed to bind to form %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	log.Info().Msgf("Adding '%s', '%s' as new subscriber", form.Email, form.Name)

	id, err := uuid.NewRandom()
	if err != nil {
		log.Error().Msgf("Failed to create random id: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(context.Background(), "INSERT INTO subscriptions (id, email, name, subscribed_at) VALUES ($1, $2, $3, $4)", id, form.Email, form.Name, time.Now().UTC())
	if err != nil {
		log.Error().Msgf("Failed to execute query: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	log.Info().Msg("New subscriber details have been saved")
}
