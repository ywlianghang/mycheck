package PDF

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
)

func (out *OutputWayStruct) pdfTitleModule(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	titleStr := "MySQL深度巡检报告"
	pdf.SetTitle(titleStr,false)
	pdf.SetHeaderFuncMode(func() {
		pdf.SetFont("simfang", "", 20)
		wd := pdf.GetStringWidth(titleStr) + 6
		pdf.SetY(0.6)            //先要设置 Y，然后再设置 X。否则，会导致 X 失效
		pdf.SetX((210 - wd) / 2) //水平居中的算法
		pdf.SetDrawColor(0, 80, 180)  //frame color
		pdf.SetFillColor(230, 230, 0) //background color
		pdf.SetTextColor(220, 50, 50) //text color
		pdf.SetLineWidth(1)
		pdf.CellFormat(wd, 10, titleStr, "1", 1, "CM", true, 0, "")
		//第 5 个参数，实际效果是：指定下一行的位置
		pdf.Ln(5)
	}, false)
	return pdf
}
//设置页眉页脚
func (out *OutputWayStruct) pdfHeaderModule(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("simfang", "", 8)
		pdf.SetTextColor(128, 128, 128)
		pdf.CellFormat(
			0, 5,
			fmt.Sprintf("Page %d", pdf.PageNo()),
			"", 0, "C", false, 0, "",
		)
	})
	return pdf
}

//设置一级标题
func (out *OutputWayStruct) pdfPrimaryTitleModule(pdf *gofpdf.Fpdf,titleStr string) *gofpdf.Fpdf {
	pdf.SetFont("simfang", "", 12)
	pdf.SetFillColor(200, 220, 255) //background color
	pdf.CellFormat(0, 6, titleStr, "", 1, "L", true, 0, "")
	pdf.Ln(2)
	return pdf
}