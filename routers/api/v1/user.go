package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"set-flags/models"
)

// check the total rewards received by the user for the flag
func CheckRewards(c *gin.Context) {
	userId := c.Param("user_id")
	flagId := c.Param("flag_id")

	data := models.FindEvidenceByFlagIdAndAttachmentId(flagId, userId)
	c.PureJSON(http.StatusOK, data)
}
