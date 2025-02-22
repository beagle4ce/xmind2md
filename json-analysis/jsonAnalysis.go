package jsonanalysis

import (
	"encoding/json"
	errhandle "xmind2md/err-handle"
	"xmind2md/models"
)

/*
优先解析relationships块
*/
func JsonAnalysis(jsonBytes []byte) []models.Sheet {

	var sheets []models.Sheet
	// 解析json数据
	err := json.Unmarshal(jsonBytes, &sheets)
	errhandle.HandleError(err)

	return sheets
}
