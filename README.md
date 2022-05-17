## 功能
1. 生成dci项目结构
2. 根据mysql生成数据库操作脚手架
3. 初始化githooks。功能如下
   * git commit 提交时检查注释规范，不符合不能commit。
   * git commit 提交检查go代码质量。
   * git commit 提交检查单测覆盖率(默认达到80%)。
4. makefile。功能如下：
   * make init 初始化项目，安装githook等
   * make build 编译项目。生成在dist目录下
   * make run 编译运行项目。
5. Dockerfile。docker镜像
   
## 安装
* mac
```
brew install fghosth/fghosth/peep
```
* 源码安装
```
go get github.com/fghosth/peep
```
## 使用
```
peep -h

Usage:
  mysqlTool [command]

Available Commands:
  help        Help about any command
  mcreate     create code By mysql
  pcreate     create project struct
```
* 生成项目结构
```
peep pcreate --module "fghosth.net/reportssss" --path /Users/derek/project/demo/gomybatis/tmp
```
* mysql脚手架生成
```
peep mcreate --mp /Users/derek/project/demo/gomybatis/model --mpn "gomybatis/model" --uri "root:zaqwedcxs@tcp(localhost:3306)/fghosth_reports?charset=utf8mb4&parseTime=true" --xp /Users/derek/project/demo/gomybatis/mysqlxml
```

>【mp】生成的文件目录；【mpn】生成的文件包名；【rui】mysql的连接；【xp】生成xml目录 

## 项目版本号
在build时指定版本号
```
1. make build Version=v1.0.0
2. docker build --build-arg version=v1.0.2 -t XXXXX .
```
获取版本,输出程序版本号，go编译版本信息，编译时间
```bazaar
server version
```

## 代码质量检查
请安装ide的代码检查插件，githook使用的是golangci-lint。
支持的IDE：
* Go for Visual Studio Code
* Sublime
* GoLand
* GNU Emacs
* VIM
* Atom
## git comiit规范辅助程序 安装使用
注释向导，生成符合标准的规范注释
* 安装 
```
make gitcommtool
```
* 使用
```
git cz -a
```
* 根据注释生成changelog.md
```bazaar
make gitlog
```
## git commit 规范
### 格式
```
type(scope) : subject
```
### type（必须） : commit 的类别，只允许使用下面几个标识：
|类型|description|描述|
|:---:|:---:|:---:|
|feat	|A new feature	|新功能|
|fix	|A bug fix	|修复 bug|
|docs	|Documentation only changes	|文档修改|
|style	|Changes that do not affect the meaning of the code (white-space, formatting, missing semicolons, etc)	|格式（不影响代码运行的变动）|
|refactor	|A code change that neither fixes a bug nor adds a feature	|重构|
|perf|	A code change that improves performance|	提高性能|
|test|	Adding missing tests or correcting existing tests|	添加缺失测试或更正现有测试|
|build|	Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)	|依赖的外部资源变化|
|chore	|Other changes that don't modify src or test files	|构建过程或辅助工具的变动|
|revert|	Reverts a previous commit	|恢复先前的提交|

### scope（可选） : 用于说明 commit 影响的范围，比如数据层、控制层、视图层等等，视项目不同而不同。
### subject（必须） : commit 的简短描述，不超过50个字符。

* [gitflow流程说明](./docs/gitflow.md)
* [golangci 代码质量检查](https://github.com/golangci/golangci-lint)