package cmd

import (
	"fmt"
	"github.com/fghosth/peep/mysql"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
)

var createByMysqlCmd = &cobra.Command{
	Use:     "mcreate",
	Short:   "create code By mysql",
	Example: "projectCreate mcreate --mp /Users/derek/project/demo/gomybatis/model --mpn \"gomybatis/model\" --uri \"root:zaqwedcxs@tcp(localhost:3306)/fghosth_reports?charset=utf8mb4&parseTime=true\" --xp /Users/derek/project/demo/gomybatis/mysqlxml",
	Run: func(cmd *cobra.Command, args []string) {

		sl := mysql.CreateModel(mysql.ModelPath, mysql.BasePackageName, mysql.ModelFile, mysql.Mysqluri)
		err := mysql.CreateDao(mysql.ModelPackageName, mysql.ModelPath, sl, mysql.DAO_TMP, mysql.XmlPath)
		if err != nil {
			log.Println("生成DAO错误" + err.Error())
		}
		mysql.CreateXmlGoFile(sl, mysql.BasePackageName, mysql.XmlPath, mysql.ModelPath, mysql.MybatisFile)
		err = mysql.CreateBaseFile(mysql.ModelPath, mysql.BaseFileName, mysql.XmlPath)
		if err != nil {
			log.Println("生成base.go错误" + err.Error())
		}
		gocmd := exec.Command("gofmt", "-w", "-s", mysql.ModelPath)
		err = gocmd.Run()
		//gocmd = exec.Command("gofmt", "-w","-s",mysql.ModelPath+"/base")
		//err = gocmd.Run()
		if err != nil {
			fmt.Println("Execute Command failed:" + err.Error())
			return
		}
	},
}

func init() {
	createByMysqlCmd.PersistentFlags().StringVar(
		&mysql.ModelPath,
		"mp",
		"./",
		"model path",
	)
	createByMysqlCmd.PersistentFlags().StringVar(
		&mysql.Mysqluri,
		"uri",
		"root:zaqwedcxs@tcp(localhost:3306)/fghosth_reports?charset=utf8mb4&parseTime=true",
		"mysql uri",
	)
	createByMysqlCmd.PersistentFlags().StringVar(
		&mysql.ModelPackageName,
		"mpn",
		"model",
		"ModelPackageName",
	)
	createByMysqlCmd.PersistentFlags().StringVar(
		&mysql.XmlPath,
		"xp",
		"./",
		"xml path",
	)
	RootCmd.AddCommand(createByMysqlCmd)
}
