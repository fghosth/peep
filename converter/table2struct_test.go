package converter

import (
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/beevik/etree"
	"github.com/k0kubun/pp"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	modelFile       = "model.go"
	modelPath       = "/Users/derek/project/peep/dist/mysql"
	basePackageName = "base"
	daoPackageName  = "model"
	dsn             = "webdb:S3wklxG6K7P99aml@tcp(rm-uf673lf2loh08vk0s90130.mysql.rds.aliyuncs.com:3306)/webdb?charset=utf8mb4&parseTime=true"
	mybatisFile     = "createXml_test.go"
	baseName        = "base.go"
	xmlPath         = "/Users/derek/project/demo/peep/mysqlxml/"
)

func TestModifyXml(t *testing.T) {
	xmlpath := "/Users/derek/project/demo/peep/model/CountTaskMapper.xml"
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(xmlpath); err != nil {
		panic(err)
	}
	root := doc.SelectElement("mapper")
	//resultmap := doc.FindElements("//resultMap[@id='BaseResultMap']/id")
	e := doc.FindElement("./mapper/resultMap[@id='BaseResultMap']")
	index := e.Index()
	root.RemoveChild(e)
	root.InsertChildAt(index, e)
	//root.AddChild(e)
	pp.Println("================")
	doc.WriteTo(os.Stdout)
}
func TestXml(t *testing.T) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	doc.CreateProcInst("xml-stylesheet", `type="text/xsl" href="style.xsl"`)

	people := doc.CreateElement("People")
	people.CreateComment("These are all known people")

	jon := people.CreateElement("Person")
	jon.CreateAttr("name", "Jon")

	sally := people.CreateElement("Person")
	sally.CreateAttr("name", "Sally")

	doc.Indent(2)
	doc.WriteTo(os.Stdout)
}

func TestTable2Struct_Run(t1 *testing.T) {
	err := os.Mkdir(modelPath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	t2s := NewTable2Struct()
	// 个性化配置
	t2s.Config(&T2tConfig{
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
		// 是否添加结构体方法获取表名
		RealNameMethod("TableName").
		// 生成的结构体保存路径
		SavePath(modelPath + modelFile).
		// 数据库dsn,这里可以使用 t2t.DB() 代替,参数为 *sql.DB 对象
		Dsn(dsn).
		// 执行
		Run()
	var fstr []string
	for _, value := range t2s.StructList {
		line := `beans[new(` + value + `).TableName()]=*new(` + value + `)`
		fstr = append(fstr, line)
	}

	ctx := map[string]interface{}{
		"packageName": basePackageName,
		"field":       fstr,
		"xmlpath":     xmlPath,
	}

	createXmlStr, err := raymond.Render(CREATEXML_TMP, ctx)
	if err != nil {
		fmt.Println(err)
	}
	err = WriteWithIoutil(modelPath+mybatisFile, createXmlStr)
	if err != nil {
		log.Println(err)
	}
	//生成dao
	err = CreateDao(daoPackageName, "/Users/derek/project/demo/gomybatis/model/", t2s.StructList, DAO_TMP)
	if err != nil {
		log.Println(err)
	}
	//生成base.go
	err = WriteWithIoutil(modelPath+baseName, BASE_TMP)
	if err != nil {
		log.Println(err)
	}
}

func CreateDao(pName, path string, sl []string, template string) error {
	for _, v := range sl {
		daoName := v + "Dao.go"
		structName := v + "Mapper"
		xmlpath := xmlPath + structName + ".xml"
		ctx := map[string]interface{}{
			"packageName": pName,
			"structName":  structName,
			"xmlPath":     xmlpath,
		}

		createXmlStr, err := raymond.Render(template, ctx)
		if err != nil {
			log.Println(err)
			return err
		}
		if isexist, _ := PathExists(path + daoName); isexist {
			pp.Println(path + daoName + "已存在，跳过")
			continue
		}
		err = WriteWithIoutil(path+daoName, createXmlStr)
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

//使用ioutil.WriteFile方式写入文件,是将[]byte内容写入文件,如果content字符串中没有换行符的话，默认就不会有换行符
func WriteWithIoutil(name, content string) error {
	data := []byte(content)
	if err := ioutil.WriteFile(name, data, 0644); err == nil {
		log.Println("写入文件成功:", name)
	} else {
		return err
	}
	return nil
}

const (
	CREATEXML_TMP = `package {{{packageName}}}

import (
	"github.com/zhuxiujia/GoMybatis"
	"reflect"
	"testing"
    "bytes"
)
var(
	xmlPath="{{{xmlpath}}}"
)
func TestBybatis(t *testing.T){
	err:=os.Mkdir(xmlPath,os.ModePerm)
	if err!=nil{
		fmt.Println(err)
	}
	beans:=make(map[string]interface{})
	{{#each field}}
 		{{{this}}}
    {{/each}}
	for k, v := range beans {
		xmlfile:=xmlPath+reflect.TypeOf(v).Name()+"Mapper.xml"
		xml:=GoMybatis.CreateXml(k, v)
		xml = bytes.Replace(xml,[]byte("</mapper>"),[]byte(XML_TMP),-1)
 		GoMybatis.OutPutXml(xmlfile, xml)
	}
}

const (
XML_TMP=` + "`" + `
<!-- =============================！！！！以上内容不要修改！！！！！================================================= -->
<!--模板标签: columns wheres sets 支持逗号,分隔表达式，*?* 为判空表达式-->
<!--插入模板:默认id="insertTemplate,test="field != null",where自动设置逻辑删除字段,支持批量插入" -->
<insertTemplate id="Insert" />
<!--查询模板:默认id="selectTemplate,where自动设置逻辑删除字段-->
<selectTemplate id="FindByID" wheres="id?id = #{id}" />
<!--更新模板:默认id="updateTemplate,set自动设置乐观锁版本号-->
<updateTemplate id="UpdataByID"  wheres="id?id = #{id}" />
<!--删除模板:默认id="deleteTemplate,where自动设置逻辑删除字段-->
<deleteTemplate id="DeleteByID" wheres="id?id= #{id}" />
<!--批量插入: 因为上面已经有id="insertTemplate" 需要指定id -->
<insertTemplate id="InsertBatch"/>
<!--统计模板:-->
<!--	<selectTemplate id="selectCountTemplate" columns="count(*)" wheres="reason?reason = #{reason}"/>-->
</mapper>
` + "`)"

	DAO_TMP = `
package {{{packageName}}}
import (
	"gomybatis/model/base"
	"github.com/zhuxiujia/GoMybatis"
)
//支持基本类型和指针(int,string,time.Time,float...且需要指定参数名称` + "`" + `mapperParams:"name"以逗号隔开，且位置要和实际参数相同)
//参数中包含有*GoMybatis.Session的类型，用于自定义事务
//自定义结构体参数（属性必须大写）
//方法 return 必须包含有error ,为了返回错误信息
type {{{structName}}} struct {
	GoMybatis.SessionSupport                                   //session事务操作 写法1.  ExampleActivityMapper.SessionSupport.NewSession()
	NewSession               func() (GoMybatis.Session, error) //session事务操作.写法2   ExampleActivityMapper.NewSession()
	//模板示例
	FindByID      func(id int) ([]base.CountTask, error) ` + "`" + `mapperParams:"id"` + "`" + `
	Insert      func(arg base.CountTask) (int64, error) ` + "`" + `mapperParams:"arg"` + "`" + `
	InsertBatch func(args []base.CountTask) (int64, error) ` + "`" + `mapperParams:"args"` + "`" + `
	UpdataByID      func(arg base.CountTask) (int64, error)    ` + "`" + `mapperParams:"id"` + "`" + `
	DeleteByID      func(id int) (int64, error)     ` + "`" + `mapperParams:"id"` + "`" + `
}

func New{{{structName}}}()(*{{{structName}}},error){
	bm:={{{structName}}}{}
	db,err:=base.GetDBBase()
	if err!=nil {
		return nil,err
	}
	err=db.WriteMapper(&bm,"{{{xmlPath}}}")
	if err!=nil {
		return nil,err
	}
	return &bm,nil
}
`
	BASE_TMP = `
package base

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zhuxiujia/GoMybatis"
	"io/ioutil"
	"sync"
)

var (
	once sync.Once
	dbBase *DBBase
	defaultMaxIdleConns = 10
	defaultMaxOpenConns = 50
)

type DBBase struct{
	engine *GoMybatis.GoMybatisEngine
	mysqlUri string
	db *sql.DB
	opt Option
}

type Option struct{
	MaxIdleConns int
	MaxOpenConns int
	SetLogEnable bool //是否自定义日志系统
	LogFun func(msg []byte)//自定义日志
}

type DBOption func(option *Option)


func GetDBBase()(*DBBase,error){
	if dbBase==nil || dbBase.db==nil {
		return nil,errors.New("DBBase 未初始化，请先调用NewDBBase")
	}
	return dbBase,nil
}


func NewDBBase(mysqlUri string, opt ...func(option *Option))(*DBBase,error) {
	if mysqlUri == "" {
		return nil,errors.New("数据库链接地址错误")
	}
	var err error
	once.Do(func() {
		option:=Option{
			MaxIdleConns: defaultMaxIdleConns,
			MaxOpenConns: defaultMaxOpenConns,
			SetLogEnable: false,
			LogFun: func(msg []byte) {},
		}
		for _, f := range opt {
			f(&option)
		}
		engine:=GoMybatis.GoMybatisEngine{}.New()
		dbBase=&DBBase{
			engine:&engine,
			mysqlUri:mysqlUri,
			opt: option,
		}
		dbBase.db, err = dbBase.engine.Open("mysql", mysqlUri) //此处请按格式填写你的mysql链接，这里用*号代替
		if err!=nil {
			return
		}
		dbBase.db.SetMaxIdleConns(dbBase.opt.MaxIdleConns)
		dbBase.db.SetMaxOpenConns(dbBase.opt.MaxOpenConns)
		//自定义日志实现
		dbBase.engine.SetLogEnable(dbBase.opt.SetLogEnable)
		dbBase.engine.SetLog(&GoMybatis.LogStandard{
			PrintlnFunc: func(messages []byte) {
				dbBase.opt.LogFun(messages)
			},
		})
	})
	if err != nil {
		return nil,err
	}
	return dbBase,nil
}
//返回db链接以便自定义
func (this DBBase)GetDB()*sql.DB{
	return this.db
}
//返回engin以便自定义
func (this DBBase)GetEngin()*GoMybatis.GoMybatisEngine{
	return this.engine
}

func (this *DBBase)WriteMapper(obj interface{},xmlfile string)error{
	//读取mapper xml文件
	bytes, err := ioutil.ReadFile(xmlfile)
	if err!=nil {
		return err
	}
	//设置对应的mapper xml文件
	this.engine.WriteMapperPtr(obj, bytes)
	return nil
}


//设置最大空闲连接数
func WithMaxIdleConns(num int) DBOption {
	return func(option *Option) {
		option.MaxIdleConns = num
	}
}
//设置最大连接数
func WithMaxOpenConns(num int) DBOption {
	return func(option *Option) {
		option.MaxOpenConns = num
	}
}
//设置是否使用自定义日志
func WithSetLogEnable(flag bool) DBOption {
	return func(option *Option) {
		option.SetLogEnable = flag
	}
}
//设置自定义日志方法
func WithLogFun(f func(msg []byte)) DBOption {
	return func(option *Option) {
		option.LogFun= f
	}
}

`
)

func Test_removeBrackets(t *testing.T) {
	s := "bigint(23) unsigned"
	res := removeBrackets(s)
	log.Println(res)
}
