# Pricking

[Pricking](https://github.com/Rvn0xsy/Pricking) 是一个自动化部署水坑和网页钓鱼的项目。

想要了解更多可以阅读：

- [红队技巧：基于反向代理的水坑攻击](https://payloads.online/archivers/2021-02-16/1)
- [Pricking 项目（一） ：使用介绍](https://payloads.online/archivers/2021-02-18/1)
- [Pricking 项目（二） ：JS模块开发](https://payloads.online/archivers/2021-02-18/2)

:collision: :collision: :collision: 支持HTTPS/HTTP

## 使用方法

更多使用方式可以参考 [Pricking Wiki](https://github.com/Rvn0xsy/Pricking/wiki)

> 使用本项目需要拥有一个域名，将A记录指向到当前服务器，否则只能通过IP访问。

### 安装方式 - 下载二进制文件

[Releases](https://github.com/Rvn0xsy/Pricking/releases)

### 安装方式 - 编译

```bash
$ git clone https://github.com/Rvn0xsy/Pricking
$ cd Pricking
$ make
```

### 配置文件

```yaml
filter_type:
  - "text/html" # 仅针对网页内容进行注入
exclude_file:   # 静态文件的数据包不进行注入
  - ".jpg"
  - ".css"
  - ".png"
  - ".js"
  - ".ico"
  - ".svg"
  - ".gif"
  - ".jpeg"
  - ".woff"
  - ".tff"
static_dir: "./static" # Pricking Suite 目录
pricking_prefix_url: "/pricking_static_files" # 静态目录名，不能与目标网站冲突
listen_address: ":9999" # 监听地址:端口
inject_body: "<script src='/pricking_static_files/static.js' type='module'></script>" # 注入代码
```

## [Pricking Js Suite 模块说明](./static/)

- modules/cookie.js 获取网页Cookie并打印在控制台上
- ...

## 引入方式

在static.js中添加：

```js
import * as <ModuleName> from './modules/<ModuleName>.js'
```

例如 `cookie.js`：

```js
import * as Cookie from './modules/cookie.js'
```

## 贡献

请为我提交[Pull Request](https://github.com/Rvn0xsy/Pricking/pulls)

