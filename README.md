# PasteBin-go
![Go](https://img.shields.io/github/go-mod/go-version/CharlieYu4994/PasteBin-go)![GPLv3](https://img.shields.io/github/license/CharlieYu4994/PasteBin-go) ![GitHub last commit](https://img.shields.io/github/last-commit/CharlieYu4994/PasteBin-go)

一个简易的 PasteBin 的 Go 语言实现，并支持自定义的过期时间设定

## 特点
+ 使用 Golang 编写，不使用任何外部模块，较为轻量化
+ 使用 LinkedHashMap 存储数据，并实现 LRU，不会轻易爆内存
+ 提供较多自定义设置项，方便部署

## 使用

### 环境要求
+ 任意一台服务器 (required)
+ Go 编译环境 (required)

### 编译
1. 在 [Release](https://github.com/CharlieYu4994/PasteBin-go/releases) 页面下载源码
2. 解压并进入文件夹
3. 在文件夹下运行 `go build ./`

最终你将得到一个可执行文件

### 部署
1. fork 这个存储库
2. 在你的存储库的 web 分支中找到『config.json』
3. 将里面的 `backend` 字段改成你的服务器的域名 (注意端口)
4. 在 web 分支打开 GitHub Page
5. 根据『config.json.template』创建配置文件
6. 将可执行文件和配置文件放在同一目录下运行
7. 访问你的 GitHub Page 查看效果
> Tips：这边建议你使用 Nginx 之类的软件对 PasteBin 进行反向代理

## 官方演示
更多详情请查看 [Demo 主页](https://pastebin.charlie.moe)

> Tips：这个 Demo 最大粘贴条目数量为 300条，请合理使用

> Tips：作为 PasteBin 服务，不对数据做任何可靠性保证

> 此 Demo 由 [LassiCat](https://github.com/LassiCat) 的服务器托管运行 ~~（其实就是咱的）~~

## 版权
本程序由 [@CharlieYu4994](https://blog.charlie.moe/) 编写，以 GPLv3 协议发布

**本程序不支持，不鼓励一切商业用途**