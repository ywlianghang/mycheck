package PDF

import (
	"github.com/jung-kurt/gofpdf"
)
type OutputWayStruct struct{}
type OutPutWayInter interface {
	OutPdf()
	ResultSummaryStringSlice() [][]string
}

func (out *OutputWayStruct) pdfinit() *gofpdf.Fpdf{
	pdf := gofpdf.New("P", "mm", "A4", "")
	//添加一页
	//pdf.AddPage()
	//将字体加载进来
	//AddUTF8Font("给字体起个别名", "", "fontPath")
	pdf.AddUTF8Font("simfang", "", "./lib/simfang.ttf")
	return pdf
}

func (out *OutputWayStruct) OutPdf() {
	//设置页面参数
	var pdf *gofpdf.Fpdf
	w := []float64{17.0, 70.0, 30.0, 38.0,30.0}
	titlew := []float64{50.0, 50.0, 40.0, 40.0}
	pdf = out.pdfinit()
	//主体--文字段落
	//chapterBody := func(bodyStr string) {
	//	pdf.SetFont("simfang", "", 8)
	//	//输出对齐文本
	//	pdf.MultiCell(0, 5, string(bodyStr), "", "", false)
	//	pdf.Ln(-1)
	//}

	pdf.AddPage()
	out.pdfPrimaryTitleModule(pdf,"一、巡检介绍")
	var dc1 = []string{"巡检时间：","2021-8-3 16:23:00","巡检人员：","golang"}
	var dc2 = []string{"巡检级别：","重保巡检","巡检耗时(h)：","2h"}
	var dc = [][]string{dc1,dc2}
	pdf = out.pdfTableInsert(pdf,titlew,dc)
	pdf.Ln(-1)
	out.pdfPrimaryTitleModule(pdf,"二、巡检结果概览")
	c := []string{" ","检测项","检测数量","正常","异常"}

	pdf = out.pdfTableBodyColorsFormat(pdf,w,c)
	pdf = out.pdfTableInsert(pdf,w,out.ResultSummaryStringSlice())

	pdf.CellFormat(0, 8, "", "", 1, "LM", false, 0, "")
	pdf.Ln(-1)
	out.pdfPrimaryTitleModule(pdf,"三、巡检结果详情")
	//pdf = out.pdfTableBodyFormat(pdf,out.ResultSummaryStringSlice())

	pdf.Ln(-1)
	//pdf.CellFormat(0, 10, "MySQL深度巡检报告", "0", 1, "CM", false, 0, "")
	//正文标题一
	//pdf.CellFormat(0, 10, "一、巡检时间", "0", 1, "LM", false, 0, "")
	//pdf.SetFont("simfang", "", 15)

	if err := pdf.OutputFileAndClose("MySQL DepthInspection Result Report.pdf"); err != nil {
		panic(err.Error())
	}

}