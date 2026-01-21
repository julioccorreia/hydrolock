package router

import (
	"github.com/gin-gonic/gin"
	"github.com/julioccorreia/hydrolock/internal/adapters/http/handlers"
)

func NewRouter(waterHandler *handlers.WaterHandler) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		v1.POST("/intake", waterHandler.Register)
	}

	return r
}
