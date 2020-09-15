## j2pay
 多用户钱包管理系统服务端

* 项目目录临时结构
```
             ├─conf               配置文件目录
             ├─controller         控制器
             ├─docs               swagger接口文档
             ├─log                日志信息
             ├─middleware         中间件
             ├─model              实体
             │  ├─request         请求参数对象
             │  └─response        返回参数对象
 J2PAY       ├─myerr              错误处理
             ├─pkg                包      
             │  ├─casbin          鉴权
             │  ├─logger          日志
             │  ├─setting         项目配置
             │  └─util            工具
             ├─routers            网关配置
             ├─service            业务处理
             └─validate           参数验证
 ```
* 开发框架
```
1.Gin
2.JWT
3.Session
4.Casbin
5.Gorm
6.Swagger
7.Logrus
8.Mysql
```

## swagger 接口文档地址
* [api文档](http://192.168.3.55:8088/swagger/index.html)


