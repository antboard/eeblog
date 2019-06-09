package handle

import (
	"net/http"

	"github.com/antboard/eeblog/model"

	"github.com/gin-gonic/gin"
)

// Index ...
func Index(c *gin.Context) {
	vbs := model.GetResentBlog(0)
	c.JSON(http.StatusOK, vbs)
	// c.String(http.StatusOK, "...")
}
