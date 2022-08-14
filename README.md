# gatherInfo

功能如下：
- [x] 子域名收集
- [ ] 目录扫描
- [ ] 端口扫描
- [ ] 指纹识别

### 依赖
`windows` 需要安装`WinPcap` [点这下载](https://www.winpcap.org/install/default.htm) 
`linux` 下需安装`libpcap-dev`
### 子域名模块

- [x] 分为API查找和子域名爆破
- [x] 支持自定API
- [x] 支持强制爆破泛解析域名
- [x] 支持自定义带宽
- [x] 支持自定义字典
- [x] 支持验证域名是否存活

参数：
- `-d` 后面跟主域名，如`go run ./main.go subdomain -d google.com`
- `-b` 设置带宽，默认值为30000,每秒大约300个DNS数据包
- `-p` payload爆破字典路径,默认为`dict`目录下`subdomain.txt` 
- `-v` 加上此参数会验证域名是否存活
- `-w` 加上此参数，会强制爆破泛解析域名。默认遇到泛解析会跳过
- `-m` 可用此参数指定模式,默认先查询API再进行爆破。`-m api` | `-m enu` 只使用API查询或只进行爆破
- `-o` 使用此选项，可指定输出文件名。文件输出在项目的`output`目录下，如过不存在会自动创建。输出文件应为`xlsx`格式

例子：`go run ./main.go subdomain -d google.com -m api -v -o test.xlsx`

### 目录扫描模块

- [ ] 自定义头部


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

