package GoMybatis

import (
	"database/sql"
)

//数据源路由接口
type DataSourceRouter interface {
	//路由规则
	//参数：mapperName mapper文件包名+名称例如（example.ExampleActivityMapper）
	//返回（session,error）路由选择后的session，error异常
	Router(mapperName string, engine SessionEngine) (Session, error)
	//设置sql.DB，该方法会被GoMybatis框架内调用
	SetDB(driverName string, url string, db *sql.DB)

	Name() string
}
