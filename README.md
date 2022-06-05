# gatherInfo

integrate 'collect subdomain','directory scan','port scan',and 'fingerprint scan'

子域名模块

- 子域名爆破，支持强制爆破泛解析域名，比较粗糙，搞个黑名单直接丢弃,不能只把A记录弄进黑名单，其他记录也要,因此要解析数据包
- 状态表用以记录每个数据包的状态，超时则会重发，重发两次,尽可能解决丢包问题。由于每个数据包都会存在状态表,因此会牺牲内存，用内存换准确性。
- [ ] 备案号查询。先通过查询系统域名备案号，再通过备案号反查与备案号相关的域名 
- http://www.beianbeian.com
- http://icp.bugscaner.com

- [ ] SSL证书
- https://myssl.com/ssl.html
- https://www.chinassl.net/ssltools/ssl-checker.html

- [ ] 测试本地最大发包数

目录扫描模块

- [ ] 自定义头部


指纹识别
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
- 链表
- 问题：目前使用了双向循环链表优化，基本完成
- bug:有时候会查询不到表(查询出错，或者删早le) 初步判断在checktimeout出错


