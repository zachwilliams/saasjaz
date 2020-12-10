package core

import (
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Router -- public pointer to http router (in this case gin engine)
var Router *gin.Engine

// API -- public pointer to http router group (api endpoints are attached to this group from /api/init.go)
var API *gin.RouterGroup

func init() {
	Router = gin.Default()

	Router.Use(static.Serve("/*", static.LocalFile("./frontend/build", true)))

	Router.Use(corsMiddleware)
	Router.Use(errorMiddleware)

	API = Router.Group("/api")
}

func errorMiddleware(c *gin.Context) {
	c.Next()

	ginErr := c.Errors.Last()
	if ginErr == nil {
		return
	}

	switch err := ginErr.Err.(type) {
	default:
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}
}

func corsMiddleware(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	if origin == "" {
		return
	}

	if c.Request.Method == http.MethodOptions {
		h := c.Writer.Header()
		h.Set("Access-Control-Allow-Origin", origin)
		h.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,HEAD")
		h.Set("Access-Control-Allow-Headers", "authorization,content-type")
		h.Set("Access-Control-Max-Age", "86400")
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	h := c.Writer.Header()
	h.Set("Access-Control-Allow-Origin", origin)
}
