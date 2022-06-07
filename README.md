# gatherInfo

功能如下：
- [x] 子域名收集
- [ ] 目录扫描
- [ ] 端口扫描
- [ ] 指纹识别

### 子域名模块

- 分为API查找(需填写部分平台的key)和子域名爆破
- 支持强制爆破泛解析域名
- 支持自定义带宽
- 支持自定义字典
- 支持验证域名是否存活

- [ ] 备案号查询。先通过查询系统域名备案号，再通过备案号反查与备案号相关的域名 
- http://www.beianbeian.com
- http://icp.bugscaner.com

- [ ] SSL证书
- https://myssl.com/ssl.html
- https://www.chinassl.net/ssltools/ssl-checker.html

- [ ] 测试本地最大发包数

目录扫描模块

- [ ] 自定义头部


### 指纹识别

- [ ] 访问大量路径匹配特征
- [ ] 根据响应内容匹配特征，不需要大量请求


端口扫描

- [ ] random SrcAdress
- [ ] 分片传输
- [ ] 速度太快会报错
- [ ] 发送包时候，缓冲区复用

## NOTE
- 赶快学习context
- 原子操作

## Reference
[ksubdomain](https://github.com/knownsec/ksubdomain)

