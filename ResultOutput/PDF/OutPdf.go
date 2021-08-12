package PDF

import (
	pub "DepthInspection/api/PublicClass"
	"fmt"
	"github.com/jung-kurt/gofpdf"
)
type OutputWayStruct struct{}
type OutPutWayInter interface {
	OutPdf()
	ResultSummaryStringSlice() [][]string
}

func (out *OutputWayStruct) pdfinit() *gofpdf.Fpdf{
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("simfang", "", "./lib/simfang.ttf")
	return pdf
}

func (out *OutputWayStruct) OutPdf() {
	//设置页面参数
	var pdf *gofpdf.Fpdf
	var tmpCheckTypeSliect []string
	pdf = out.pdfinit()

	pdf.AddPage()

	//标题头
	//pdf = out.pdfTitleModule(pdf)

	//标题一
	out.pdfPrimaryTitleModule(pdf,"一、巡检介绍")
	titlew := []float64{50.0, 50.0, 40.0, 40.0}
	var dc1 = []string{"巡检时间：",fmt.Sprintf("%s",pub.CheckBeginTime),"巡检人员：",pub.ResultOutput.InspectionPersonnel}
	var dc2 = []string{"巡检级别：",pub.ResultOutput.InspectionLevel,"巡检耗时(s)：",fmt.Sprintf("%d",pub.CheckTimeConsuming)}
	var dc = [][]string{dc1,dc2}
	pdf = out.pdfTableInsert(pdf,titlew,dc)
	pdf.CellFormat(0, 2, "", "", 1, "LM", false, 0, "")
	pdf.Ln(-1)
	//标题二
	out.pdfPrimaryTitleModule(pdf,"二、巡检结果概览")
	w := []float64{17.0, 70.0, 30.0, 38.0,30.0}  //定义每行表格的宽度
	c := []string{" ","检测项","检测数量","正常","异常"}
	pdf = out.pdfTableBodyColorsFormat(pdf,w,c)
	pdf = out.pdfTableInsert(pdf,w,out.ResultSummaryStringSlice())
	pdf.CellFormat(0, 2, "", "", 1, "LM", false, 0, "")
	pdf.Ln(-1)

	//标题三
	out.pdfPrimaryTitleModule(pdf,"三、巡检结果详情")  //TNP
	w3 := []float64{10,60.0,30.0,20.0, 70.0}  //定义每行表格的宽度
	cd := []string{" ","巡检项名称","阈值","错误码","异常相关信息"}
	pdf.SetFont("simfang", "", 10)
	pdf.MultiCell(0, 5, string("3.1 巡检数据库环境"), "", "", false)

	//子标题3.2内容
	pdf.MultiCell(0, 5, string("3.2 巡检数据库配置"), "", "", false)
	pdf = out.pdfTableBodyColorsFormat(pdf,w3,cd)
	//tmpResultConfig := out.tmpConfigCheckResultSummary("configParameter",pub.InspectionResult.DatabaseConfigCheck.ConfigParameter)
	tmpCheckTypeSliect = []string{"configParameter"}
	tmpResultConfig := out.tmpResultSummary(tmpCheckTypeSliect)
	pdf = out.pdfTableInsert(pdf,w3,tmpResultConfig)
	pdf.Ln(-1)

	//子标题3.3内容
	pdf.MultiCell(0, 5, string("3.3 巡检数据库性能"), "", "", false)
	pdf = out.pdfTableBodyColorsFormat(pdf,w3,cd)
	tmpCheckTypeSliect = []string{"binlogDiskUsageRate","historyConnectionMaxUsageRate","tmpDiskTableUsageRate",
		"tmpDiskfileUsageRate","innodbBufferPoolUsageRate","innodbBufferPoolDirtyPagesRate","innodbBufferPoolHitRate",
	"openFileUsageRate","openTableCacheUsageRate","openTableCacheOverflowsUsageRate","selectScanUsageRate","selectfullJoinScanUsageRate",
	"tableAutoPrimaryKeyUsageRate","tableRows","diskFragmentationRate","bigTable","coldTable"}
	tmpResultPerformance := out.tmpResultSummary(tmpCheckTypeSliect)
	pdf = out.pdfTableInsert(pdf,w3,tmpResultPerformance)
	pdf.Ln(-1)

	//子标题3.4内容
	pdf.MultiCell(0, 5, string("3.4 巡检数据库基线"), "", "", false)
	pdf = out.pdfTableBodyColorsFormat(pdf,w3,cd)
	tmpCheckTypeSliect = []string{"tableCharset","tableEngine","tableForeign","tableNoPrimaryKey","tableAutoIncrement",
		"tableBigColumns","indexColumnIsNull","indexColumnIsEnumSet","indexColumnIsBlobText","tableIncludeRepeatIndex",
	"tableProcedure","tableFunc","tableTrigger"}
	tmpResultBaselineResult := out.tmpResultSummary(tmpCheckTypeSliect)
	pdf = out.pdfTableInsert(pdf,w3,tmpResultBaselineResult)
	pdf.Ln(-1)

	//子标题3.5内容
	pdf.MultiCell(0, 5, string("3.5 巡检数据库安全"), "", "", false)
	pdf = out.pdfTableBodyColorsFormat(pdf,w3,cd)
	tmpCheckTypeSliect = []string{"anonymousUsers","emptyPasswordUser","rootUserRemoteLogin","normalUserConnectionUnlimited",
		"userPasswordSame","normalUserDatabaseAllPrivilages","normalUserSuperPrivilages","databasePort"}
	tmpResultUserSecurityResult := out.tmpResultSummary(tmpCheckTypeSliect)
	pdf = out.pdfTableInsert(pdf,w3,tmpResultUserSecurityResult)
	pdf.Ln(-1)

	pdf.MultiCell(0, 5, string("3.6 巡检数据库空间"), "", "", false)
	pdf.MultiCell(0, 5, string("3.7 巡检数据库备份"), "", "", false)
	pdf.Ln(-1)

	//将内容写入到pdf中
	if err := pdf.OutputFileAndClose(pub.ResultOutput.OutputPath+pub.ResultOutput.OutputFile); err != nil {
		panic(err.Error())
	}

}