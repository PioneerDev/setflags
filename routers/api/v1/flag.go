package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
	"set-flags/models"
	"set-flags/pkg/cloud/aws"
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
	flagId := c.Param("id")

	_, err := uuid.FromString(flagId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"info": fmt.Sprintf("flagId: %s is not a valid UUID.", flagId),
		})
		return
	}

	op := c.Param("op")

	if op != "yes" && op != "no" && op != "done" {
		c.JSON(http.StatusBadRequest, gin.H{
			"info": fmt.Sprintf("op: %s is invalid.", op),
		})
		return
	}

	if !models.FLagExists(flagId) {
		c.JSON(http.StatusNotFound, gin.H{
			"info": "Flag not found.",
		})
		return
	}

	flag := models.FindFlagByID(flagId)

	if flag.Status != "done" {
		c.JSON(http.StatusBadRequest, gin.H{
			"info": "not yet upload the evidence.",
		})
		return
	}

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
	dst := fmt.Sprintf("./%s", file.Filename)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"info": err,
		})
		return
	}

	// todo
	// upload media to s3, need test in actual environment
	url, err := aws.S3Upload(dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"info": err,
		})
		return
	}

	attachmentID, _ := uuid.FromString(attachmentId)
	flagID, _ := uuid.FromString(flagId)
	models.CreateEvidence(attachmentID, flagID, type_, url)

	// update flag status to `done`
	models.UpdateFlagStatus(flagId, "done")

	c.JSON(http.StatusOK, gin.H{
		"info": fmt.Sprintf("'%s' uploaded!", file.Filename),
	})
}

// list all the evidences since yesterday
func ListEvidences(c *gin.Context) {
	flagId := c.Param("flag_id")
	flagID, _ := uuid.FromString(flagId)
	data := models.FindEvidencesByFlag(flagID)

	c.PureJSON(http.StatusOK, data)
}
