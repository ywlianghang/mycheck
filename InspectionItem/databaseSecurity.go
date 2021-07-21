package InspectionItem

import (
	"DepthInspection/api/PublicDB"
	"fmt"
	"strings"
)

func (baselineCheck *DatabaseBaselineCheckStruct) BaselineCheckUserPriDesign(aa *PublicDB.ConfigInfo) {
	aa.Loggs.Info("Begin to check for any vulnerabilities in database user privileges")
	strSql := fmt.Sprintf("select user,host,authentication_string password from mysql.user")
	cc := aa.DatabaseExecInterf.DBQueryDateJson(aa,strSql)
	var tmpPassword,tmpUser,tmpHost interface{}
	for i := range cc{
		//检查匿名用户
		if cc[i]["user"] == "" {
			aa.Loggs.Error(fmt.Sprintf("Anonymous users currently exist. The information is as follows: user: \"%s\" host: \"%s\"",cc[i]["user"],cc[i]["host"]))
		}
		//检查空密码用户
		if cc[i]["password"] == ""{
			aa.Loggs.Error(fmt.Sprintf("The current username password is empty. The information is as follows: user: \"%s\" host: \"%s\"",cc[i]["user"],cc[i]["host"]))
		}
		//检查root用户远端登录，只能本地连接
		if cc[i]["user"] == "root" && cc[i]["host"] != "localhost" && cc[i]["host"] != "127.0.0.1"{
			aa.Loggs.Error(fmt.Sprintf("The root user is currently in remote login danger. The information is as follows: user: \"%s\" host: \"%s\"",cc[i]["user"],cc[i]["host"]))
		}
		//检查普通用户远端连接的限制，不允许使用%
		if cc[i]["user"] != "" && cc[i]["user"] != "root" && cc[i]["password"] != "" && cc[i]["host"] == "%"{
			aa.Loggs.Error(fmt.Sprintf("The current user name has no connection IP limit. The information is as follows: user: \"%s\" host: \"%s\"",cc[i]["user"],cc[i]["host"]))
		}
		//检查不同用户使用相同密码
		if cc[i]["password"] == tmpPassword{
			aa.Loggs.Error(fmt.Sprintf("Different users in the current database use the same password, please change it. The information is as follows: user1: \"%v@%v\"  user2: \"%s@%s\"",tmpUser,tmpHost,cc[i]["user"],cc[i]["host"]))
		}
		tmpUser = cc[i]["user"]
		tmpHost = cc[i]["host"]
		tmpPassword = cc[i]["password"]
		//检查跨用户权限*.*
		strSql = fmt.Sprintf("show grants for '%s'@'%s'",cc[i]["user"],cc[i]["host"])
		cd := aa.DatabaseExecInterf.DBQueryDateString(aa,strSql)
		if cc[i]["user"] != "root" && cc[i]["host"] != "localhost" && cc[i]["host"] != "127.0.0.1"{
			//检查当前用户是否存在ON *.*
			if strings.Contains(cd,"ON *.*") {
				aa.Loggs.Error(fmt.Sprintf("Cross-user permissions currently exist (ON *.*). The information is as follows: user@host: \"%s@%s\"",cc[i]["user"],cc[i]["host"]))
				//检查当前用户是否WITH GRANT OPTION
			}else if strings.Contains(cd,"WITH GRANT OPTION"){
				aa.Loggs.Error(fmt.Sprintf("The current user has permission transfer (WITH GRANT OPTION). The information is as follows: user@host: \"%s@%s\"",cc[i]["user"],cc[i]["host"]))
			}
		}
	}
}

func (baselineCheck *DatabaseBaselineCheckStruct) BaselineCheckPortDesign(aa *PublicDB.ConfigInfo) {
	aa.Loggs.Info("Begin checking whether the current database is using the default port 3306")
	strSql := fmt.Sprintf("show global variables like 'port'")
	cc := aa.DatabaseExecInterf.DBQueryDateMap(aa,strSql)
	if cc["port"] == "3306"{
		aa.Loggs.Error(fmt.Sprintf("The MySQL service uses the default port. The information is as follows: using port: %s.",cc["port"]))
	}
}