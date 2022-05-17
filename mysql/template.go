package mysql

const (
	CREATEXML_TMP = `package {{{packageName}}}

import (
	"github.com/zhuxiujia/GoMybatis"
    "github.com/beevik/etree"
    "github.com/aymerick/raymond"
	"log"
	"reflect"
    "fmt"
	"testing"
    "bytes"
	"os"
)
var(
	xmlPath="{{{xmlpath}}}/"
)
func TestBybatis(t *testing.T){
	err:=os.Mkdir(xmlPath,os.ModePerm)
	if err!=nil{
		log.Println(err)
	}
	beans:=make(map[string]interface{})
	{{#each field}}
 		{{{this}}}
    {{/each}}
	for k, v := range beans {
		xmlfile := xmlPath + "/" + reflect.TypeOf(v).Name() + "Mapper.xml"
		xml := GoMybatis.CreateXml(k, v)
		ctx := map[string]interface{}{
			"tableName": k,
		}
		teml, err := raymond.Render(XML_TMP, ctx)
		if err != nil {
			fmt.Println(err)
		}
		xml = bytes.Replace(xml, []byte("</mapper>"), []byte(teml), -1)
		srcdoc:=getResultMapElement(xml)
		if isexist, _ := PathExists(xmlfile); isexist {
			log.Println(xmlfile + "已存在，只更新resultmap")
			doc:=updateXmlFile(srcdoc,xmlfile)
			//doc.WriteTo(os.Stdout)
			doc.WriteToFile(xmlfile)
		}else{
			srcdoc.WriteToFile(xmlfile)
		}
	}
}

//xml已存在更新
func updateXmlFile(srcdoc *etree.Document,xmlfile string)*etree.Document{
	doc:=etree.NewDocument()
	doc.ReadFromFile(xmlfile)
	index:=doc.FindElement("//mapper/resultMap").Index()
	doc.FindElement("//mapper").RemoveChildAt(index)
	doc.FindElement("//mapper").InsertChildAt(index,srcdoc.FindElement("//mapper/resultMap"))
	return doc
}


// 处理resultmap id langtype 并返回响应element
func getResultMapElement(data []byte)*etree.Document{
	obj:=etree.NewDocument()
	err := obj.ReadFromBytes(data)
	if err!=nil {
		log.Println("obj.ReadFromBytes(xml)",err)
		return nil
	}
	rmap:=obj.FindElement("//mapper/resultMap")
	rele:=rmap.FindElement("//id")
	if  rele!=nil {//如果找到id element
		rele.CreateAttr("langType","int64")
	}
	return obj
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
<select id="List">
	select * from {{{tablestr}}}
	where delete_flag = 1
	<if test="where != null">
		and  ${where}
	</if>
	order by ${orderField} ${order}
	<if test="page != 0 and size != 0">limit #{offset}, #{size}</if>
</select>
</mapper>
` + "`)"

	DAO_TMP = `package {{{pName}}}
import (
	"{{{basePackageName}}}/base"
	"github.com/zhuxiujia/GoMybatis"
	"sync"
)
//支持基本类型和指针(int,string,time.Time,float...且需要指定参数名称` + "`" + `mapperParams:"name"以逗号隔开，且位置要和实际参数相同)
//参数中包含有*GoMybatis.Session的类型，用于自定义事务
//自定义结构体参数（属性必须大写）
//方法 return 必须包含有error ,为了返回错误信息
type {{{structName}}} struct {
	GoMybatis.SessionSupport                                   //session事务操作 写法1.  ExampleActivityMapper.SessionSupport.NewSession()
	NewSession               func() (GoMybatis.Session, error) //session事务操作.写法2   ExampleActivityMapper.NewSession()
	//模板示例
	FindByID      func(id int) ([]base.{{{TableStruct}}}, error) ` + "`" + `args:"id"` + "`" + `
	Insert      func(arg base.{{{TableStruct}}}) (int64, error) ` + "`" + `args:"arg"` + "`" + `
	InsertBatch func(args []base.{{{TableStruct}}}) (int64, error) ` + "`" + `args:"args"` + "`" + `
	UpdataByID      func(arg base.{{{TableStruct}}}) (int64, error)    ` + "`" + `args:"id"` + "`" + `
	DeleteByID      func(id int) (int64, error)     ` + "`" + `args:"id"` + "`" + `
    List func(offset, size int64,orderField,order,where string)([]base.{{{TableStruct}}}, error) ` + "`" + `args:"offset,size,orderField,order,where"` + "`" + `
}
var (
	{{{sName}}}Once sync.Once
	{{{sName}}} *{{{structName}}}
)
func New{{{structName}}}()(*{{{structName}}},error){
	var err error
	var db *base.DBBase
	{{{sName}}}Once.Do(func() {
		{{{sName}}}=&{{{structName}}}{}
		db,err=base.GetDBBase()
		if err!=nil {
			return
		}
		err=db.WriteMapper({{{sName}}},base.Xmlpath+"/{{{xmlFile}}}")
		if err!=nil {
			return
		}
	})
	if err!=nil {
		return nil,err
	}
	return {{{sName}}},nil
}
`
	BASE_TMP = `package base

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zhuxiujia/GoMybatis"
	"io/ioutil"
	"sync"
	"time"
)

var (
	once sync.Once
	dbBase *DBBase
	defaultMaxIdleConns = 10
	defaultMaxOpenConns = 50
	defaultConnMaxLifetime = 60
	Xmlpath = "{{{xmlpath}}}/"
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
	ConnMaxLifetime int
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
			ConnMaxLifetime: defaultConnMaxLifetime,
			SetLogEnable: false,
			LogFun: nil,
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
		dbBase.db.SetConnMaxLifetime(time.Duration(dbBase.opt.ConnMaxLifetime) * time.Second)
		dbBase.engine.SetLogEnable(dbBase.opt.SetLogEnable)
		//自定义日志实现
		if dbBase.opt.LogFun != nil {
			dbBase.engine.SetLog(&GoMybatis.LogStandard{
				PrintlnFunc: func(messages []byte) {
					dbBase.opt.LogFun(messages)
				},
			})
		}
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
