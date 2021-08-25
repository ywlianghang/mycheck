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
var abcd = make(map[string][]map[string]string)
func tmpInit() {
	abcd["configParameter"] = pub.InspectionResult.DatabaseConfigCheck.ConfigParameter
	abcd["tableCharset"] = pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableCharset
	abcd["tableEngine"] = pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableEngine
	abcd["tableNoPrimaryKey"] = pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableNoPrimaryKey
	abcd["tableForeign"] = pub.InspectionResult.DatabaseBaselineCheck.TableDesign.TableForeign
	abcd["tableAutoIncrement"] = pub.InspectionResult.DatabaseBaselineCheck.ColumnDesign.TableAutoIncrement
	abcd["tableBigColumns"] = pub.InspectionResult.DatabaseBaselineCheck.ColumnDesign.TableBigColumns
	abcd["indexColumnIsNull"] = pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnIsNull
	abcd["indexColumnType"] = pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnType
	abcd["tableIncludeRepeatIndex"] = pub.InspectionResult.DatabaseBaselineCheck.IndexColumnsDesign.IndexColumnIsRepeatIndex
	abcd["tableProcedureFunc"] = pub.InspectionResult.DatabaseBaselineCheck.ProcedureTriggerDesign.TableProcedure
	abcd["tableTrigger"] = pub.InspectionResult.DatabaseBaselineCheck.ProcedureTriggerDesign.TableTrigger
	abcd["anonymousUsers"] = pub.InspectionResult.DatabaseSecurity.UserPriDesign.AnonymousUsers
	abcd["emptyPasswordUser"] = pub.InspectionResult.DatabaseSecurity.UserPriDesign.EmptyPasswordUser
	abcd["rootUserRemoteLogin"] = pub.InspectionResult.DatabaseSecurity.UserPriDesign.RootUserRemoteLogin
	abcd["normalUserConnectionUnlimited"] = pub.InspectionResult.DatabaseSecurity.UserPriDesign.NormalUserConnectionUnlimited
	abcd["userPasswordSame"] = pub.InspectionResult.DatabaseSecurity.UserPriDesign.UserPasswordSame
	abcd["normalUserDatabaseAllPrivilages"] = pub.InspectionResult.DatabaseSecurity.UserPriDesign.NormalUserDatabaseAllPrivilages
	abcd["normalUserSuperPrivilages"] = pub.InspectionResult.DatabaseSecurity.UserPriDesign.NormalUserSuperPrivilages
	abcd["databasePort"] = pub.InspectionResult.DatabaseSecurity.PortDesign.DatabasePort
	abcd["binlogDiskUsageRate"] = pub.InspectionResult.DatabasePerformance.PerformanceStatus.BinlogDiskUsageRate
	abcd["historyConnectionMaxUsageRate"] = pub.InspectionResult.DatabasePerformance.PerformanceStatus.HistoryConnectionMaxUsageRate
	abcd["tmpDiskTableUsageRate"] = pub.InspectionResult.DatabasePerformance.PerformanceStatus.TmpDiskTableUsageRate
	abcd["tmpDiskfileUsageRate"] = pub.InspectionResult.DatabasePerformance.PerformanceStatus.TmpDiskfileUsageRate
	abcd["innodbBufferPoolUsageRate"] = pub.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolUsageRate
	abcd["innodbBufferPoolDirtyPagesRate"] = pub.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolDirtyPagesRate
	abcd["innodbBufferPoolHitRate"] = pub.InspectionResult.DatabasePerformance.PerformanceStatus.InnodbBufferPoolHitRate
	abcd["openFileUsageRate"] = pub.InspectionResult.DatabasePerformance.PerformanceStatus.OpenFileUsageRate
	abcd["openTableCacheUsageRate"] = pub.InspectionResult.DatabasePerformance.PerformanceStatus.OpenTableCacheUsageRate
	abcd["openTableCacheOverflowsUsageRate"] = pub.InspectionResult.DatabasePerformance.PerformanceStatus.OpenTableCacheOverflowsUsageRate
	abcd["selectScanUsageRate"] = pub.InspectionResult.DatabasePerformance.PerformanceStatus.SelectScanUsageRate
	abcd["selectfullJoinScanUsageRate"] = pub.InspectionResult.DatabasePerformance.PerformanceStatus.SelectfullJoinScanUsageRate
	abcd["tableAutoPrimaryKeyUsageRate"] = pub.InspectionResult.DatabasePerformance.PerformanceTableIndex.TableAutoPrimaryKeyUsageRate
	abcd["tableRows"] = pub.InspectionResult.DatabasePerformance.PerformanceTableIndex.TableRows
	abcd["diskFragmentationRate"] = pub.InspectionResult.DatabasePerformance.PerformanceTableIndex.DiskFragmentationRate
	abcd["bigTable"] = pub.InspectionResult.DatabasePerformance.PerformanceTableIndex.BigTable
	abcd["coldTable"] = pub.InspectionResult.DatabasePerformance.PerformanceTableIndex.ColdTable
}
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
	//var dd []map[string]string
	var cc []string
	tmpInit()
	var tmpeqInt int
	var tmpeqStr string
	for i := range CheckTypeSliect{
		if vi,ok := abcd[CheckTypeSliect[i]];ok{
			if vi != nil{
				tmpeqInt ++
				if tmpeqInt <10{
					tmpeqStr = fmt.Sprintf("0%d",tmpeqInt)
				}else{
					tmpeqStr = fmt.Sprintf("%d",tmpeqInt)
				}
				cc = out.tmpcc(vi)
				if cc != nil{
					cd := append([]string{tmpeqStr},cc...)
					bc = append(bc,cd)
				}
			}
		}
	}
	return bc
}

func (out *OutputWayStruct) ResultSummaryStringSlice()  [][]string{
	var resultProfile [][]string
	tmpInit()
	var tmpeqInt int
	var tmpeqStr string
	for k,v := range abcd {
		if v != nil{
			tmpeqInt ++
			if tmpeqInt <10{
				tmpeqStr = fmt.Sprintf("0%d",tmpeqInt)
			}else{
				tmpeqStr = fmt.Sprintf("%d",tmpeqInt)
			}
			tmpRes := out.tmpaa(tmpeqStr,k,v)
			resultProfile = append(resultProfile,tmpRes)
		}
	}
	return resultProfile
}
