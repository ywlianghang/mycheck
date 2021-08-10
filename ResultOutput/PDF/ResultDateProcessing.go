package PDF

import (
	"DepthInspection/api/PublicClass"
	"strconv"
)

func (out *OutputWayStruct) tmpaa(checkNum,checkType string,ast []map[string]string) []string{
	var abnormalCount int = 0
	var normalCount int = 0
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
func (out *OutputWayStruct) ResultSummaryStringSlice()  [][]string{
	var resultProfile [][]string
	aa := out.tmpaa("01","configParameter",PublicClass.InspectionResult.DatabaseConfigCheck.ConfigParameter)
	resultProfile = append(resultProfile,aa)
	ba := out.tmpaa("02","tableCharset",PublicClass.InspectionResult.BaselineCheckTablesDesign.TableCharset)
	resultProfile = append(resultProfile,ba)
	ca := out.tmpaa("03","tableEngine",PublicClass.InspectionResult.BaselineCheckTablesDesign.TableEngine)
	resultProfile = append(resultProfile,ca)
	da := out.tmpaa("04","tableNoPrimaryKey",PublicClass.InspectionResult.BaselineCheckTablesDesign.TableNoPrimaryKey)
	resultProfile = append(resultProfile,da)
	ea := out.tmpaa("05","tableForeign",PublicClass.InspectionResult.BaselineCheckTablesDesign.TableForeign)
	resultProfile = append(resultProfile,ea)
	fa := out.tmpaa("06","tableAutoIncrement",PublicClass.InspectionResult.BaselineCheckColumnsDesign.TableAutoIncrement)
	resultProfile = append(resultProfile,fa)
	ga := out.tmpaa("07","tableBigColumns",PublicClass.InspectionResult.BaselineCheckColumnsDesign.TableBigColumns)
	resultProfile = append(resultProfile,ga)
	ha := out.tmpaa("08","indexColumnIsNull",PublicClass.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsNull)
	resultProfile = append(resultProfile,ha)
	ia := out.tmpaa("09","indexColumnIsEnumSet",PublicClass.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsEnumSet)
	resultProfile = append(resultProfile,ia)
	ja := out.tmpaa("10","indexColumnIsBlobText",PublicClass.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsBlobText)
	resultProfile = append(resultProfile,ja)
	ka := out.tmpaa("11","tableIncludeRepeatIndex",PublicClass.InspectionResult.BaselineCheckIndexColumnDesign.IndexColumnIsRepeatIndex)
	resultProfile = append(resultProfile,ka)
	la := out.tmpaa("12","tableProcedure",PublicClass.InspectionResult.BaselineCheckProcedureTriggerDesign.TableProcedure)
	resultProfile = append(resultProfile,la)
	ma := out.tmpaa("13","tableFunc",PublicClass.InspectionResult.BaselineCheckProcedureTriggerDesign.TableFunc)
	resultProfile = append(resultProfile,ma)
	na := out.tmpaa("14","tableTrigger",PublicClass.InspectionResult.BaselineCheckProcedureTriggerDesign.TableTrigger)
	resultProfile = append(resultProfile,na)
	oa := out.tmpaa("15","anonymousUsers",PublicClass.InspectionResult.BaselineCheckUserPriDesign.AnonymousUsers)
	resultProfile = append(resultProfile,oa)
	pa := out.tmpaa("16","emptyPasswordUser",PublicClass.InspectionResult.BaselineCheckUserPriDesign.EmptyPasswordUser)
	resultProfile = append(resultProfile,pa)
	qa := out.tmpaa("17","rootUserRemoteLogin",PublicClass.InspectionResult.BaselineCheckUserPriDesign.RootUserRemoteLogin)
	resultProfile = append(resultProfile,qa)
	ra := out.tmpaa("18","normalUserConnectionUnlimited",PublicClass.InspectionResult.BaselineCheckUserPriDesign.NormalUserConnectionUnlimited)
	resultProfile = append(resultProfile,ra)
	ta := out.tmpaa("19","userPasswordSame",PublicClass.InspectionResult.BaselineCheckUserPriDesign.UserPasswordSame)
	resultProfile = append(resultProfile,ta)
	ua := out.tmpaa("20","normalUserDatabaseAllPrivilages",PublicClass.InspectionResult.BaselineCheckUserPriDesign.NormalUserDatabaseAllPrivilages)
	resultProfile = append(resultProfile,ua)
	va := out.tmpaa("21","normalUserSuperPrivilages",PublicClass.InspectionResult.BaselineCheckUserPriDesign.NormalUserSuperPrivilages)
	resultProfile = append(resultProfile,va)
	wa := out.tmpaa("22","databasePort",PublicClass.InspectionResult.BaselineCheckPortDesign.DatabasePort)
	resultProfile = append(resultProfile,wa)
	xa := out.tmpaa("23","binlogDiskUsageRate",PublicClass.InspectionResult.DatabasePerformanceCheck.BinlogDiskUsageRate)
	resultProfile = append(resultProfile,xa)
	ya := out.tmpaa("24","historyConnectionMaxUsageRate",PublicClass.InspectionResult.DatabasePerformanceCheck.HistoryConnectionMaxUsageRate)
	resultProfile = append(resultProfile,ya)
	za := out.tmpaa("25","tmpDiskTableUsageRate",PublicClass.InspectionResult.DatabasePerformanceCheck.TmpDiskTableUsageRate)
	resultProfile = append(resultProfile,za)
	ab := out.tmpaa("26","tmpDiskfileUsageRate",PublicClass.InspectionResult.DatabasePerformanceCheck.TmpDiskfileUsageRate)
	resultProfile = append(resultProfile,ab)
	bb := out.tmpaa("27","innodbBufferPoolUsageRate",PublicClass.InspectionResult.DatabasePerformanceCheck.InnodbBufferPoolUsageRate)
	resultProfile = append(resultProfile,bb)
	cb := out.tmpaa("28","innodbBufferPoolDirtyPagesRate",PublicClass.InspectionResult.DatabasePerformanceCheck.InnodbBufferPoolDirtyPagesRate)
	resultProfile = append(resultProfile,cb)
	db := out.tmpaa("29","innodbBufferPoolHitRate",PublicClass.InspectionResult.DatabasePerformanceCheck.InnodbBufferPoolHitRate)
	resultProfile = append(resultProfile,db)
	eb := out.tmpaa("30","openFileUsageRate",PublicClass.InspectionResult.DatabasePerformanceCheck.OpenFileUsageRate)
	resultProfile = append(resultProfile,eb)
	fb := out.tmpaa("31","openTableCacheUsageRate",PublicClass.InspectionResult.DatabasePerformanceCheck.OpenTableCacheUsageRate)
	resultProfile = append(resultProfile,fb)
	gb := out.tmpaa("32","openTableCacheOverflowsUsageRate",PublicClass.InspectionResult.DatabasePerformanceCheck.OpenTableCacheOverflowsUsageRate)
	resultProfile = append(resultProfile,gb)
	hb := out.tmpaa("33","selectScanUsageRate",PublicClass.InspectionResult.DatabasePerformanceCheck.SelectScanUsageRate)
	resultProfile = append(resultProfile,hb)
	ib := out.tmpaa("34","selectfullJoinScanUsageRate",PublicClass.InspectionResult.DatabasePerformanceCheck.SelectfullJoinScanUsageRate)
	resultProfile = append(resultProfile,ib)
	jb := out.tmpaa("35","tableAutoPrimaryKeyUsageRate",PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.TableAutoPrimaryKeyUsageRate)
	resultProfile = append(resultProfile,jb)
	kb := out.tmpaa("36","tableRows",PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.TableRows)
	resultProfile = append(resultProfile,kb)
	lb := out.tmpaa("37","diskFragmentationRate",PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.DiskFragmentationRate)
	resultProfile = append(resultProfile,lb)
	mb := out.tmpaa("38","bigTable",PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.BigTable)
	resultProfile = append(resultProfile,mb)
	nb := out.tmpaa("39","coldTable",PublicClass.InspectionResult.DatabasePerformanceTableIndexCheck.ColdTable)
	resultProfile = append(resultProfile,nb)
	return resultProfile
}
