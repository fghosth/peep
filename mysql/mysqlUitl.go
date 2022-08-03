package mysql

import (
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/fghosth/peep/converter"
	"github.com/fghosth/peep/util"
	"github.com/k0kubun/pp"
	"log"
	"os"
	"strings"
)

var (
	ModelPath        string
	XmlPath          string
	BasePackageName  = "base"
	ModelPackageName string
	ModelFile        = "model.go"
	Mysqluri         string
	BaseFileName     = "base.go"
	MybatisFile      = "createXml_test.go"
	BaseName         = "/base/"
)

func CreateModel(modelPath, basePackageName, modelFile, mysqluri string) []string {
	err := os.Mkdir(modelPath, os.ModePerm)
	if err != nil {
		log.Println("创建目录错误", err.Error())
	}
	err = os.Mkdir(modelPath+BaseName, os.ModePerm)
	if err != nil {
		log.Println("创建目录错误", err.Error())
	}
	t2s := converter.NewTable2Struct()
	// 个性化配置
	t2s.Config(&converter.T2tConfig{
		// 如果字段首字母本来就是大写, 就不添加tag, 默认false添加, true不添加
		RmTagIfUcFirsted: false,
		// tag的字段名字是否转换为小写, 如果本身有大写字母的话, 默认false不转
		TagToLower: false,
		// 字段首字母大写的同时, 是否要把其他字母转换为小写,默认false不转换
		UcFirstOnly: false,
		//// 每个struct放入单独的文件,默认false,放入同一个文件(暂未提供)
		//SeperatFile: false,
		StructNameToHump: true,
	})
	// 开始迁移转换
	err = t2s.
		// 指定某个表,如果不指定,则默认全部表都迁移
		//Table("reports").
		// 表前缀
		//Prefix("prefix_").
		// 是否添加json tag
		EnableJsonTag(true).
		//日期用time.Time
		DateToTime(true).
		// 生成struct的包名(默认为空的话, 则取名为: package model)
		PackageName(basePackageName).
		// tag字段的key值,默认是orm
		//TagKey("orm").
		//TagKey("gm").
		// 是否添加结构体方法获取表名
		RealNameMethod("TableName").
		// 生成的结构体保存路径
		SavePath(modelPath + BaseName + modelFile).
		// 数据库dsn,这里可以使用 t2t.DB() 代替,参数为 *sql.DB 对象
		Dsn(mysqluri).
		// 执行
		Run()
	return t2s.StructList
}

func CreateBaseFile(modelPath, baseName, xmlPath string) error {
	ctx := map[string]interface{}{
		"xmlpath": xmlPath,
	}

	createXmlStr, err := raymond.Render(BASE_TMP, ctx)
	//生成base.go
	err = util.WriteWithIoutil(modelPath+BaseName+baseName, createXmlStr)
	if err != nil {
		return err
	}
	return nil
}

func CreateXmlGoFile(structList []string, basePackageName, xmlPath, modelPath, mybatisFile string) {
	var fstr []string
	for _, value := range structList {
		line := `beans[new(` + value + `).TableName()]=*new(` + value + `)`
		fstr = append(fstr, line)
	}

	ctx := map[string]interface{}{
		"packageName": basePackageName,
		"field":       fstr,
		"xmlpath":     xmlPath,
		"tablestr":    "{{{tableName}}}",
	}

	createXmlStr, err := raymond.Render(CREATEXML_TMP, ctx)
	if err != nil {
		fmt.Println(err)
	}
	err = util.WriteWithIoutil(modelPath+BaseName+mybatisFile, createXmlStr)
	if err != nil {
		log.Println(err)
	}
}

func CreateDao(pName, path string, sl []string, template, xmlPath string) error {
	for _, v := range sl {
		daoName := v + "Dao.go"
		structName := v + "Mapper"
		xmlFile := structName + ".xml"
		tmp := strings.Split(pName, "/")
		packageName := tmp[len(tmp)-1]
		sname, _ := util.FUPer(structName)
		ctx := map[string]interface{}{
			"pName":           packageName, //dao 包名
			"basePackageName": pName,       //要引入的完整路径包名
			"structName":      structName,
			"TableStruct":     v,
			"xmlFile":         xmlFile,
			"sName":           sname,
		}

		createXmlStr, err := raymond.Render(template, ctx)
		if err != nil {
			log.Println(err)
			return err
		}
		if isexist, _ := PathExists(path + "/" + daoName); isexist {
			pp.Println(path + "/" + daoName + "已存在，跳过")
			continue
		}
		err = util.WriteWithIoutil(path+"/"+daoName, createXmlStr)
		if err != nil {
			log.Println(err)
			return nil
		}
	}
	return nil
}

//判断文件或目录是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
