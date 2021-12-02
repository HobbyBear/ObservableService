package api

import (
	"ObservableService/services/postservice/export"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetPostHandler(c *gin.Context) {

	postResp, err := export.GetPost(c.Request.Context(), "12331")
	if err != nil {
		log.Println(err)
		return
	}
	resMap := make(map[string]interface{})
	resMap["code"] = 200
	resMap["uid"] = postResp.Uid
	resMap["text"] = postResp.Text

	c.JSON(http.StatusOK, resMap)
}
