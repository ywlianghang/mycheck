package InspectionItem

import (
	"DepthInspection/api/PublicClass"
	"fmt"
	"strings"
)

func (baselineCheck *DatabaseBaselineCheckStruct) BaselineCheckUserPriDesign() {
	PublicClass.Loggs.Info("Begin to check for any vulnerabilities in database user privileges")
	var tmpPassword,tmpUser,tmpHost interface{}
	for i := range PublicClass.MysqlUser{
		cc := PublicClass.MysqlUser[i]
		var d = make(map[string]string)
		d["user"] = cc["user"].(string)
		d["host"] = cc["host"].(string)
		//检查匿名用户
		if cc["user"] == "" {
			d["checkStatus"] = "abnormal"    //异常
			d["checkType"] = "anonymousUsers"
			PublicClass.InspectionResult.BaselineCheckUserPriDesign.AnonymousUsers = append(PublicClass.InspectionResult.BaselineCheckUserPriDesign.AnonymousUsers,d)
			PublicClass.Loggs.Error(fmt.Sprintf("Anonymous users currently exist. The information is as follows: user: \"%s\" host: \"%s\"",cc["user"],cc["host"]))
		}else{
			d["checkStatus"] = "normal"    //正常
			d["checkType"] = "anonymousUsers"
			PublicClass.InspectionResult.BaselineCheckUserPriDesign.AnonymousUsers = append(PublicClass.InspectionResult.BaselineCheckUserPriDesign.AnonymousUsers,d)
		}
		m := newMap(d)
		//检查空密码用户
		if cc["password"] == ""{
			m["checkStatus"] = "abnormal"    //异常
			m["checkType"] = "emptyPasswordUser"
			PublicClass.InspectionResult.BaselineCheckUserPriDesign.EmptyPasswordUser = append(PublicClass.InspectionResult.BaselineCheckUserPriDesign.EmptyPasswordUser,m)
			PublicClass.Loggs.Error(fmt.Sprintf("The current username password is empty. The information is as follows: user: \"%s\" host: \"%s\"",cc["user"],cc["host"]))
		}else{
			m["checkStatus"] = "normal"    //异常
			m["checkType"] = "emptyPasswordUser"
			PublicClass.InspectionResult.BaselineCheckUserPriDesign.EmptyPasswordUser = append(PublicClass.InspectionResult.BaselineCheckUserPriDesign.EmptyPasswordUser,m)
		}
		n := newMap(d)
		//检查root用户远端登录，只能本地连接
		if cc["user"] == "root" && cc["host"] != "localhost" && cc["host"] != "127.0.0.1"{
			n["checkStatus"] = "abnormal"    //异常
			n["checkType"] = "rootUserRemoteLogin"
			PublicClass.InspectionResult.BaselineCheckUserPriDesign.RootUserRemoteLogin = append(PublicClass.InspectionResult.BaselineCheckUserPriDesign.RootUserRemoteLogin,n)
			PublicClass.Loggs.Error(fmt.Sprintf("The root user is currently in remote login danger. The information is as follows: user: \"%s\" host: \"%s\"",cc["user"],cc["host"]))
		}else{
			n["checkStatus"] = "normal"    //异常
			n["checkType"] = "rootUserRemoteLogin"
			PublicClass.InspectionResult.BaselineCheckUserPriDesign.RootUserRemoteLogin = append(PublicClass.InspectionResult.BaselineCheckUserPriDesign.RootUserRemoteLogin,m)
		}
		o := newMap(d)
		//检查普通用户远端连接的限制，不允许使用%
		if cc["user"] != "" && cc["user"] != "root" && cc["password"] != "" && cc["host"] == "%"{
			o["checkStatus"] = "abnormal"    //异常
			o["checkType"] = "normalUserConnectionUnlimited"
			PublicClass.InspectionResult.BaselineCheckUserPriDesign.NormalUserConnectionUnlimited = append(PublicClass.InspectionResult.BaselineCheckUserPriDesign.NormalUserConnectionUnlimited,o)
			PublicClass.Loggs.Error(fmt.Sprintf("The current user name has no connection IP limit. The information is as follows: user: \"%s\" host: \"%s\"",cc["user"],cc["host"]))
		}else{
			o["checkStatus"] = "normal"    //异常
			o["checkType"] = "normalUserConnectionUnlimited"
			PublicClass.InspectionResult.BaselineCheckUserPriDesign.NormalUserConnectionUnlimited = append(PublicClass.InspectionResult.BaselineCheckUserPriDesign.NormalUserConnectionUnlimited,o)
		}

		//检查不同用户使用相同密码
		p := newMap(d)
		if cc["password"] == tmpPassword{
			p["checkStatus"] = "abnormal"    //异常
			p["checkType"] = "userPasswordSame"
			PublicClass.InspectionResult.BaselineCheckUserPriDesign.UserPasswordSame = append(PublicClass.InspectionResult.BaselineCheckUserPriDesign.UserPasswordSame,p)
			PublicClass.Loggs.Error(fmt.Sprintf("Different users in the current database use the same password, please change it. The information is as follows: user1: \"%v@%v\"  user2: \"%s@%s\"",tmpUser,tmpHost,cc["user"],cc["host"]))
		}else{
			p["checkStatus"] = "normal"    //异常
			p["checkType"] = "userPasswordSame"
			PublicClass.InspectionResult.BaselineCheckUserPriDesign.UserPasswordSame = append(PublicClass.InspectionResult.BaselineCheckUserPriDesign.UserPasswordSame,p)
		}
		tmpPassword = cc["password"]
		tmpUser = cc["user"]
		tmpHost = cc["host"]

		//检查跨用户权限*.*
		strSql := fmt.Sprintf("show grants for '%s'@'%s'",cc["user"],cc["host"])
		cd := PublicClass.DBexecInter.DBQueryDateString(strSql)
		if cc["user"] != "root" && cc["host"] != "localhost" && cc["host"] != "127.0.0.1"{
			//检查当前用户是否存在ON *.*
			q := newMap(d)
			if strings.Contains(cd,"ON *.*") {
				q["checkStatus"] = "abnormal"    //异常
				q["checkType"] = "normalUserDatabaseAllPrivilages"
				PublicClass.InspectionResult.BaselineCheckUserPriDesign.NormalUserDatabaseAllPrivilages = append(PublicClass.InspectionResult.BaselineCheckUserPriDesign.NormalUserDatabaseAllPrivilages,q)
				PublicClass.Loggs.Error(fmt.Sprintf("Cross-user permissions currently exist (ON *.*). The information is as follows: user@host: \"%s@%s\"",cc["user"],cc["host"]))
			} else{
				q["checkStatus"] = "abnormal"    //异常
				q["checkType"] = "normalUserDatabaseAllPrivilages"
				PublicClass.InspectionResult.BaselineCheckUserPriDesign.NormalUserDatabaseAllPrivilages = append(PublicClass.InspectionResult.BaselineCheckUserPriDesign.NormalUserDatabaseAllPrivilages,q)
			}
			r := newMap(d)
			//检查当前用户是否WITH GRANT OPTION
			if strings.Contains(cd,"WITH GRANT OPTION"){
				r["checkStatus"] = "abnormal"    //异常
				r["checkType"] = "normalUserSuperPrivilages"
				PublicClass.InspectionResult.BaselineCheckUserPriDesign.NormalUserSuperPrivilages = append(PublicClass.InspectionResult.BaselineCheckUserPriDesign.NormalUserSuperPrivilages,r)
				PublicClass.Loggs.Error(fmt.Sprintf("The current user has permission transfer (WITH GRANT OPTION). The information is as follows: user@host: \"%s@%s\"",cc["user"],cc["host"]))
			}else{
				r["checkStatus"] = "normal"    //正常
				r["checkType"] = "normalUserSuperPrivilages"
				PublicClass.InspectionResult.BaselineCheckUserPriDesign.NormalUserSuperPrivilages = append(PublicClass.InspectionResult.BaselineCheckUserPriDesign.NormalUserSuperPrivilages,r)
			}
		}
	}
}

func (baselineCheck *DatabaseBaselineCheckStruct) BaselineCheckPortDesign() {
	PublicClass.Loggs.Info("Begin checking whether the current database is using the default port 3306")
	var d = make(map[string]string)
	cc := PublicClass.GlobalVariables
	d["checkStatus"] = "normal"    //正常
	d["checkType"] = "nil"
	if cc["port"] == "3306"{
		d["checkStatus"] = "abnormal"    //异常
		d["checkType"] = "databasePort"
		PublicClass.Loggs.Error(fmt.Sprintf("The MySQL service uses the default port. The information is as follows: using port: %s.",cc["port"]))
	}
	PublicClass.InspectionResult.BaselineCheckPortDesign.DatabasePort = append(PublicClass.InspectionResult.BaselineCheckPortDesign.DatabasePort,d)
}