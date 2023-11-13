# MOCK SERVER

http 服务器模拟程序

## 使用方法

1. 保证执行目录下存在serverConfig.json文件，可以用同名文件替换serverConfig.json
1. 执行对应操作系统版本的可执行文件，默认会加载serverConfig.json文件
1. 可执行文件可以通过 `-f` 参数指定加载文件，例如 `mockServer-darwin -f serverConfig.json`

## serverConfig.json的说明

1. `port` 表示启动服务器时监听的端口。请直接填写端口号，不需要加冒号。
1. `path` 表示需要注册的路由器路径，相同的路径只会注册一个路由。