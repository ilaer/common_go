package office

import (
	"fmt"
	"github.com/ilaer/common_go/standard"
	"log"
	"path/filepath"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// @title    ExcelSheetExportAsPdf
// @description    把excel文件指定的sheet导出为pdf文件
// @param xlsxPath string excel文件路径,pdfPath string 存放生成的pdf的目录
// @return

func ExcelSheetsExportAsPdf(excelFilePath, pdfPath string, indexs []int) error {

	//CoInitialize是Windows提供的API函数，用来告诉 Windows以单线程的方式创建com对象。
	//应用程序调用com库函数（除CoGetMalloc和内存分配函数）之前必须初始化com库。
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	//创建com对象
	iunk, err := oleutil.CreateObject("Excel.Application")
	if err != nil {
		log.Printf("error creating Wps object: %s", err)
		return fmt.Errorf("error creating Wps object: %s", err)

	}
	defer iunk.Release()

	// 获取能够遍历的对象
	excel := iunk.MustQueryInterface(ole.IID_IDispatch)
	defer excel.Release()

	//获取工作簿对象
	workbooks := oleutil.MustGetProperty(excel, "Workbooks").ToIDispatch()
	defer workbooks.Release()

	//打开指定工作簿
	workbook := oleutil.MustCallMethod(workbooks, "Open", excelFilePath).ToIDispatch()
	defer workbook.Release()

	//遍历需要导出的sheet序号
	for _, index := range indexs {

		//获取指定sheet
		worksheet := oleutil.MustGetProperty(workbook, "Worksheets", index).ToIDispatch()
		defer worksheet.Release()

		//获取指定sheet的PageSetup
		PageSetup, err := oleutil.GetProperty(worksheet, "PageSetup")
		if err != nil {
			log.Printf("Worksheets(1) get PageSetup error: %s", err)
		} else {
			////设置指定sheet的PageSetup.FitToPagesTall
			_, err = oleutil.PutProperty(PageSetup.ToIDispatch(), "FitToPagesWide", 1)
			if err != nil {
				log.Printf("Worksheets(1).PageSetup put FitToPagesTall error: %s", err)
			}
			_, err = oleutil.PutProperty(PageSetup.ToIDispatch(), "LeftMargin", 0)
			if err != nil {
				log.Printf("Worksheets(1).PageSetup put LeftMargin error: %s", err)
			}
			_, err = oleutil.PutProperty(PageSetup.ToIDispatch(), "RightMargin", 0)
			if err != nil {
				log.Printf("Worksheets(1).PageSetup put RightMargin error: %s", err)
			}
			_, err = oleutil.PutProperty(PageSetup.ToIDispatch(), "TopMargin", 0)
			if err != nil {
				log.Printf("Worksheets(1).PageSetup put TopMargin error: %s", err)
			}
			_, err = oleutil.PutProperty(PageSetup.ToIDispatch(), "BottomMargin", 0)
			if err != nil {
				log.Printf("Worksheets(1).PageSetup put BottomMargin error: %s", err)
			}
		}

		// 导出为指定格式
		fileName, _ := standard.GetFileNameFromFilePath(excelFilePath)

		pdfFileName := fmt.Sprintf("%s_%d.pdf", fileName, index)
		pdfFilePath := filepath.Join(pdfPath, pdfFileName)
		oleutil.MustCallMethod(worksheet, "ExportAsFixedFormat", 0, pdfFilePath)

	}

	//todo 如果不保存wps会弹窗,必须先保存.暂时未找到解决方案
	oleutil.PutProperty(workbook, "Saved", true)
	oleutil.MustCallMethod(workbook, "Close")
	oleutil.MustCallMethod(excel, "Quit")
	return nil
}
