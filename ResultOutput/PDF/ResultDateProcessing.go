package PDF

import (
	pub "DepthInspection/api/PublicClass"
	"fmt"
	"strconv"
)

func (out *OutputWayStruct) tmpaa(checkNum,checkType string,ast []map[string]string) []string{
	var abnormalCount = 0
	var normalCount = 0
	var checkNumberTotal string
	var aa []string
	for v := range ast{
		if ast[v]["checkStatus"] == "abnormal" && ast[v]["checkType"] == checkType{
			abnormalCount++
		}
		if ast[v]["checkStatus"] == "normal" && ast[v]["checkType"] == checkType {
			normalCount++
		}
	}
	checkNumberTotal = strconv.Itoa(abnormalCount+normalCount)
	aa = []string{checkNum,checkType,checkNumberTotal,strconv.Itoa(normalCount),strconv.Itoa(abnormalCount)}
	return aa
}
//配置参数结果汇总
//func (out *OutputWayStruct) tmpConfigCheckResultSummary(checkType string,ast []map[string]string) [][]string{
//	var bb [][]string
//	var tmpThreshold,tmpCheckValue,tmpCheckName string
//	var tmpEque = 0
//	for v := range ast{
//		if ast[v]["checkStatus"] == "abnormal" && ast[v]["checkType"] == checkType{
//			var aa []string
//			tmpEque ++
//			tmpCheckName = ast[v]["configVariableName"]
//			tmpThreshold = ast[v]["configValue"]
//			tmpCheckValue = fmt.Sprintf("%s=%s",ast[v]["configVariableName"],ast[v]["configVariable"])
//			aa = append(aa,strconv.Itoa(tmpEque))
//			aa = append(aa,tmpCheckName)
//			aa = append(aa,tmpThreshold)
//			aa = append(aa," ")
//			aa = append(aa,tmpCheckValue)
//			bb = append(bb,aa)
//		}
//	}
//	return bb
//}

func (out *OutputWayStruct) tmpcc(checkRulest []map[string]string) []string {
	var bc []string
	var tmpCheckType, tmpThreshold, tmpErrorCode, tmpAbnormalInformation = "","","",""
	var tmpeq int
	for i := range checkRulest {
		if checkRulest[i]["checkStatus"] == "abnormal" {
			tmpCheckType = checkRulest[i]["checkType"]
			tmpThreshold = checkRulest[i]["threshold"]
			tmpErrorCode = checkRulest[i]["errorCode"]
			tmpeq++
			if tmpeq > 1 {
				tmpAbnormalInformation = fmt.Sprintf("%s 等", tmpAbnormalInformation)
				break
			}
			if tmpAbnormalInformation != "" {
				tmpAbnormalInformation = fmt.Sprintf("%s,%s", tmpAbnormalInformation, checkRulest[i]["currentValue"])
			} else {
				tmpAbnormalInformation = fmt.Sprintf("%s", checkRulest[i]["currentValue"])
			}
		}
	}
	if tmpCheckType != "" && tmpThreshold != "" && tmpErrorCode != "" && tmpAbnormalInformation != ""{
		bc = []string{tmpCheckType, tmpThreshold, tmpErrorCode, tmpAbnormalInformation}
	}
	return bc
}

func (out *OutputWayStruct) tmpResultSummary(CheckTypeSliect []string) [][]string{
	var bc [][]string
	var dd []map[string]string
	var cc []string
	var tmpEQ int
	for i := range CheckTypeSliect{
		switch CheckTypeSliect[i] {
		case "configParameter":
			dd = pub.InspectionResult.DatabaseConfigCheck.ConfigParameter
		case "binlogDiskUsageRate":
			dd = pub.InspectionResult.DatabasePerformance.PerformanceStatus.BinlogDiskUsageRate
		case "historyConnectionMaxUsageRat":
			dd = pub.InspectionResult.DatabasePerformance.PerformanceStatus.HistoryConnectionMaxUsageRate
		case "tmpDiskTableUsageRate":
			dd = pub.InspectionResult.DatabasePerformance.PerformanceStatus.TmpDiskTableUsageRate
		case "tmpDiskfileUsageRate":
			dd = pub.InspectionResult.DatabasePerformance.PerformanceStatus.TmpDiskfileUsageRate
		case "innodbBufferPoolUsageRate":
			dd = pub.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolUsageRate
		case "innodbBufferPoolDirtyPagesRate":
			dd = pub.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolDirtyPagesRate
		case "innodbBufferPoolHitRate":
			dd = pub.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolHitRate
		case "openFileUsageRate":
			dd = pub.InspectionResult.DatabasePerformance.PerformanceStatus.OpenFileUsageRate
		case "openTableCacheUsageRate":
			dd = pub.InspectionResult.DatabasePerformance.PerformanceStatus.OpenTableCacheUsageRate
		case "openTableCacheOverflowsUsageRate":
			dd = pub.InspectionResult.DatabasePerformance.PerformanceStatus.OpenTableCacheOverflowsUsageRate
		case "selectScanUsageRate":
			dd = pub.InspectionResult.DatabasePerformance.PerformanceStatus.SelectScanUsageRate
		case "selectfullJoinScanUsageRate":
			dd = pub.InspectionResult.DatabasePerformance.PerformanceStatus.SelectfullJoinScanUsageRate
		case "tableAutoPrimaryKeyUsageRate":
			dd = pub.InspectionResult.DatabasePerformance.PerformanceTableIndex.TableAutoPrimaryKeyUsageRate
		case "tableRows":
			dd = pub.InspectionResult.DatabasePerformance.PerformanceTableIndex.TableRows
		case "diskFragmentationRate":
			dd = pub.InspectionResult.DatabasePerformance.PerformanceTableIndex.DiskFragmentationRate
		case "bigTable":
			dd = pub.InspectionResult.DatabasePerformance.PerformanceTableIndex.BigTable
		case "coldTable":
			dd = pub.InspectionResult.DatabasePerformance.PerformanceTableIndex.ColdTable
		case "tableCharset":
			dd = pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableCharset
		case "tableEngine":
			dd = pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableEngine
		case "tableForeign":
			dd = pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableForeign
		case "tableNoPrimaryKey":
			dd = pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableNoPrimaryKey
		case "tableAutoIncrement":
			dd = pub.InspectionResult.DatabaseBaselineCheck.ColumnDesign.TableAutoIncrement
		case "tableBigColumns":
			dd = pub.InspectionResult.DatabaseBaselineCheck.ColumnDesign.TableBigColumns
		case "indexColumnIsNull":
			dd = pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnIsNull
		case "indexColumnIsEnumSet":
			dd = pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnIsEnumSet
		case "indexColumnIsBlobText":
			dd = pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnIsBlobText
		case "tableIncludeRepeatIndex":
			dd = pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnIsRepeatIndex
		case "tableProcedure":
			dd = pub.InspectionResult.DatabaseBaselineCheck.ProcedureTriggerDesign.TableProcedure
		case "tableFunc":
			dd = pub.InspectionResult.DatabaseBaselineCheck.ProcedureTriggerDesign.TableFunc
		case "tableTrigger":
			dd = pub.InspectionResult.DatabaseBaselineCheck.ProcedureTriggerDesign.TableTrigger
		case "anonymousUsers":
			dd = pub.InspectionResult.DatabaseSecurity.UserPriDesign.AnonymousUsers
		case "emptyPasswordUser":
			dd = pub.InspectionResult.DatabaseSecurity.UserPriDesign.EmptyPasswordUser
		case "rootUserRemoteLogin":
			dd = pub.InspectionResult.DatabaseSecurity.UserPriDesign.RootUserRemoteLogin
		case "normalUserConnectionUnlimited":
			dd = pub.InspectionResult.DatabaseSecurity.UserPriDesign.NormalUserConnectionUnlimited
		case "userPasswordSame":
			dd = pub.InspectionResult.DatabaseSecurity.UserPriDesign.UserPasswordSame
		case "normalUserDatabaseAllPrivilages":
			dd = pub.InspectionResult.DatabaseSecurity.UserPriDesign.NormalUserDatabaseAllPrivilages
		case "normalUserSuperPrivilages":
			dd = pub.InspectionResult.DatabaseSecurity.UserPriDesign.NormalUserSuperPrivilages
		case "databasePort":
			dd = pub.InspectionResult.DatabaseSecurity.PortDesign.DatabasePort
		}
		cc = out.tmpcc(dd)
		if cc != nil{
			tmpEQ ++
			cd := append([]string{strconv.Itoa(tmpEQ)},cc...)
			bc = append(bc,cd)
		}
	}
	return bc
}

func (out *OutputWayStruct) ResultSummaryStringSlice()  [][]string{
	var resultProfile [][]string
	aa := out.tmpaa("01","configParameter",pub.InspectionResult.DatabaseConfigCheck.ConfigParameter)
	resultProfile = append(resultProfile,aa)
	ba := out.tmpaa("02","tableCharset",pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableCharset)
	resultProfile = append(resultProfile,ba)
	ca := out.tmpaa("03","tableEngine",pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableEngine)
	resultProfile = append(resultProfile,ca)
	da := out.tmpaa("04","tableNoPrimaryKey",pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableNoPrimaryKey)
	resultProfile = append(resultProfile,da)
	ea := out.tmpaa("05","tableForeign",pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableForeign)
	resultProfile = append(resultProfile,ea)
	fa := out.tmpaa("06","tableAutoIncrement",pub.InspectionResult.DatabaseBaselineCheck.ColumnDesign.TableAutoIncrement)
	resultProfile = append(resultProfile,fa)
	ga := out.tmpaa("07","tableBigColumns",pub.InspectionResult.DatabaseBaselineCheck.ColumnDesign.TableBigColumns)
	resultProfile = append(resultProfile,ga)
	ha := out.tmpaa("08","indexColumnIsNull",pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnIsNull)
	resultProfile = append(resultProfile,ha)
	ia := out.tmpaa("09","indexColumnIsEnumSet",pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnIsEnumSet)
	resultProfile = append(resultProfile,ia)
	ja := out.tmpaa("10","indexColumnIsBlobText",pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnIsBlobText)
	resultProfile = append(resultProfile,ja)
	ka := out.tmpaa("11","tableIncludeRepeatIndex",pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnIsRepeatIndex)
	resultProfile = append(resultProfile,ka)
	la := out.tmpaa("12","tableProcedure",pub.InspectionResult.DatabaseBaselineCheck.ProcedureTriggerDesign.TableProcedure)
	resultProfile = append(resultProfile,la)
	ma := out.tmpaa("13","tableFunc",pub.InspectionResult.DatabaseBaselineCheck.ProcedureTriggerDesign.TableFunc)
	resultProfile = append(resultProfile,ma)
	na := out.tmpaa("14","tableTrigger",pub.InspectionResult.DatabaseBaselineCheck.ProcedureTriggerDesign.TableTrigger)
	resultProfile = append(resultProfile,na)
	oa := out.tmpaa("15","anonymousUsers",pub.InspectionResult.DatabaseSecurity.UserPriDesign.AnonymousUsers)
	resultProfile = append(resultProfile,oa)
	pa := out.tmpaa("16","emptyPasswordUser",pub.InspectionResult.DatabaseSecurity.UserPriDesign.EmptyPasswordUser)
	resultProfile = append(resultProfile,pa)
	qa := out.tmpaa("17","rootUserRemoteLogin",pub.InspectionResult.DatabaseSecurity.UserPriDesign.RootUserRemoteLogin)
	resultProfile = append(resultProfile,qa)
	ra := out.tmpaa("18","normalUserConnectionUnlimited",pub.InspectionResult.DatabaseSecurity.UserPriDesign.NormalUserConnectionUnlimited)
	resultProfile = append(resultProfile,ra)
	ta := out.tmpaa("19","userPasswordSame",pub.InspectionResult.DatabaseSecurity.UserPriDesign.UserPasswordSame)
	resultProfile = append(resultProfile,ta)
	ua := out.tmpaa("20","normalUserDatabaseAllPrivilages",pub.InspectionResult.DatabaseSecurity.UserPriDesign.NormalUserDatabaseAllPrivilages)
	resultProfile = append(resultProfile,ua)
	va := out.tmpaa("21","normalUserSuperPrivilages",pub.InspectionResult.DatabaseSecurity.UserPriDesign.NormalUserSuperPrivilages)
	resultProfile = append(resultProfile,va)
	wa := out.tmpaa("22","databasePort",pub.InspectionResult.DatabaseSecurity.PortDesign.DatabasePort)
	resultProfile = append(resultProfile,wa)
	xa := out.tmpaa("23","binlogDiskUsageRate",pub.InspectionResult.DatabasePerformance.PerformanceStatus.BinlogDiskUsageRate)
	resultProfile = append(resultProfile,xa)
	ya := out.tmpaa("24","historyConnectionMaxUsageRate",pub.InspectionResult.DatabasePerformance.PerformanceStatus.HistoryConnectionMaxUsageRate)
	resultProfile = append(resultProfile,ya)
	za := out.tmpaa("25","tmpDiskTableUsageRate",pub.InspectionResult.DatabasePerformance.PerformanceStatus.TmpDiskTableUsageRate)
	resultProfile = append(resultProfile,za)
	ab := out.tmpaa("26","tmpDiskfileUsageRate",pub.InspectionResult.DatabasePerformance.PerformanceStatus.TmpDiskfileUsageRate)
	resultProfile = append(resultProfile,ab)
	bb := out.tmpaa("27","innodbBufferPoolUsageRate",pub.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolUsageRate)
	resultProfile = append(resultProfile,bb)
	cb := out.tmpaa("28","innodbBufferPoolDirtyPagesRate",pub.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolDirtyPagesRate)
	resultProfile = append(resultProfile,cb)
	db := out.tmpaa("29","innodbBufferPoolHitRate",pub.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolHitRate)
	resultProfile = append(resultProfile,db)
	eb := out.tmpaa("30","openFileUsageRate",pub.InspectionResult.DatabasePerformance.PerformanceStatus.OpenFileUsageRate)
	resultProfile = append(resultProfile,eb)
	fb := out.tmpaa("31","openTableCacheUsageRate",pub.InspectionResult.DatabasePerformance.PerformanceStatus.OpenTableCacheUsageRate)
	resultProfile = append(resultProfile,fb)
	gb := out.tmpaa("32","openTableCacheOverflowsUsageRate",pub.InspectionResult.DatabasePerformance.PerformanceStatus.OpenTableCacheOverflowsUsageRate)
	resultProfile = append(resultProfile,gb)
	hb := out.tmpaa("33","selectScanUsageRate",pub.InspectionResult.DatabasePerformance.PerformanceStatus.SelectScanUsageRate)
	resultProfile = append(resultProfile,hb)
	ib := out.tmpaa("34","selectfullJoinScanUsageRate",pub.InspectionResult.DatabasePerformance.PerformanceStatus.SelectfullJoinScanUsageRate)
	resultProfile = append(resultProfile,ib)
	jb := out.tmpaa("35","tableAutoPrimaryKeyUsageRate",pub.InspectionResult.DatabasePerformance.PerformanceTableIndex.TableAutoPrimaryKeyUsageRate)
	resultProfile = append(resultProfile,jb)
	kb := out.tmpaa("36","tableRows",pub.InspectionResult.DatabasePerformance.PerformanceTableIndex.TableRows)
	resultProfile = append(resultProfile,kb)
	lb := out.tmpaa("37","diskFragmentationRate",pub.InspectionResult.DatabasePerformance.PerformanceTableIndex.DiskFragmentationRate)
	resultProfile = append(resultProfile,lb)
	mb := out.tmpaa("38","bigTable",pub.InspectionResult.DatabasePerformance.PerformanceTableIndex.BigTable)
	resultProfile = append(resultProfile,mb)
	nb := out.tmpaa("39","coldTable",pub.InspectionResult.DatabasePerformance.PerformanceTableIndex.ColdTable)
	resultProfile = append(resultProfile,nb)
	return resultProfile
}
