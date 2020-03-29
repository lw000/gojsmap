package routers

import (
	"gojsmap/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterService(engine *gin.Engine) {
	api := engine.Group("/api")
	api.GET("/log", logHandler)
}

func logHandler(c *gin.Context) {
	module := c.Query("module")
	if module == "" {
		c.JSON(http.StatusOK, gin.H{"c": 1, "m": "module field is empty", "d": gin.H{}})
		return
	}

	line := c.Query("line")
	if line == "" {
		c.JSON(http.StatusOK, gin.H{"c": 2, "m": "line field is empty", "d": gin.H{}})
		return
	}

	column := c.Query("column")
	if column == "" {
		c.JSON(http.StatusOK, gin.H{"c": 3, "m": "column field is empty", "d": gin.H{}})
		return
	}

	info := c.Query("info")
	if info == "" {
		c.JSON(http.StatusOK, gin.H{"c": 4, "m": "info field is empty", "d": gin.H{}})
		return
	}

	cfg := models.NewLogConfig()
	cfg.SetModule(module)
	cfg.SetLine(line)
	cfg.SetColumn(column)
	cfg.SetInfo(info)
	if err := cfg.Save(c.ClientIP()); err != nil {
		c.JSON(http.StatusOK, gin.H{"c": 5, "m": err.Error(), "d": gin.H{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"c": 0, "m": "", "d": gin.H{}})
}
