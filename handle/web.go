package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Text ...
func Text(c *gin.Context) {
	c.String(http.StatusOK, "...")
}
