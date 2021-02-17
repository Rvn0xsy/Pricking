# Pricking

[Pricking](https://github.com/Rvn0xsy/Pricking) 是一个自动化部署水坑和网页钓鱼的项目。

想要了解更多可以阅读：

- [红队技巧：基于反向代理的水坑攻击](https://payloads.online/archivers/2021-02-16/1)
- [Pricking 项目（一） ：使用介绍](https://payloads.online/archivers/2021-02-18/1)
- [Pricking 项目（二） ：JS模块开发](https://payloads.online/archivers/2021-02-18/2)

## Usage

> 使用本项目需要拥有一个域名，将A记录指向到当前服务器，否则只能通过IP访问。

```
$ git clone https://github.com/Rvn0xsy/Pricking
$ cd Pricking
# 修改Docker-Compose环境变量设置为你要克隆的网站
$ docker-compose up -d
```

### 查看日志

```
$ tail -f access.log
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

