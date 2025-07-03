# zaplogger

`zaplogger` 是一个生产可用的日志包，基于 `zap` 包封装。具有如下特性：

- 支持日志级别：`Debug`、`Info`、`Warn`、`Error`、`Panic`、`Fatal`。默认Info级别
- 支持Caller文件名和行号。
- 支持输出到标准输出和文件，可以同时输出到多个对象。
- 支持error级别以上日志输出到指定对象
- 支持 `JSON` 和 `Console` 两种日志格式。
- 支持结构化日志记录。
- **支持Context传递（业务定制）**



## 功能特性

- WithFields 直接修改全局日志结构
  - 添加一些默认通用的字段到每行日志，方便日志查询和分析。
- WithValues 修改logger字段，返回一个新的logger对象
  - 应用场景例如RequestID追踪，使用RequestID串联一次请求的所有日志，这些日志可能分布在不同的组件，不同的机器上。支持RequestID可以大大提高排障的效率，降低排障难度。在一些大型分布式系统中，没有RequestID排障简直就是灾难。
- WithContext 可以通过context传递子logger进行调用



## 快速开始

### 简单示例

开箱即用，代码使用 `logger` 包默认的全局 `SLogger`，

```go
	// 调用包里的logger
	logger.Log.Info("info")
	defer logger.Log.Sync()

	// 调用SugaredLogger
	lg := logger.Log.Sugar()
	lg.Debug("debug")
	lg.Warn("warning")
	lg.Error("error")
	lg.Fatal("fatal")
```



### 自定义初始化配置

可以使用 `New` 来初始化一个日志包，如下：

```go
	lg := logger.New(
		logger.WithFormat("console"),                                         // 日志格式
		logger.WithEnableColor(true),                                         // 是否显示日志级别颜色，只有在console下才支持
		logger.WithDisableCaller(false),                                      // 是否关闭调用信息，默认为false
		logger.WithLevel("info"),                                             // 日志级别
		logger.WithDisableStacktrace(false),                                  // 是否关闭Stacktrace，默认false
		logger.WithFields(zap.String("app", "myApp"), zap.Int("version", 1)), // 日志结构化
		//logger.WithOutputPaths([]string{"test.log"}),
		logger.WithOutputPaths([]string{"stdout"}),
		//logger.WithErrorOutputPaths([]string{"stderr", "error.log"}), // 如果使用同时使用stdout和stderr，并且没有屏蔽stdout的情况下会导致error信息输出两次
	)

	lg.Info("test")
	lg.Error("error")
```



配置变更

```go
	lg := logger.New(
		logger.WithFields(zap.String("k", "v")),
		logger.WithLevel("info"),
	)
	lg.Debug("test")

	lg.SetOptions(logger.WithLevel("debug"))

	lg.Debug("test1")
```



#### 配置说明

**logger配置修改**

修改配置属性查看with相关函数

```go
// 修改日志级别
logger.SetOptions(logger.WithLevel("debug"))

// 修改全局对象字段
logger.SetOptions(logger.WithFields(zap.Int("userID", 10), zap.String("requestID", "fbf54504")))

// 开启调用信息，默认开启
logger.SetOptions(logger.WithDisableCaller(true))
```

指定输出对象

```shell
logger.SetOptions(logger.WithOutput(logger.GetFileLogWriter("test.log")))
```

修改输出模式

```go
// 默认为json格式，并且标准输出为的输出模式不可变更
logger.SetOptions(logger.WithFormat("console"))
```

Format 支持 `console` 和 `json` 2 种格式：

- console：输出为 text 格式。例如：`2020-12-05 08:12:02.324	DEBUG	example/example.go:43	This is a debug message`
- json：输出为 json 格式，例如：`{"level":"debug","time":"2020-12-05 08:12:54.113","caller":"example/example.go:43","msg":"This is a debug message"}`

> 标准输出的日志格式不可变更

函数调用链

```shell
logger.SetOptions(logger.WithDisableStacktrace(false))

lg.Error("error")
```

只有error级别以上的日志才会输出trace信息

```shell
2025-07-03T17:07:15.128+0800	ERROR	options/options_test.go:23	error	{"app": "myApp", "version": 1}
github.com/yuanbaopig/logger/example/options.TestOptions
	/Users/mj/Documents/Project_go/logger/example/options/options_test.go:23
testing.tRunner
	/usr/local/go/src/testing/testing.go:1690
--- PASS: TestOptions (0.00s)
```





### 结构化日志输出

**永久结构化**

```go
package main

import (
	"github.com/yuanbaopig/logger"
	"go.uber.org/zap"
)

func main() {

  lg := logger.New(
		logger.WithFields(zap.String("app", "myApp"), zap.Int("version", 1)), // 日志结构化
	)
  
	lg.Debug("debug")
	lg.Info("info")
	lg.Warn("warning")
	lg.Error("error")
}
```



**临时结构化**

`marmotedu/log` 也支持结构化日志打印，例如：

```go
// logger 的使用方式
log.Info("This is a info message", log.Int32("int_key", 10))
// sugarlogger 的使用方式
log.Infow("Message printed with Errorw", "X-Request-ID", "fbf54504-64da-4088-9b86-67824a7fb508") 
```

对应的输出结果为：

```shell
2020-12-05 08:16:18.749	INFO	example/example.go:44	This is a info message	{"int_key": 10}
2020-12-05 08:16:18.749	ERROR	example/example.go:46	Message printed with Errorw	{"X-Request-ID": "fbf54504-64da-4088-9b86-67824a7fb508"}
```

`log.Info` 这类函数需要指定具体的类型，以最大化的提高日志的性能。`log.Infow` 这类函数，不用指定具体的类型，底层使用了反射，性能会差些。建议用在低频调用的函数中。





### Context 传递logger

修改logger字段，返回一个新的logger对象

```go
package main

import (
	"context"
	"github.com/yuanbaopig/logger"
	"go.uber.org/zap"
)

func main() {
	defer logger.Log.Sync()

	// 定义字段
	lv := logger.WithValues(zap.Int("userID", 10))

	lv.Info("test")

	// 将logger存储到context中
	ctx := logger.WithContext(context.Background(), lv)

	// 进行context传递
	PrintString(ctx, "World")

	// 原结构不受影响
	logger.Log.Sugar().Infof("Hello World")
}

func PrintString(ctx context.Context, str string) {
	// 从context中获取logger
	lc := logger.FromContext(ctx)
	lc.Sugar().Infof("Hello %s", str)
}
```



