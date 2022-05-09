## 并发并行

并发：多线程程序在一个核的CPU上运行

并行：多线程程序在多个核的CPU上的运行

Go可以充分发挥多核优势，高效运行

### Goroutine

协程：用户态，轻量级线程，栈KB级别

线程：内核态，线程跑多个协程，栈MB级别

```go
func hello(i int) {
    println("hello goroutine : " + fmt.Sprint(i))
}
func HelloGoRoutine() {
    for i:= 0; i < 5; i++ {
        go func(j int) {
            hello(j)
        }(i)
    }
    time.Sleep(time.Second)
}
```

### CSP(Communicating Sequential Processes)

通过通信共享内存，在Go语言中使用通道实现

Go提倡通过通信共享内存而不是通过共享内存而实现通信

### Channel

make(chan 元素类型,[缓冲大小])

- 无缓冲通道 make(chan int)

- 有缓冲通道 make(chan int, 2)

```go
// A 子协程发送0~9数字
// B 子协程计算输入数字的平方
// 主协程输出最后的平方数
func CalSquare() {
    // 无缓冲通道
    src := make(chan int)
    // 有缓冲通道
    dest := make(chan int, 3)
    go func() {
        defer close(src)
        for i := 0; i < 10; i++ {
            src <- i
        }
    }()
    go func() {
        defer close(dest)
        for i := range src {
            dest <- i*i
        }
    }()
    for i := range dest {
        // 复杂操作
        println(i)
    }
}
```

### 并发安全 Lock

```go
lock sync.Mutex
lock.Lock()
lock.Unlock()
```

### WaitGroup

```go
var wg sync.WaitGroup
wg.Add(delta int)		// 计数器+delta
wg.Done()				// 计数器减一
wg.Wait()				// 阻塞直到计数器为0
```

## 依赖管理

Go依赖管理演进

- GOPATH
- Go Vendor
- Go Module

#### GOPATH

环境变量 $GOPATH

```
.
|--bin			// 项目编译的二进制文件
|--pkg			// 项目编译的中间产物，加速编译
|--src			// 项目源码
```

项目代码直接依赖src下的代码

go get 下载最新版本的包到src目录下

弊端：如果项目依赖于某一个包的不同版本，则会出现问题。无法实现package的多版本控制

#### Go Vendor

- 项目目录下增加vendor文件，所有依赖包副本形式放在vendor中
- 依赖寻址方式：vendor=>GOPATH

通过每个项目引入一份依赖的副本，解决了多个项目需要同一个package依赖的冲突问题

问题：无法控制依赖的版本，更新项目又可能出现依赖冲突

#### Go Module

- 通过go.mod文件管理依赖包版本
- 通过go get/go mod指令工具管理依赖包

### 依赖管理的三要素

配置文件，描述依赖 go.mod

中心仓库管理依赖库 proxy

本地工具 go get/mod

#### 依赖配置-go.mod

依赖标识： \[Module Path][Version/Pseudo-version]

```go
module example/project/app		依赖管理基本单元

go 1.16							原生库

require (						单元依赖
	example/lib1 v1.0.2
    example/lib2 v1.0.0 // indirect
)
```

#### 依赖配置-version

语义化版本

${MAJOR}.${MINOR}.${PATCH}

v1.3.0

基于commit伪版本

vX.0.0-yyyymmddhhmmss-abcdefgh123

v0.0.0-20220401081311-c38fb59326b7

#### 依赖配置-indirect

使用indirect标识表示非直接依赖

主版本2+模块会在模块路径后加/vN

对于没有go.mod文件并且主版本2+的依赖，会+incompatible

#### 依赖配置-依赖图

会选择最低的兼容版本

#### 依赖分发-回源 Proxy

例如可以直接依赖Github、SVN等

但是这样会无法保证构建稳定性，无法保证依赖可用性，增加第三方压力

使用Proxy会缓存从Github、SVN中的依赖版本

#### 依赖分发-变量 GOPROXY

GOPROXY="https://proxy1.cn,https://proxy2.cn,direct"

服务站点URL列表，"direct"表示源站

#### 工具-go get/mod

```go
go get exapmle.org/pkg  	@update  	// 默认
							@none		// 删除依赖
							@v1.1.2		// tag版本，语义版本
							@23dfdd5	// 特定的commit
							@master		// 分支的最新commit
                            
go mod  init 			// 初始化，创建go.mod文件
		download		// 下载模块到本地缓存
		tidy			// 增加需要的依赖，删除不需要的依赖
```

## 测试

- 回归测试
- 集成测试
- 单元测试

从上到下，覆盖率逐层变大，成本却逐层降低

### 单元测试

包括输入、测试单元、输出’校对

#### 单元测试-规则

- 所有测试文件以_test.go结尾

- func TestXxx(*testing.T)

```go
func TestPublishPost(t *testing.T) {}
```

- 初始化逻辑放到TestMain中

```go
func TestMain(m *testing.M) {
	// 测试前：数据装载、配置初始化等前置工作
	code := m.Run()
	// 测试后：释放资源等收尾工作
	os.Exit(code)
}
```

```go
func HelloTom() string {
	return "Jerry"
}

func TestHelloTom(t *testing.T) {
    output := HelloTom()
    expectOutput := "Tom"
    if output != expectOutput {
        t.Errorf("Expected %s do not match actual %s", expectOutput, output)
    }
}
```

#### 单元测试-运行

```
go test [flags][packages]
RUN TestHelloTom
```

#### 单元测试-assert

```go
import "github.com/stretchr/testify/assert"
assert.Equal(t, expectOutput, output)
```

#### 单元测试-覆盖率

代码覆盖率

```
go test 			--cover
```

- 一般覆盖率：50%~60%，较高覆盖率80%+
- 测试分支相互独立、全面覆盖
- 测试单元粒度足够小，函数单一职责

#### 单元测试-依赖

外部依赖 => 稳定&幂等

#### 单元测试-Mock

monkey：https://github.com/bouk/monkey

快读Mock函数

- 为一个函数打桩
- 为一个方法打桩

### 基准测试

```go
func BenchmarkXxx(b *testing.B) { }
func BenchmarkXxxParallel(b *testing.B) { }
```

## 项目实践

### 需求背景

社区话题页面

- 展示话题（标题，文字描述）和回帖列表
- 暂不考虑前端页面实现，仅仅实现一个本地web服务
- 话题和回帖数据用文件存储

### ER图

Topic id title ...

Post id topic_id ...

### 分层结构

- 数据层：数据Model，外部数据的增删改查
- 逻辑层：业务Entity，处理核心业务逻辑输出
- 视图层：试图view，处理和外部的交互逻辑

```
		Model		Entity		View
File->Repository -> Service -> Controller->Client
		数据层			逻辑层		视图层
```

### 组件工具

Gin高性能go web框架

https://github.com/gin-gonic/gin#installation

Go Mod

go mod init

go get gopkg.in/gin-gonic/gin.v1@v1.3.0 (建议使用github上最新版的gin，否则可能会出错)

## 社区话题页面需求描述

* 展示话题（标题，文字描述）和回帖列表
* 暂时不考虑前端页面实现，仅仅实现一个本地web服务
* 话题和回帖数据用文件存储

### data包

元数据文件

含有两个文件topic和post，使用json格式存储

### repository包

数据层

将data包中文件的内容读取成结构体存储在map中

### service包

服务层

根据给定的topicId来获取所在的topic和posts列表

### controller包

视图层

将数据暴露到接口中

### go.mod文件

go module 依赖配置管理

### server.go

web服务main入口文件
