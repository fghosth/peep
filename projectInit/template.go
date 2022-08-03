package projectInit

const (
	DOCKERFILE_TMP = `FROM golang:1.17.4-alpine3.15 as builder

RUN apk add make
RUN mkdir /go/src/app
ADD . /go/src/app
RUN cd /go/src/app \
    && make Version=v1.0.0 linux_build


FROM alpine:latest
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache ca-certificates  bash
RUN mkdir -p /app/logs
RUN mkdir -p /app/profile
COPY --from=builder /go/src/app/dist/server  /app/

COPY entrypoint.sh /bin
RUN chmod +x /bin/entrypoint.sh
RUN chmod +x /app/server
WORKDIR /app
EXPOSE 8000
CMD ["entrypoint.sh","-"]
`

	ENTRYPOINT_TMP = `#!/bin/bash
set -e


if [ "${1:0:1}" = '-' ]; then
    set -- app "$@" #如果第一个参数的第一个字符是【-】,在所有参数前添加segment 以空格分割
fi

if [ "$1" = 'app' ]; then
    mkdir -p /app/logs
    mkdir -p /app/profile
    touch /app/logs/app.log
    touch /app/logs/costEngin.log
	mkdir -p /app/logs/
	/app/server up --conf /app/config/config.yaml >> /app/logs/demo.log 2>&1 &
	sleep 1
	tail -qf /app/logs/*.log
fi

`

	README_TMP = `
[![Workflow](https://github.com/Andrew-M-C/go.jsonvalue/actions/workflows/go_test_general.yml/badge.svg)](https://github.com/fghosth/xlib/actions/workflows/go_yml)
[![Coverage Status](https://coveralls.io/repos/github/fghosth/utils/badge.svg?branch=main)](https://coveralls.io/github/fghosth/xlib?branch=main)
[![codebeat badge](https://codebeat.co/badges/432bc9f4-aee8-493b-ba38-fe480fb5fb4f)](https://codebeat.co/projects/github-com-fghosth-xlib-main)
[![Go Report Card](https://goreportcard.com/badge/github.com/fghosth/utils)](https://goreportcard.com/report/github.com/fghosth/xlib)

[![Go Reference](https://pkg.go.dev/badge/github.com/fghosth/utils.svg)](https://pkg.go.dev/github.com/fghosth/xlib)
[![License](https://img.shields.io/badge/license-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
[![Latest](https://img.shields.io/badge/latest-v1.1.1-blue.svg)](https://github.com/Andrew-M-C/go.jsonvalue/tree/v1.0.0)

# 项目名称


## 一.目的

## 二.设计目标


## 目录结构
1. ui: 是用户接口层，主要用于处理用户发送的Restful请求和解析用户输入的配置文件等，并将信息传递给层application的接口。
2. Application: 负责多进程管理及调度、多线程管理及调度、多协程调度和维护业务实例的状态模型。当调度层收到用户接口层的请求后，委托Context层与本次业务相关的上下文进行处理
3. context 是环境层，以上下文为单位，将Domain层的领域对象cast成合适的role，让role交互起来完成业务逻辑
4. Domain 领域层，定义领域模型，不仅包括领域对象及其之间关系的建模，还包括对象的角色role的显式建模
5. Infrastructure 是基础实施层，为其他层提供通用的技术能力：业务平台，编程框架，持久化机制，消息机制，第三方库的封装，通用算法，等等。
6. doc 文档
7. cmd 主程序和其他可编译的程序

## ui接口采用grpc

[详细文档](doc/技术方案.md)

### 待修改
`

	MAKEFILE_TMP = `default:
	@echo 'Usage of make: [ build | linux_build | windows_build | clean ]'
build: confile
	@go build -mod=vendor -tags netgo  -o ./dist/server ./cmd
linux_build: confile
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor  -tags netgo -o ./dist/server ./cmd
windows_build: confile
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -mod=vendor -o ./dist/server.exe ./cmd
run: build
	@./dist/server up --conf ./dist/conf/config.yaml
stop:
	@kill -9 $(ps -ef | grep "./dist/server" |gawk '$0 !~/grep/ {print $2}' |tr -s '\n' ' ')
init:
	@git init
	@chmod +x .githooks/go_pre_commit.sh
	@cp .githooks/pre-commit .git/hooks && chmod +x ./.git/hooks/pre-commit && go mod tidy && go mod vendor
	@cp .githooks/commit-msg .git/hooks && chmod +x ./.git/hooks/commit-msg
gitcommtool:
	@brew install node
	@npm install -g cnpm --registry=https://registry.npm.taobao.org
	@cnpm install -g commitizen
	@cnpm install -g conventional-changelog
	@cnpm install -g conventional-changelog-cli
	@cnpm init --yes
	@commitizen init cz-conventional-changelog --save --save-exact
gitlog:
	@conventional-changelog -p angular -i CHANGELOG.md -s
install: build
	@mv ./dist/server $(GOPATH)/dist/

clean:
	@rm -rf ./dist
confile:
	@mkdir -p ./dist/conf
	@cp ./conf/config.yaml ./dist/conf

`
	LOG_TMP = `package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)


func InitLogger(logpath string, loglevel string,write2File bool) *zap.Logger {

	hook := lumberjack.Logger{
		Filename:   logpath, // 日志文件路径
		MaxSize:    128,     // megabytes
		MaxBackups: 300,      // 最多保留300个备份
		MaxAge:     7,       // days
		Compress:   true,    // 是否压缩 disabled by default
	}
	var w zapcore.WriteSyncer
	if write2File {
		w = zapcore.AddSync(&hook)
	}else{
		w =zapcore.AddSync(os.Stdout)
	}
	// 设置日志级别,debug可以打印出info,debug,warn；info级别可以打印warn，info；warn只能打印warn
	// debug->info->warn->error
	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}
	// 设置日志级别
	//atom := zap.NewAtomicLevelAt(level)
	//
	//config := zap.Config{
	//	Level:            atom,                                                // 日志级别
	//	Development:      true,                                                // 开发模式，堆栈跟踪
	//	Encoding:         "json",                                              // 输出格式 console 或 json
	//	EncoderConfig:    encoderConfig,                                       // 编码器配置
	//	InitialFields:    map[string]interface{}{"serviceName": "spikeProxy"}, // 初始化字段，如：添加一个服务器名称
	//	OutputPaths:      []string{"stdout", "./logs/spikeProxy.log"},         // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
	//	ErrorOutputPaths: []string{"stderr"},
	//	DisableCaller:false,
	//}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		level,
	)
	logger := zap.New(core,zap.AddCaller())
	logger.Info("DefaultLogger init success")

	return logger
}
`
	CONF_YAML_TMP = `Application:
  RateLimit: 1000 #每秒QPS  单位 个/秒
  MaxWorkerTask: 50 #每个worker最大同时处理任务
  #本程序最大连接数
  Maxconn: 10
  #程序缓存大小100MB 100*1024*1024
  CacheSize: 104857600
  #grpc传输消息大小15MB 15*1024*1024
  GrpcMaxSize: 15728640
  #callback 超时时间 小时
  CallBackMaxElapsedTime: 24
  #最大缓存字节数 100MB 100*1024*1024
  MaxObject: 104857600
 
Log:
  #设置日志级别,debug可以打印出info,debug,warn；info级别可以打印warn，info；warn只能打印warn
  Loglevel: "debug"
  LogFile: "XXX.log"
  AppLogFile: "app.log"
  #日志是否写入文件,如果不写如文件则输出到console
  WriteLog2File: true
Debug:
  ISCPUPprof: true
  CPUPprofPath: './profile/cpu.pprof'
  #CPUPprofPath: '/app/profile/cpu.pprof'
  IsMemPprof: true
  MemPprofPath: './profile/mem.pprof'
#MemPprofPath: '/app/profile/mem.pprof'
Mysql:
  #mysql 链接
  DBName: "fghosth_xxx"
  MysqlURI: "root:zaqwedcxs@tcp(localhost:3306)/fghosth?charset=utf8mb4&parseTime=true"
  #mysql最大连接数
  MysqlMaxConn: 100
  #mysql最大空闲链接
  MysqlMaxIdleConns: 10
  MysqlXMLFile: "./conf/mysqlxml" #mysql xml 文件目录

`

	CONF_TMP = `package conf

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
    "github.com/shima-park/agollo"
	"{{{module}}}/infra/util"
	"go.uber.org/zap"
	"log"
	"strconv"
	"bytes"
	"os"
    "{{{module}}}/infra/apollo"
	"{{{module}}}/infra/util"
)

var (
	Viper  = viper.New()
	CFG    *Config
	Logger *zap.Logger
    ApolloAddr      = "apolloconfig-dev.fghosth.cn"
	ApolloNamespace = "valhalla.yml"
	ApolloAppID     = "devops-app-20211228"
	ApolloToken     = "8edd9ed05d85706efa6c7b0ecd710fa7cb574045"
	ApolloUser      = "fanhua"
)

type Config struct {
	Application
	Log
	Debug
	Kafka
	Redis
	Rabbitmq
	Mysql
}

type Application struct {
	Host                   string //本机的唯一id
	LockTime               int    //卡券锁定时间
	MaxWorkerTask          int    //#每个worker最大同时处理任务
	CheckInterval          int    //检查任务时间间隔
	OpentracAddress        string
	OpentraceServiceName   string
	RateLimit              int  //访问频率 单位 个/秒
	GrpcMaxSize            int  //grpc传输消息大小
	MaxObject              int  //最大对象缓存数量
	LeaderAsWorker         bool //leader节点是否也承担worker的任务
	CallBackMaxElapsedTime int  //callback最大尝试时间
	CacheSize              int  //程序缓存大小

	
}

type Log struct {
	WriteLog2File bool   //日志是否写入文件,如果不写如文件则输出到console
	Loglevel      string //debug可以打印出info,debug,warn；info级别可以打印warn，info；warn只能打印warn
	LogFile       string //业务日志文件
	AppLogFile    string //程序运行访问日志文件
}
type Debug struct {
	ISCPUPprof   bool
	IsMemPprof   bool
	CPUPprofPath string
	MemPprofPath string
}
type Kafka struct {
	KafkaADDR string
	User      string
	Password  string
}
type Redis struct {
	RedisURL string
	RedisPWD string
	RedisDB  int //数据库
}
type Rabbitmq struct {
	RabbitMQURI string //链接地址
	Exname      string //exchange name
	Extype      string //exchange type
	QueueName   string //queue name
	RouteKey    string //routekey
}

type Mysql struct {
	DBName            string
	MysqlMaxConn      int    //mysql最大连接数
	MysqlMaxIdleConns int    //mysql最大空闲链接
	MysqlURI          string //mysql uri
	MysqlXMLFile      string
}

// InitConfig 初始化配置文件,初始化全局logger。会根据配置文件变化动态改变响应配置,全局logger.
func (cfg Config) InitConfig(file string) {
	Viper.SetConfigType("yaml")
	Viper.AddConfigPath(".")
	cfg.load(file)
	uuid, _ := util.GetUUID()
	CFG.Host = strconv.FormatUint(uuid, 10)
	//监控配置文件变化
	Viper.WatchConfig()
	Viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("配置文件变更，重新生效")
		cfg.load(file)
		log.Println(CFG)
	})
}

func (cfg Config) load(file string) {
	Viper.SetConfigFile(file)
	err := Viper.ReadInConfig() // 读取配置文件
	if err != nil {             // 加载配置文件错误
		log.Println("配置文件加载错误",err)
	}
	CFG = &Config{}
	err = Viper.Unmarshal(CFG)
	if err != nil {
		log.Println(err)
	} else {
		Logger = util.InitLogger(CFG.LogFile, CFG.Loglevel,CFG.WriteLog2File)
	}
}

//更新配置文件
func updateConfig(ns string, oldVal, newVal agollo.Configurations, err error) error {
	//pp.Println(ns, oldVal, newVal, err)
	Viper.SetConfigType("yaml")
	for k, v := range newVal {
		if oldVal[k] == v { //跳过未变动配置
			continue
		}
		if v1, ok := v.(string); ok {
			err = Viper.ReadConfig(bytes.NewBufferString(v1))
			if err != nil {
				return err
			}
			//更新全局配置
			err = Viper.Unmarshal(&CFG)
			if err != nil {
				return err
			}
			//更新logger
			Logger = util.InitLogger(CFG.Log.LogFile, CFG.Log.Loglevel, CFG.Log.WriteLog2File)
			Logger.Debug("=======config changed=========:", zap.Any("conf", CFG))
		}
	}
	return nil
}

// InitConfigApollo 初始化配置文件
func (cfg Config) InitConfigApollo() {
	//加载环境变量中的配置
	if os.Getenv("ApolloAddr") != "" {
		ApolloAddr = os.Getenv("ApolloAddr")
	}
	if os.Getenv("ApolloNamespace") != "" {
		ApolloNamespace = os.Getenv("ApolloNamespace")
	}
	if os.Getenv("ApolloAppID") != "" {
		ApolloAppID = os.Getenv("ApolloAppID")
	}
	if os.Getenv("ApolloToken") != "" {
		ApolloToken = os.Getenv("ApolloToken")
	}
	if os.Getenv("ApolloUser") != "" {
		ApolloUser = os.Getenv("ApolloUser")
	}
	Viper.SetConfigType("yaml")
	apo := apollo.NewApolloConfig(ApolloAddr, ApolloToken, ApolloUser, true)
	res, err := apo.GetConfigByNS(ApolloAppID, ApolloNamespace)
	if err != nil {
		log.Fatal(err)
	}
	err = Viper.ReadConfig(bytes.NewBufferString(res["content"]))
	if err != nil {
		log.Fatal(err)
	}
	err = Viper.Unmarshal(&CFG)
	if err != nil {
		log.Fatal(err)
	} else {
		//pp.Println(conf.CFG.LogFile, conf.CFG.Loglevel, conf.CFG.WriteLog2File)
		Logger = util.InitLogger(CFG.Log.LogFile, CFG.Log.Loglevel, CFG.Log.WriteLog2File)
		//conf.Logger.Info("=========",zap.Any("asd","asdf"))
	}
	//监听
	go func() {
		defer func() {
			err = apo.WatchConfig(ApolloAppID, ApolloNamespace, updateConfig)
			if err != nil {
				log.Println(err)
			}
		}()
		err = apo.WatchConfig(ApolloAppID, ApolloNamespace, updateConfig)
		if err != nil {
			log.Println(err)
		}
	}()
}
`
	GOMOD_TMP = `module {{{module}}} 

go 1.17

`
	GITIGNORE_TMP = `*.idea
dist
*.log
demo
.DS_Store
reports.log
*.pdf
node_modules
package.json
package-lock.json
cover.out
`
	CMD_MAIN_TMP = `package main

import (
   "{{{cmdPName}}}"
	_ "github.com/mkevac/debugcharts"
	"log"
	"net/http"
	//_ "net/http/pprof"
)
var (
	Version   string
	BuildTime string
	GoVersion string
)
func main(){
	cmd.GoVersion = GoVersion
	cmd.BuildTime = BuildTime
	cmd.Version = Version
	go func() {
		err := http.ListenAndServe("0.0.0.0:6060", nil)
		if err != nil {
			log.Println(err)
		}
	}()
	cmd.Execute()
}
`

	CMD_VERSION_TMP = `package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Version   string
	BuildTime string
	GoVersion string
)
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "show version",
	Example: "./server version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:"+Version)
		fmt.Println("BuildTime:"+BuildTime)
		fmt.Println("GoVersion:"+GoVersion)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

`

	CMD_ROOT_TMP = `package cmd

import (
	"fmt"
	"github.com/fghosth/lile"
	"os"
)

var RootCmd = lile.BaseCommand("peep", "创建")
var cfgFile string
func Execute() {
	RootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"conf",
		"./conf/Default_config.yaml",
		"config file path",
	)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
`

	CMD_UP_TMP = `package cmd

import (
	"github.com/spf13/cobra"
    "log"
	//"{{{module}}}/infra/swagger"
	"{{{module}}}/conf"
)

var createByMysqlCmd = &cobra.Command{
	Use:   "create",
	Short: "create code By mysql",
	Example:"projectCreate create --mp /Users/derek/project/demo/gomybatis/model --mpn \"gomybatis/model\" --uri \"root:zaqwedcxs@tcp(localhost:3306)/fghosth_reports?charset=utf8mb4&parseTime=true\" --xp /Users/derek/project/demo/gomybatis/mysqlxml",
	Run:func(cmd *cobra.Command, args []string) {
		//捕捉退出信号
		c := make(chan os.Signal, 1)
		//signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGUSR1)
		go func() {
			si := <-c
			pp.Println("退出信号", si)
			//lile.Shutdown()
			os.Exit(0)
		}()
		var err error
		//配置文件 初始化配置，全局log
		log.Println("初始化配置文件")
		cnf := conf.Config{}
		if cfgFile == "./conf/Default_config.yaml" { //apollo 模式
			log.Println("apollo 加载配置文件")
			cnf.InitConfigApollo()
		} else { //配置文件模式
			log.Println("文件 加载配置文件")
			cnf.InitConfig(cfgFile)
		}

		//swagger ui
		//go func() {
		//	log.Println("启动swagger")
		//	uri := "http://localhost:1323/docs/valhalla.swagger.json"
		//	//path := "/Users/derek/project/devops/valhalla/ui/grpc/valhalla"
		//	host := ":1323"
		//	swagger.SwaggerUi(uri, host)
		//}()

	},
}
func init() {
	//createByMysqlCmd.PersistentFlags().StringVar(
	//	&mysql.XmlPath,
	//	"xp",
	//	"./",
	//	"xml path",
	//)
	RootCmd.AddCommand(createByMysqlCmd)
}
`
	UTIL_UUID_TMP = `package util

import (
	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

func init() {
	//st.MachineID = awsutil.AmazonEC2MachineID
	sf = sonyflake.NewSonyflake(sonyflake.Settings{})
	if sf == nil {
		panic("sonyflake not created")
	}
}

func GetUUID()(uint64,error){
	return sf.NextID()
}
`
	GRPC_OPENTRACE_SERVER_TMP = `package grpc_mw

import (
	"context"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
    "github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"sync"
	"time"
)

var (
	servertrace *ServerTrace
	once sync.Once
)

type ServerTrace struct{
	address string
	serverName string
	Tracer opentracing.Tracer
	Closer io.Closer
	TraceOption
}

type TraceOption struct {
	opentraceType string
	opentraceParam float64
}

type OpenTraceOption func(option *TraceOption)

func GetServerTrace()(*ServerTrace,error){
	if servertrace !=nil && servertrace.Tracer!=nil {
		return servertrace, nil
	}
	return nil,errors.New("ServerTrace 未初始化,请先运行NewServerTrace")
}


func NewServerTrace(address,serverName string,opt ...OpenTraceOption)(*ServerTrace,error){
	var err error
	once.Do(func() {
		option := TraceOption{
			opentraceType: "const",
			opentraceParam:  1,
		}
		for _, fn := range opt {
			fn(&option)
		}
		servertrace=&ServerTrace{
			address:    address,
			serverName: serverName,
			TraceOption:option,
		}
		servertrace.Tracer,servertrace.Closer,err=InitJaeger(serverName,address,option.opentraceType,option.opentraceParam)
	})
	if err!=nil {
		return nil,err
	}
	return servertrace,nil
}

func (this ServerTrace)ServerOption(tracer opentracing.Tracer) grpc.ServerOption {
	return grpc.UnaryInterceptor(this.JaegerGrpcServerInterceptor)
}

type TextMapReader struct {
	metadata.MD
}


//读取metadata中的span信息
func (t TextMapReader) ForeachKey(handler func(key, val string) error) error { //不能是指针
	for key, val := range t.MD {
		for _, v := range val {
			if err := handler(key, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (this ServerTrace)JaegerGrpcServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	//从context中获取metadata。md.(type) == map[string][]string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		//如果对metadata进行修改，那么需要用拷贝的副本进行修改。（FromIncomingContext的注释）
		md = md.Copy()
	}
	carrier := TextMapReader{md}
	tracer := this.Tracer
	spanContext, e := tracer.Extract(opentracing.TextMap, carrier)
	if e != nil {
		fmt.Println("Extract err:", e)
	}
	//span := tracer.StartSpan(info.FullMethod, ext.RPCServerOption(spanContext),ext.SpanKindRPCServer,opentracing.FollowsFrom(spanContext))
	span := tracer.StartSpan(info.FullMethod, ext.RPCServerOption(spanContext),opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},ext.SpanKindRPCServer)
	span.SetTag("req",req)
	span.SetTag("md",carrier)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return handler(ctx, req)
}

func InitJaeger(service string, jaegerAgentHost string,OpentraceType string,opentraceParam float64) (tracer opentracing.Tracer, closer io.Closer, err error) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  OpentraceType,
			Param: opentraceParam,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:jaegerAgentHost,
		},
	}
	tracer, closer, err = cfg.NewTracer(config.Logger(jaeger.StdLogger), config.Metrics(metrics.NullFactory))
	//opentracing.SetGlobalTracer(tracer)
	return tracer, closer, err
}

//设置 OpentraceType
func WithOpentraceType(t string) OpenTraceOption {
	return func(option *TraceOption) {
		option.opentraceType = t
	}
}


//设置 OpentraceParam
func WithOpentraceParam(p float64) OpenTraceOption {
	return func(option *TraceOption) {
		option.opentraceParam = p
	}
}
`
	GRPC_OPENTRACE_CLIENT_TMP = `package grpc_mw

import (
	"context"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"sync"
)

var (
	clientTrace *ClientTrace
	onceCT sync.Once
)
type ClientTrace struct{
	address string
	serverName string
	Tracer opentracing.Tracer
	Closer io.Closer
	TraceOption
}
func GetClientTrace()(*ClientTrace,error){
	if clientTrace !=nil && clientTrace.Tracer!=nil {
		return clientTrace, nil
	}
	return nil,errors.New("ClientTrace 未初始化,NewClientTrace")
}
func NewClientTrace(address,serverName string,opt ...OpenTraceOption)(*ClientTrace,error){
	var err error
	onceCT.Do(func() {
		option := TraceOption{
			opentraceType: "const",
			opentraceParam:  1,
		}
		for _, fn := range opt {
			fn(&option)
		}
		clientTrace=&ClientTrace{
			address:    address,
			serverName: serverName,
			TraceOption:option,
		}
		clientTrace.Tracer,clientTrace.Closer,err=InitJaeger(serverName,address,option.opentraceType,option.opentraceParam)
	})
	if err!=nil {
		return nil,err
	}
	return clientTrace,nil
}


func (this ClientTrace)ClientDialOption() grpc.DialOption {
	return grpc.WithUnaryInterceptor(this.JaegerGrpcClientInterceptor)
}
type TextMapWriter struct {
	metadata.MD
}
//重写TextMapWriter的Set方法，我们需要将carrier中的数据写入到metadata中，这样grpc才会携带。
func (t TextMapWriter) Set(key, val string) {
	//key = strings.ToLower(key)
	t.MD[key] = append(t.MD[key], val)
}

func (this ClientTrace)JaegerGrpcClientInterceptor (ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	var parentContext opentracing.SpanContext
	//先从context中获取原始的span
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		parentContext = parentSpan.Context()
	}
	//tracer := opentracing.GlobalTracer()
	tracer := this.Tracer
	span := tracer.StartSpan(method, opentracing.ChildOf(parentContext),opentracing.Tag{Key: string(ext.Component), Value: "gRPC"}, ext.SpanKindRPCClient)
	defer func() {
		span.Finish()
	}()
	//从context中获取metadata。md.(type) == map[string][]string
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		//如果对metadata进行修改，那么需要用拷贝的副本进行修改。（FromIncomingContext的注释）
		md = md.Copy()
	}
	//定义一个carrier，下面的Inject注入数据需要用到。carrier.(type) == map[string]string
	//carrier := opentracing.TextMapCarrier{}
	carrier := TextMapWriter{md}
	//将span的context信息注入到carrier中
	e := tracer.Inject(span.Context(), opentracing.TextMap, carrier)
	if e != nil {
		fmt.Println("tracer Inject err,", e)
	}
	//创建一个新的context，把metadata附带上
	ctx = metadata.NewOutgoingContext(ctx, md)

	return invoker(ctx, method, req, reply, cc, opts...)
}
`
	GRPC_LIMITER_TMP = `package grpc_mw

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Limiter defines the interface to perform request rate limiting.
// If Limit function return true, the request will be rejected.
// Otherwise, the request will pass.
type Limiter interface {
	Limit() bool
}

// UnaryServerInterceptor returns a new unary server interceptors that performs request rate limiting.
func UnaryServerInterceptor(limiter Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if limiter.Limit() {
			return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected by grpc_ratelimit middleware, please retry later.", info.FullMethod)
		}
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new stream server interceptor that performs rate limiting on the request.
func StreamServerInterceptor(limiter Limiter) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if limiter.Limit() {
			return status.Errorf(codes.ResourceExhausted, "%s is rejected by grpc_ratelimit middleware, please retry later.", info.FullMethod)
		}
		return handler(srv, stream)
	}
}
`

	GRPC_RATE_LIMITER_TMP = `package grpc_mw

import (
	"go.uber.org/ratelimit"
)


type PassLimiter struct{
	limiter ratelimit.Limiter
}

func NewPassLimiter(rate int)*PassLimiter{
	pl:=PassLimiter{limiter: ratelimit.New(rate)}
	return &pl
}
func (this *PassLimiter) Limit() bool {
	this.limiter.Take()
	return false
}

`
)
