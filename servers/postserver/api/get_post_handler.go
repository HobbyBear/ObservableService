package api

import (
	"ObservableService/pkg/logger"
	"ObservableService/pkg/monitor"
	"ObservableService/services/postservice/export"
	"encoding/json"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func GetPostHandler(resp http.ResponseWriter, req *http.Request) {

	monitor.TestSetGauge(100)
	postResp, err := export.GetPost(req.Context(), "12331")
	if err != nil {
		log.Println(err)
		return
	}
	resMap := make(map[string]interface{})
	resMap["code"] = 200
	resMap["uid"] = postResp.Uid
	resMap["text"] = postResp.Text

	data, err := json.Marshal(resMap)
	if err != nil {
		logger.Error("marsh fail", zap.Error(err))
		return
	}
	resp.Write(data)
}
