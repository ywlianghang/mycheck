package flag

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"runtime"
	"strings"
)

type parameter struct {
	Config string
	helpbool bool
}
var CheckParameter = &parameter{}
func cliHelp(){
	app := cli.NewApp()
	app.Name = "mycheck"                         //应用名称
	app.Usage = "In-depth inspection of MySQL for system guarantee and daily inspection during major festivals" //应用功能说明
	app.Author = "lianghang"                           //作者
	app.Email = "ywlianghang@gmail.com"                //邮箱
	app.Version = "1.0"                              //版本
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "config,c", //命令名称
			Usage: "Loading a Configuration File. for example： --config /tmp/mycheck.ymal", //命令说明
			Value:  "nil",                                                            //默认值
			Destination: &CheckParameter.Config,                                                                 //赋值
		},
	}
	app.Action = func(c *cli.Context) { //应用执行函数
	}
	app.Run(os.Args)

	aa := os.Args
	for i:= range aa {
		if aa[i] == "--help" || aa[i] == "-h" {
			CheckParameter.helpbool = true
			os.Exit(1)
		}
		if aa[i] == "--version" || aa[i] == "-v"{
			CheckParameter.helpbool = true
			os.Exit(1)
		}
		if strings.Contains(aa[i],"--config") {
			if !strings.Contains(CheckParameter.Config,"/") && !strings.Contains(CheckParameter.Config,"\\"){
				pathdir,err := os.Getwd()
				if err != nil{
					fmt.Println(err)
					os.Exit(1)
				}
				sysType := runtime.GOOS
				if sysType == "linux"{
					CheckParameter.Config = fmt.Sprintf("%s/%s",pathdir,CheckParameter.Config)
				}
				if sysType == "windows"{
					CheckParameter.Config = fmt.Sprintf("%s\\%s",pathdir,CheckParameter.Config)
				}
			}
		}
	}

}
//判断文件是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func ParameterCheck() {
	cliHelp()
	if !CheckParameter.helpbool {
		if CheckParameter.Config == "nil" {
			fmt.Println("The myCheck configuration file is not specified. Use --help to view the parameters")
			os.Exit(0)
		} else {
			if ok, err := pathExists(CheckParameter.Config); !ok && err != nil {
				fmt.Println("The configuration file is not specified correctly. Use --help to view the parameters")
				fmt.Println("error info：", err)
				os.Exit(0)
			}
		}
	}
}


