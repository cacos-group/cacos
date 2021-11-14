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

![tech](https://f11.baidu.com/it/u=352677692,158135625&fm=30&app=106&f=JPEG&access=215967316?w=640&h=316&s=8451427C07725C2048C409580200C0F2)

### 设计方案
![design](https://f12.baidu.com/it/u=2976030758,158135616&fm=30&app=106&f=JPEG&access=215967316?w=640&h=969&s=59A81D72190B504F1EF560CA0000E0B2)
