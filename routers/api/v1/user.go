package v1

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/gofrs/uuid"
	"net/http"
	"set-flags/models"
)

// check the total rewards received by the user for the flag
func CheckRewards(c *gin.Context) {
	userId := c.Param("user_id")
	flagId := c.Param("flag_id")
	userID, _ := uuid.FromString(userId)
	flagID, _ := uuid.FromString(flagId)
	data := models.FindEvidenceByFlagIdAndAttachmentId(flagID, userID)
	c.PureJSON(http.StatusOK, data)
}
