**CloudFolder**

云文件系统

可以实现用户文件的上传下载以及分享，支持秒传功能。
使用了go语言的**Gin**框架实现。


使用了阿里云的oss，使用生产者消费者模式添加了异步转储文件到oss的功能（通过channel实现），在消费消息时使用channel控制并发数。


数据库：**redis,mongoDB**

mongoDB用来存储用户信息以及文件元信息，Redis保存用户token实现单用户登录。

安装：
git clone https://github.com/guozhao-coder/CloudFolder.git


运行：

1,把config包中的database文件里的数据库配置改为自己的配置。

2,MySQL表结构很简单，两个表（现在没有表，全部迁移到了Mongo）：

~~~~
CREATE TABLE `file` (
  `fileId` varchar(50) NOT NULL,
  `fileName` varchar(500) NOT NULL,
  `fileSize` float NOT NULL,
  `filePath` varchar(5000) NOT NULL,
  `userId` varchar(200) NOT NULL,
  `fileTime` varchar(40) NOT NULL,
  `fileHash` varchar(40) NOT NULL,
  `ossPath` varchar(500) NOT NULL,
  PRIMARY KEY (`fileId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

CREATE TABLE `user` (
  `userId` varchar(200) NOT NULL,
  `password` varchar(200) NOT NULL,
  `username` varchar(200) NOT NULL,
  PRIMARY KEY (`userId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

~~~~

3,设置项目为mod模式，设置proxy为https://goproxy.io

4，go run main.go

5，访问本地5656端口

