package api

import (
	"github.com/gin-gonic/gin"
)

type Form struct {
	Name  string `form:"name" binding:"required"`
	Email string `form:"email" binding:"required"`
}

func Subscribe(c *gin.Context) {
	var form Form
	err := c.Bind(&form)
	if err == nil {
		c.JSON(200, gin.H{"name": form.Name, "email": form.Email})
	}
}
