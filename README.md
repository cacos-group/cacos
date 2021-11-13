# cacos-server
cacos server

## 设计文档

[设计文档地址](https://mbd.baidu.com/newspage/data/landingshare?preview=1&pageType=1&isBdboxFrom=1&context=%7B%22nid%22%3A%22news_9418462740030256263%22%2C%22sourceFrom%22%3A%22bjh%22%7D)

## Cacos是什么？

Cacos是配置管理中心，使用etcd、mysql存储，go语言开发。

## Cacos架构

### 架构图

![jiatou](https://pic.rmb.bdstatic.com/bjh/news/2d19255654bdd5234cf81a048b6d3a78.jpeg@wm_2,t_55m+5a625Y+3L2FscGhh,fc_ffffff,ff_U2ltSGVp,sz_27,x_17,y_17)

### 数据模型

![shujumoxing](https://pic.rmb.bdstatic.com/bjh/news/bc43bca30bdbaa04fa6cc0deed62de74.jpeg@wm_2,t_55m+5a625Y+3L2FscGhh,fc_ffffff,ff_U2ltSGVp,sz_32,x_20,y_20)

### 数据一致性

![eventSourcing](https://f11.baidu.com/it/u=2183350210,158022252&fm=30&app=106&f=JPEG&access=215967316?w=640&h=180&s=18A05D30857644225ECD65DE000080B2)

## 业务逻辑
![yewulouji](https://pic.rmb.bdstatic.com/bjh/news/3a46c36a62372aedd8f044741df4e0f4.png@wm_2,t_55m+5a625Y+3L2FscGhh,fc_ffffff,ff_U2ltSGVp,sz_28,x_18,y_18)

## 设计方案
![tech](https://pic.rmb.bdstatic.com/bjh/news/bd417b0ac98f9bc89ed40ba15ff356a2.png@wm_2,t_55m+5a625Y+3L2FscGhh,fc_ffffff,ff_U2ltSGVp,sz_28,x_18,y_18)