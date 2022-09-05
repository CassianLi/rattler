# Getting started

> 当前应用可监听`soft pak export`文件的生成路径，并将`export` 内容以`Json` 字符串的形式发送到消息队列中。监听`Import` 报关文档的消息队列，将报关文档内容以`xml` 文件形式保存在报关路径中，完成报关的**发送/接收** 功能，信息的传递均通过消息队列完成。同时还可启动**文件服务器** 以提供`api`访问下载`export` 文件和**税金单（PDF）** 。



## 命令介绍

直接运行`./rattler` 查看命令选项

```shell
Rattler is used to communicate with SoftPak software to send or receive customs documents. For example:

Usage:
  rattler [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  listen      Monitor the messages in the message queue and save the corresponding messages as XML files to the specified path
  serve       Start a web server
  watch       Monitor changes to files in the specified file path

Flags:
      --config string   config file (default is $HOME/.rattler/config.yaml)
  -h, --help            help for rattler

Use "rattler [command] --help" for more information about a command.

```

命令参数：

- `rattler watch` 监听`export` 路径，并发送`export` 文件内容及文件名到指定消息队列，完成报关结果的获取
- `rattler listen` 监听`import` 报关文件的消息队列，并保存报关文档到指定路径
- `rattler serve`  启动`web` 文件服务器，提供接口访问**税金单文件和Export xml文件**



## 子命令启动

下面分别介绍各子命令启动及配置文件内容

### rattler watch

`Export` 文件监听命令，需要监听不同的`declare country` 则需要以不同的配置文件来启动命令。如：

```shell
# 监听 NL export 报关结果
rattler watch --config .rattler/watch-nl.yaml

```



**.rattler/watch-nl.yaml**

```shell
# Export file listener configuration file

# Declare country
declare-country: NL


log:
  # log level: debug | info | warn | error
  level: debug
  # log path
  filename: log/rattler-watch-nl.log

watcher:
  # file listening path
  watch-dir:
  # file backup path
  backup-dir:

rabbitmq:
  url: amqp://USER:PASSWORD@MQ_HOST:5672
  exchange: softpak.export.topic
  queue-prefix: softpak.export
  queue-count: 1

```



**.rattler/watch-be.yaml**

```shell
# Export file listener configuration file

# Declare country
declare-country: BE


log:
  # log level: debug | info | warn | error
  level: info
  # log file path
  filename: log/rattler-watch-be.log

watcher:
  # file listening path
  watch-dir:
  # file backup path
  backup-dir:

rabbitmq:
  url: amqp://USER:PASSWORD@MQ_HOST:5672
  exchange: softpak.export.topic
  queue-prefix: softpak.export
  queue-count: 1


```



### rattler listen

监听报关消息队列`softpak.import` 将报关内容保存到指定报关文件存放路径。

**启动命令**

```shell
# 监听Import队列，获取报关文档

rattler listen --config .rattler/listen.yaml
```



**.rattler/listen.yaml**

```shell
# This is the customs declaration data monitoring configuration file
log:
  # log level
  level: info
  # log file path
  filename: log/rattler-listen.log

# Import XML save directory
import-dir:

# Customs declaration file monitoring queue
rabbitmq:
  url: amqp://USER:PASSWORD@HOST_IP:5672
  exchange: softpak.import.topic
  queue: softpak.import

```

<u>*==注意：密码中带有特殊字符时注意特殊符号的转义，使用**单引号**包裹包含特殊字符的字符串（如：‘amqp://user:@123!123@127.0.0.1:5672’）==*</u>

### rattler serve

`soft pak` 税金单以及`export xml` 文件访问`web`服务器。



**启动命令**

```shell
rattler serve --config .rattler/serve.yaml
```



**.rattler/serve.yaml**

```shell
# Soft pak file server's port
port: 7003

log:
  # log level debug | info | warn | error
  level: info
  # log file pat
  filename: log/rattler-serve.log

directory:
  # Tax bill file path
  tax-bill: 
  	nl: ${tax_bill_dir_for_nl}
  	be: ${tax_bill_dir_for_be}
  export:
    # Export file backup path: nl | be
    nl: ${export_dir_for_nl}
    be: ${export_dir_for_be}
```



## API 接口说明

### Download tax bill PDF

下载税金单文件，用户指定访问源文件名和下载后的文件名进行文件下载。

#### Link：http://127.0.0.1:7003/download/pdf/:origin/:target?dc=nl

| 参数   | 必填 | 类型 | 说明                                             |
| ------ | ---- | ---- | ------------------------------------------------ |
| origin | y    | Str  | 需要下载的源文件名（***不包含文件后缀***）       |
| Target | y    | Str  | 下载后保存为指定的文件名（***不包含文件后缀***） |
| Dc     | y    | Str  | declare country (**NL\|BE**)                     |

***例如：http://localhost:7003/download/pdf/tax-bill/save_tax_bill?dc=nl***

### Download export XML

下载或访问`export` 结果文件

#### Link: http://127.0.0.1:7003/download/xml/:dc/:filename?download=1

| 参数     | 必填 | 类型 | 说明                                                         |
| -------- | ---- | ---- | ------------------------------------------------------------ |
| Dc       | y    | Str  | declare country (**NL\|BE**)                                 |
| Filename | y    | Str  | 文件名（***带文件后缀的完整文件名***）                       |
| Download | n    | Str  | 是否以下载的方式访问（***1｜0 ***），如果为1 则表示直接下载文件，其他值可在浏览器预览文件内容 |

***例如：http://localhost:7003/download/xml/be/17960_AI-2021-77635_18.xml?download=***



## 部署服务

### windows 部署

以` rattler serve --config .rattler/serve.yaml`  为例进行介绍，其他命令于此相同。

**1. 修改配置文件`.rattler/serve.yaml` **

```yaml
# Soft pak file server's port
port: 7003

log:
  # log level debug | info | warn | error
  level: info
  # log file pat
  filename: log/rattler-serve.log

directory:
  # Tax bill file path
  tax-bill: 
  	nl: ${tax_bill_dir_for_nl}
  	be: ${tax_bill_dir_for_be}
  export:
    # Export file backup path: nl | be
    nl: ${export_dir_for_nl}
    be: ${export_dir_for_be}
```

**2. 命令启动**

```shell
./rattler.exe serve --config .rattler/serve.yaml
```



**3. 以系统服务方式启动**

`windows` 非系统服务的命令不能直接使用`windows` 系统自带的命令创建系统服务。因此需要借助第三方工具来代替我们运行系统服务。

如: [winsw.exe](https://github.com/winsw/winsw/releases/tag/v2.11.0) ,下面以 `winsw.exe` 进行介绍:

- 首先将`winsw.exe` **重命名**（`exp:rattler-serve.exe`），创建`xml`配置文件**（配置文件名需要和执行文件名相同：`rattler-serve.xml`）**

```shell
<service>
    <id>rattler-file-service</id>
    <name>RattlerFileService</name>
    <description>This is a service that can be used to access SoftPak declaration documents</description>
    <executable>..\rattler.exe</executable>
    <arguments>serve --config ..\.rattler\serve.yaml</arguments>
    <logmode>reset</logmode>
</service>

```

- 安装系统服务

```pow
rattler-serve.exe install
```

- 其他命名

```postgresql
# 启动服务
rattler-serve.exe start

# 停止服务
rattler-serve.exe stop

# 卸载服务
rattler-serve.exe uninstall
```



### Rattler 所有服务配置内容

`rattler` 打包文件中已经包含`winsw` 针对所有命令的配置，存放在路径`winsw/`中。文件结构如下：

```shell
➜  winsw ✗ tree
.
├── rattler-listen.exe
├── rattler-listen.xml
├── rattler-serve.exe
├── rattler-serve.xml
├── rattler-watch-be.exe
├── rattler-watch-be.xml
├── rattler-watch-nl.exe
└── rattler-watch-nl.xml

```

