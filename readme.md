### 超轻量web框架

- 服务器引擎

  基于go官方http包，进行再封装，支持Gzip压缩

  - 静态文件服务器
    - 可以指定任意资源目录输出目录内文件
    - 自定义默认文件 
    - 在Red Hat 4.4.7 1核1G配置下支持5000并发文件请求

  - 动态逻辑处理服务器
    - 支持多种HTTP请求方式(GET、POST、PUT、DELETE)
    - 支持websocket处理

- 路由配置
  - restful风格路由自动加载
    - 支持以Get|GET|Post|POST|Put|PUT|Delete|DELETE为前缀的HandleFunc自动注册为对应请求方式的资源路径处理器
    - 支持单独注册路径及对应处理器HandleFunc
    - 支持注册404处理方法
    - 支持注册中间件方法

- 日志
  - 支持终端打印请求日志
