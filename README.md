# gatherInfo

功能如下：
- [x] 子域名收集
- [x] 目录扫描
- [ ] 端口扫描
- [ ] 指纹识别

### 子域名模块

#### 依赖
`windows` 需要安装`winpcap` [点这下载](https://www.winpcap.org/install/default.htm) 
`linux` 下需安装`libpcap-dev`

- [x] 分为api查找和子域名爆破
- [x] 支持自定api
- [x] 支持强制爆破泛解析域名
- [x] 支持自定义带宽
- [x] 支持自定义字典
- [x] 支持验证域名是否存活

参数：
- `-d` 后面跟主域名，如`go run ./main.go subdomain -d google.com`
- `-b` 设置带宽，默认值为30000,每秒大约300个dns数据包
- `-p` payload爆破字典路径,默认为`dict`目录下`subdomain.txt` 
- `-v` 加上此参数会验证域名是否存活
- `-w` 加上此参数，会强制爆破泛解析域名。默认遇到泛解析会跳过
- `-m` 可用此参数指定模式,默认先查询api再进行爆破。`-m api` | `-m enu` 只使用api查询或只进行爆破
- `-o` 使用此选项，可指定输出文件名。文件输出在项目的`output`目录下，如过不存在会自动创建。输出文件应为`xlsx`格式

例子：`go run ./main.go subdomain -d example.com -m api -v -o test.xlsx`

### 目录扫描模块

- [x] 自定义请求头部
- [x] 自定义请求方法(不支持携带请求体,目前也没碰到过这种需求)
- [x] 自定义并发量
- [x] 自定义错误，使用部分响应内容即可。多个错误用逗号分开
- [x] 自定义合法状态码
- [x] 自定义字典

参数：

- `-c` 或 `--codes ints`         默认所有状态码都会打印，使用此参数则只会列出指定状态码的结果，请使用逗号 `,` 分开
- `-d` 或 `--dictionary string`  payload 字典路径 (默认 "dict/dir.txt")
- `-f` 或 `--filter string`      自定义错误，过滤字符串,使用部分相应内容即可，多个错误用逗号分开
- `-H` 或 `--header string`      自定义http头部，`key` 和 `value` 之前使用 `:` 分开例如: -H "Authorization: sercretxxxxxx"
- `-m` 或 `--method string`      自定义请求方法，一般就用`HEAD`和`GET` (默认 "GET")
- `-p` 或 `--proxy string`       定义代理，例如 `-p http://127.0.0.1:8888`
- `-t` 或 `--thread int`         并发量 (默认 30)
- `-U` 或 `--urldict string`     URL字典路径,每行一个URL
- `-u` 或 `--urls strings`       URL,支持多个URL，用 `,` 分开
- `-o` 或 `--output string`       结果输出文件，请使用`xlsx`后缀


例如: `go run .\main.go dirscan -u https://www.example.com -m HEAD -p http://127.0.0.1:8080 -H "Authorization: xxxxxxx" -t 20 -f "not found" -c 200,301,302 -o output.xlsx`

### 指纹识别

- [ ] 访问大量路径匹配特征
- [ ] 根据响应内容匹配特征，不需要大量请求


### 端口扫描

- [ ] random SrcAdress
- [ ] 分片传输
- [ ] 速度太快会报错
- [ ] 发送包时候，缓冲区复用

### 其他
- cmd颜色不支持，windows 平台请使用windows terminal 亦或linux (非常好用)



## NOTE
- API remove duplicate
- progressbar
- removedTabChan

## Reference
- [ksubdomain](https://github.com/knownsec/ksubdomain)

