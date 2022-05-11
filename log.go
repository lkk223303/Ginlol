package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.gmc303.com/convolvulus/backend_gateway/cores/admin"
	"gitlab.gmc303.com/convolvulus/backend_gateway/cores/system"
)

func OperationLog(c *gin.Context, note string) {
	// 紀錄操作異動

	var act string
	inURL := c.FullPath()

	op := c.MustGet(admin.TokenAdminInfo).(*admin.Info)
	sCore := system.New()
	os := &system.OperationSaveInfo{}
	os.OperationItem = system.OTAM
	os.SourceOperation = op.Account

	switch c.Request.Method {
	case "GET":
		act = "Read"
	case "POST":
		os.ActType = system.OTActInsert
		act = "Create"
	case "PUT":
		os.ActType = system.OTActUpdate
		act = "Update"
	case "DELETE":
		os.ActType = system.OTActDelete
		act = "Delete"
	}

	// 去掉參數
	uStr := strings.Split(inURL, "/")
	var r []string
	for _, k := range uStr {
		if strings.HasPrefix(k, ":") || strings.HasPrefix(k, "*") {
			k = ""
		}
		if k != "" {
			r = append(r, k)
		}
	}

	/*
		儲存格式 VERB resource properties
	*/

	os.Note = fmt.Sprintf("%s: %s %s, Note:%s", act, uStr[1], r[1:], note) //"Insert: chieftain [manual card info], Note:{id:1565454,name:"name"}"
	sCore.OperationSave(os)
}
