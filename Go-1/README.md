# Go语言

## go语言开发环境
进入官网下载Go

VSCode编辑器/GoLand IDE(收费）
    
## go语言基础语法和标准库
    package: 
    
    import: "fmt" 为输入输出包
    
    func main() { }:
    
    命令行：
        go run .\helloworld.go
        go build 
        
    变量：
    if else：if后面没有小括号，必须要有大括号
    for：什么都不写代表死循环
         使用C++类似的循环
         同样不加括号
    switch：switch后面不加括号
            底下的case不需要家break，不会跳到其他分支
            
    数组：var a[5]int
    
    切片：动态数组
         s := make([]string, 3)
         s[0] = "a"
         s = append(s, "d")
         
    map: m := make(map[string]int) // string为key， int为val
         m["one"] = 1
         delete(m, "one")   
         
    range:范围索引
        for k, v := range m {
            fmt.Println(k, v)
        }
        
    函数：func add(a int, b int) int {
             return a + b
         }
         数据类型通常是后置的
         支持返回多个值，第一个为正真返回的值，后面的相当于错误信息等
         如果希望引用实参，形参加星号*表示指针，实参加取址&，调用时加星号
         
    结构体：
        type user struct {
            name string
            password string
        }
        a := user{name: "wang", password: "1024"}
        var d user
        d.name = "wang"
        d.password = "1024"
        
    结构体方法：
        func (u *user) resetPassword(password string) {
            u.password = password
        }
        使用变量加点调用函数
        
    错误处理：
        import "errors"
        在参数后面加上 err error
        如果函数返回有error，则要先处理错误
        
    字符串操作：
        import "strings"
        
    字符串格式化：
        import "fmt"
        fmt.Printf("s=%v\n", s)
        fmt.Printf("s=%+v\n", p)
        fmt.Printf("s=%#v\n", p)
        fmt.Printf("%.2f\n", f)
        
    JSON处理：
        import "encoding/json"
        
    时间处理：
        import "time"
        now := time.Now()
        time.Unix() //获取时间戳
        
    数字解析：
        import "strconv"
        strconv.ParseFloat("1.234", 64)
        
    进程信息：
        import "os"
        import "os/exec"
        os.Args
        os.Getenv("PATH")
        os.Setenv("AA", "BB")
        exec.Command("grep", "127.0.0.1", "/etc/hosts").CombinedOutput


## 猜谜游戏
    在Windows命令行下最后运行会一直跳出错误输入
    原因是
        input = strings.TrimSuffix(input, "\r\n")
        最后需要以\r\n结尾，否则会多出来一个回车符

## 命令行词典
    抓包
    复制cURL
    转换成Go代码： https://curlconverter.com/#go
    运行得到JSON
    将JSON转成Go代码（结构体）： https://oktools.net/json2go
    在程序中修改输入输出

## SOCKS5代理
    SOCKS5代理服务器，某些企业的内网为了安全性，可能配备了很严格的防火墙策略，所以管理员访问资源可能很麻烦。SOCKS5相当于在防火墙内部开了个口子，可以通过单个端口访问所有资源。
    curl --socks5 127.0.0.1:1080 -v http://www.qq.com
    -v 表示打印出所有细节

### SOCKS5代理原理

![image.png](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/aa8f887c2f794b0d9627a54254694d96~tplv-k3u1fbpfcp-watermark.image?)

TCP echo erver

    go关键字，协程处理
    defer关键字，函数内生命周期结束执行

解析请求
    
    将客户端发过来的请求进行解析
    // +----+----------+----------+
    // |VER | NMETHODS | METHODS  |
    // +----+----------+----------+
    // | 1  |    1     | 1 to 255 |
    // +----+----------+----------+
    // VER: 协议版本，socks5为0x05
    // NMETHODS: 支持认证的方法数量
    // METHODS: 对应NMETHODS，NMETHODS的值为多少，METHODS就有多少个字节。RFC预定义了一些值的含义，内容如下:
    // X’00’ NO AUTHENTICATION REQUIRED
    // X’02’ USERNAME/PASSWORD

与服务端建立链接

    读取请求的数据
    // +----+-----+-------+------+----------+----------+
    // |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
    // +----+-----+-------+------+----------+----------+
    // | 1  |  1  | X'00' |  1   | Variable |    2     |
    // +----+-----+-------+------+----------+----------+
    // VER 版本号，socks5的值为0x05
    // CMD 0x01表示CONNECT请求
    // RSV 保留字段，值为0x00
    // ATYP 目标地址类型，DST.ADDR的数据对应这个字段的类型。
    //   0x01表示IPv4地址，DST.ADDR为4个字节
    //   0x03表示域名，DST.ADDR是一个可变长度的域名
    // DST.ADDR 一个可变长度的值
    // DST.PORT 目标端口，固定2个字节
    
    将请求转发给服务端
    dest, err := net.Dial("tcp", fmt.Sprintf("%v:%v", addr, port))

将服务端发回来的转发给客户端，将客户端发过去的转发给服务端

    _, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
    if err != nil {
        return fmt.Errorf("write failed: %w", err)
    }
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // 两个协程分别处理客户端到服务端的转发和服务端到客户端的转发
    go func() {
        _, _ = io.Copy(dest, reader)
        cancel()
    }()
    go func() {
        _, _ = io.Copy(conn, dest)
        cancel()
    }()
    
    // 等待ctx的结束，否则会直接return
    <-ctx.Done()
    return nil

## 课后作业

![image.png](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/4b2ecd1e587540c096a919f122180edb~tplv-k3u1fbpfcp-watermark.image?)

### 第一题
使用如下两行代码替代原来的输入

```go
var guess int
fmt.Scanf("%d\n", &guess)
```

### 第二题

增加了使用火山翻译引擎的支持


```go
type HuoShanDiction struct {
	Words []struct {
		Source  int    `json:"source"`
		Text    string `json:"text"`
		PosList []struct {
			Type      int `json:"type"`
			Phonetics []struct {
				Type int    `json:"type"`
				Text string `json:"text"`
			} `json:"phonetics"`
			Explanations []struct {
				Text     string `json:"text"`
				Examples []struct {
					Type      int `json:"type"`
					Sentences []struct {
						Text      string `json:"text"`
						TransText string `json:"trans_text"`
					} `json:"sentences"`
				} `json:"examples"`
				Synonyms []interface{} `json:"synonyms"`
			} `json:"explanations"`
			Relevancys []interface{} `json:"relevancys"`
		} `json:"pos_list"`
	} `json:"words"`
	Phrases  []interface{} `json:"phrases"`
	BaseResp struct {
		StatusCode    int    `json:"status_code"`
		StatusMessage string `json:"status_message"`
	} `json:"base_resp"`
}

type HuoShanDictRequest struct {
	Text     string `json:"text"`
	Language string `json:"language"`
}

func huoshanQuery(word string, waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	client := &http.Client{}
	request := HuoShanDictRequest{Text: word, Language: "en"}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	// var data = strings.NewReader(`{"text":"good\n","language":"en"}`)
	req, err := http.NewRequest("POST", "https://translate.volcengine.com/web/dict/match/v1/?msToken=&X-Bogus=DFSzswVLQDVBKiQrSWQR4Pt/pLv9&_signature=_02B4Z6wo00001yjd.7gAAIDCo5ZkWs8NEw8o3fsAAKhLNQOsQfcxG2TsBgxwl7OK8tsxiuNz4KNSUdOHH.FY7VnK2B-b7XGr1F-h2H9gFcyNMdkZ5K46sl8nevCpsimTy5CyHUvb73Mc5YpY95", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "translate.volcengine.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", "x-jupiter-uuid=16519240543857514; i18next=zh-CN; ttcid=4dd1d6481a27405d8e088f7a6ab7d51a66; tt_scid=jlnel4RAjPdPt4fMI2z3xvfsjix73mk6PiqxOzXNcJE9sACpNxnSP4OjtY69fQQ833e2; s_v_web_id=verify_0b043b2145bfe9fab707ee95bf91ebc6; _tea_utm_cache_2018=undefined")
	req.Header.Set("origin", "https://translate.volcengine.com")
	req.Header.Set("referer", "https://translate.volcengine.com/translate?category=&home_language=zh&source_language=detect&target_language=zh&text=good%0A")
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="101", "Google Chrome";v="101"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var huoshanDiction HuoShanDiction
	err = json.Unmarshal(bodyText, &huoshanDiction)
	if err != nil {
		log.Fatal(err)
	}

	l.Lock()
	if isOutput == 0 {
		fmt.Println(word, "\nUK:", huoshanDiction.Words[0].PosList[0].Phonetics[0].Text,
			"US:", huoshanDiction.Words[0].PosList[0].Phonetics[1].Text)
		for _, items := range huoshanDiction.Words[0].PosList {
			for _, exps := range items.Explanations {
				fmt.Println(exps.Text)
			}
		}
		isOutput = 1
		fmt.Println("-------------------------------Query by 火山翻译")
	}
	l.Unlock()
	// fmt.Printf("%#v\n", huoshanDiction)
}
```

### 第三题

使用两个协程并行运行两种搜索引擎的查询，使用waitgroup让两个协程都结束了再退出程序

```go
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage: simpleDict WORD
example: simpleDict hello
		`)
		os.Exit(1)
	}
	word := os.Args[1]
	waitgroup := sync.WaitGroup{}
	waitgroup.Add(2)
	go huoshanQuery(word, &waitgroup)
	go query(word, &waitgroup)
	waitgroup.Wait()
}
```
为了最后只会输出一种查询结果，在两个查询过程的输出部分添加一个互斥量，先查询到结果的协程会获得互斥锁，根据isOutput来判断是否输出，输出完后将isOutput置为1，再放弃锁。之后第二个获得锁的协程就不会继续输出了。

```go
var l sync.Mutex
var isOutput = 0
func huoshanQuery(word string, waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	client := &http.Client{}
	request := HuoShanDictRequest{Text: word, Language: "en"}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	// var data = strings.NewReader(`{"text":"good\n","language":"en"}`)
	req, err := http.NewRequest("POST", "https://translate.volcengine.com/web/dict/match/v1/?msToken=&X-Bogus=DFSzswVLQDVBKiQrSWQR4Pt/pLv9&_signature=_02B4Z6wo00001yjd.7gAAIDCo5ZkWs8NEw8o3fsAAKhLNQOsQfcxG2TsBgxwl7OK8tsxiuNz4KNSUdOHH.FY7VnK2B-b7XGr1F-h2H9gFcyNMdkZ5K46sl8nevCpsimTy5CyHUvb73Mc5YpY95", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "translate.volcengine.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", "x-jupiter-uuid=16519240543857514; i18next=zh-CN; ttcid=4dd1d6481a27405d8e088f7a6ab7d51a66; tt_scid=jlnel4RAjPdPt4fMI2z3xvfsjix73mk6PiqxOzXNcJE9sACpNxnSP4OjtY69fQQ833e2; s_v_web_id=verify_0b043b2145bfe9fab707ee95bf91ebc6; _tea_utm_cache_2018=undefined")
	req.Header.Set("origin", "https://translate.volcengine.com")
	req.Header.Set("referer", "https://translate.volcengine.com/translate?category=&home_language=zh&source_language=detect&target_language=zh&text=good%0A")
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="101", "Google Chrome";v="101"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var huoshanDiction HuoShanDiction
	err = json.Unmarshal(bodyText, &huoshanDiction)
	if err != nil {
		log.Fatal(err)
	}

	l.Lock()
	if isOutput == 0 {
		fmt.Println(word, "\nUK:", huoshanDiction.Words[0].PosList[0].Phonetics[0].Text,
			"US:", huoshanDiction.Words[0].PosList[0].Phonetics[1].Text)
		for _, items := range huoshanDiction.Words[0].PosList {
			for _, exps := range items.Explanations {
				fmt.Println(exps.Text)
			}
		}
		isOutput = 1
		fmt.Println("-------------------------------Query by 火山翻译")
	}
	l.Unlock()
	// fmt.Printf("%#v\n", huoshanDiction)
}

func query(word string, waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	client := &http.Client{}
	request := DictRequest{TransType: "en2zh", Source: word}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("Referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")
	req.Header.Set("X-Authorization", "token:qgemv4jr1y38jyq6vhvi")
	req.Header.Set("app-name", "xy")
	req.Header.Set("os-type", "web")
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="101", "Google Chrome";v="101"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var autoGenerated AutoGenerated
	err = json.Unmarshal(bodyText, &autoGenerated)
	if err != nil {
		log.Fatal(err)
	}
	l.Lock()
	if isOutput == 0 {
		fmt.Println(word, "UK:", autoGenerated.Dictionary.Prons.En, "US:", autoGenerated.Dictionary.Prons.EnUs)
		for _, item := range autoGenerated.Dictionary.Explanations {
			fmt.Println(item)
		}
		isOutput = 1
		fmt.Println("-------------------------------Query by 彩云翻译")
	}
	l.Unlock()
	//fmt.Printf("%s\n", bodyText)
}
```