# zaplogger

`zaplogger` 是一个生产可用的日志包，基于 `zap` 包封装。具有如下特性：

- 支持日志级别：`Debug`、`Info`、`Warn`、`Error`、`Panic`、`Fatal`。默认Info级别
- 支持自定义配置。
- 支持Caller文件名和行号。
- 支持输出到标准输出和文件，可以同时输出到多个对象。
- 支持 `JSON` 和 `Console` 两种日志格式。
- 支持结构化日志记录。
- **支持Context（业务定制）**
- 支持动态开关Debug日志



## 功能特性

- WithFields 直接修改全局日志结构
  - 添加一些默认通用的字段到每行日志，方便日志查询和分析。
- WithValues 修改logger字段，返回一个新的logger对象，应用场景为
  - 支持RequestID：使用RequestID串联一次请求的所有日志，这些日志可能分布在不同的组件，不同的机器上。支持RequestID可以大大提高排障的效率，降低排障难度。在一些大型分布式系统中，没有RequestID排障简直就是灾难。
- WithContext 可以通过context传递子logger进行调用



## 快速开始

### 简单示例

开箱即用，代码使用 `logger` 包默认的全局 `SLogger`，

```go
func main() {
  defer slog.Sync()
  // 制定日志级别，默认为info级别
  logger.SetOptions(logger.WithLevel("debug"))
  logger.SLogger.Debug("debug")
}
```



### 自定义初始化配置

可以使用 `New` 来初始化一个日志包，如下：

```go
func main() {

	opt := &logger.Options{
		Output:        os.Stdout,				// 输出对象
		Level:         "debug",	// 日志级别
		DisableCaller: true,						// 是否开启调用信息
	}

	newSlog := logger.Init(opt)
	defer newSlog.Sync()

	newSlog.Info("test")
}
```



### 结构化日志输出

**结构化**

```go
package main

import (
	"github.com/yuanbaopig/logger"
	"go.uber.org/zap"
)

func main() {

	logger.SetOptions(
		logger.WithFields(zap.String("app", "myApp"), zap.Int("version", 1)), // 日志结构化

	)
	lg := zap.L().Sugar()
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



## 使用指南

logger配置修改

修改配置属性查看with相关函数

```go
// 修改日志级别
logger.SetOptions(logger.WithLevel("debug"))

// 修改全局对象字段
logger.SetOptions(logger.WithFields(zap.Int("userID", 10), zap.String("requestID", "fbf54504")))

// 开启调用信息，默认关闭
logger.SetOptions(logger.WithDisableCaller(true))
```

推荐配置

```go
logger.SetOptions(logger.WithLevel("debug"), logger.WithDisableCaller(true))
```

制定输出对象

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
	defer logger.SLogger.Sync()

	// 定义字段
	lv := logger.WithValues(zap.Int("userID", 10))

	lv.Info("test")

	// 讲logger存储到context中
	ctx := logger.WithContext(context.Background(), lv)

	// 进行context传递
	PrintString(ctx, "World")

	// 原结构不受影响
	logger.SLogger.Sugar().Infof("Hello World")
}

func PrintString(ctx context.Context, str string) {
	//从context中获取logger
	lc := logger.FromContext(ctx)
	lc.Sugar().Infof("Hello %s", str)
}
```



