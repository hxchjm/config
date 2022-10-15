#配置
支持2中配置方式，两者配置优先选择本地配置：
#### 1. 本地配置文件
- 命令行参数 -conf `./main -conf x/y/test.json`
- 环境变量 CONF_PATH `export CONF_PATH=x/y/test.json`

#### 2. 使用远程配置平台（nacos）
- 命令行参数 `-nacos.host=localhost:8849 -nacos.namespaceid=public -nacos.group=public`
- 环境变量 
    ```shell script
    export NACOS_HOST=localhost:8849 
    export NACOS_NAMESPACEID=public 
    export NACOS_GROUP=public

    ```
  
#### 3.例子
```go
package main

import (
	"config"
	"flag"
	"log"
)

type Test struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	flag.Parse()
	if err := config.Init(); err != nil {
		panic(err)
	}
	var test Test
	config.Get("Test", &test)
	log.Printf("test is %+v", test)
}
```