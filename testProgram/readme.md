# MQTT-MOCK

mqtt模拟发布程序，发布数据来源于mock.json

## 使用方法

1. 保证执行目录下存在mock.json文件，可以从mock-sample.json复制成mock.json
1. 执行对应操作系统版本的可执行文件，默认会加载mock.json文件
1. 可执行文件可以通过 `-f` 参数指定加载文件，例如 `mqtt-mock-darwin -f mock-sample.json`

## mock.json的说明

1. `config` 表示连接配置信息，根据情况填写
1. `data` 表示模拟发布的数据信息
    - `repeat` 表示循环重复发布的信息，根据逝去的时间和 `wait` 进行循环发布
    - `once` 表示根据时间 `wait` 只进行一次发布的信息
1. 两个发布信息中的 `wait` 值间隔大于100毫秒
1. 配置文件中的 `"[timestamp]"` 会替换成为发布时候的13位时间戳字符串，例如 `"1669106552199"`
1. 配置文件中的 `"[timestampINT]"` 会替换成为发布时候的13位时间戳字符串并转换成整型类型，例如 `1669106552199`
1. 配置文件中的 `payload` 根据实际情况进行填写，模拟程序只会替换 `[timestamp]` ，其它信息原样发布
1. 支持配置需要上传到MinIO的文件，可查看mock-sample.json中`config`的 `minioEndpoint`、`minioAccessKeyId`、`minioSecretAccessKey` 和`data`中的 `upload` 配置项
1. 在`upload` 配置中 `bucket`代表目标桶，`srcPath`代表本地路径，`destPath`代表桶中的目标路径，`contentType`代表类型
