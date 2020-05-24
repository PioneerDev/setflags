package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"set-flags/models"
)

// list all the flags
func ListFlags(c *gin.Context) {

	data := models.GetAllFlags()
	c.PureJSON(http.StatusOK, data)
}

// create a flag
func CreateFlag(c *gin.Context) {

	var flag map[string]interface{}

	if c.ShouldBind(&flag) == nil {
		models.CreateFlag(flag)
	}
	c.JSON(http.StatusCreated, gin.H{
		"code": 1,
		"msg":  "created flag",
	})
}

// Update an existing flag
func UpdateFlag(c *gin.Context) {
	//flagId := c.Param("id")
	//op := c.Param("op")
	//
	//if op != "yes" || op != "no" || op != "done" {
	//
	//}

}

// list all flags of the user
func FindFlagsByUserID(c *gin.Context) {
	userId := c.Param("id")

	flags := models.FindFlagsByUserID(userId)
	c.PureJSON(http.StatusOK, flags)
}

// upload evidence
func UploadEvidence(c *gin.Context) {
	flagId := c.Query("flag_id")
	attachmentId := c.Param("attachment_id")

	fmt.Sprintf("attachmentId: %s, flagId: %s", attachmentId, flagId)

	type_ := c.Query("type")

	if type_ != "image" && type_ != "audio" && type_ != "video" && type_ != "document" {
		c.JSON(http.StatusBadRequest, gin.H{
			"info": fmt.Sprintf("type: %s is invalid.", type_),
		})
		return
	}

	if !models.FLagExists(flagId) {
		c.JSON(http.StatusNotFound, gin.H{
			"info": "not found specific flag.",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"info": err,
		})
		return
	}

	log.Println(file.Filename)

	// Upload the file to specific dst.
	c.SaveUploadedFile(file, fmt.Sprintf("./%s", file.Filename))

	models.CreateEvidence(attachmentId, flagId, type_)

	c.JSON(http.StatusOK, gin.H{
		"info": fmt.Sprintf("'%s' uploaded!", file.Filename),
	})
}

// list all the evidences since yesterday
func ListEvidences(c *gin.Context) {
	flagId := c.Param("flag_id")

	data := models.FindEvidencesByFlag(flagId)

	c.PureJSON(http.StatusOK, data)
}
