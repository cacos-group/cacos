# cacos-server
cacos server

## 设计文档

[设计文档地址](https://mbd.baidu.com/newspage/data/landingshare?preview=1&pageType=1&isBdboxFrom=1&context=%7B%22nid%22%3A%22news_9477430390378252781%22%2C%22sourceFrom%22%3A%22bjh%22%7D)

## Cacos是什么？

Cacos是配置管理中心

* 1.mysql--存储事件日志
* 2.etcd--存储和查询kv；访问权限控制；
* 3.golang

## Cacos架构

### 数据模型和架构图

![tech](https://raw.githubusercontent.com/cacos-group/cacos/main/doc/sjmxjgt.jpg)

### 设计方案
![design](https://raw.githubusercontent.com/cacos-group/cacos/main/doc/design.png)