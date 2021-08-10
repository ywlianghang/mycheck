package main

import (
	"DepthInspection/InspectionItem"
	"DepthInspection/ResultOutput/PDF"
	"DepthInspection/api/PublicClass"
)

var OutputPdf PDF.OutPutWayInter =  &PDF.OutputWayStruct{}

func main() {
	//配置文件初始化
	PublicClass.ConfigInit()
	//查询数据初始化，将用到的数据初始化到内存中
	PublicClass.QueryDbDateInit()
	InspectionItem.DatabaseConfigCheck(PublicClass.Ccc)
	var c = &InspectionItem.DatabaseBaselineCheckStruct{}
	c.BaselineCheckTablesDesign()
	c.BaselineCheckColumnsDesign()
	c.BaselineCheckIndexColumnDesign()
	c.BaselineCheckProcedureTriggerDesign()
	c.BaselineCheckUserPriDesign()
	c.BaselineCheckPortDesign()
	c.DatabasePerformanceStatusCheck()
	c.DatabasePerformanceTableIndexCheck()
	OutputPdf.OutPdf()
}
