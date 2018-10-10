# Token-Static-Center

> Token团队静态资源引擎
>
> LiuFuXin @ Token Team 2018 <loli@lurenjia.in>
>
> **该文档为版本更新文档，供版本更新记录使用**

## Version 1.10 2018-10-10

- 解决通过网页上传的情况下无可避免触发反盗链模块导致的阻塞性BUG-c0d72197dfb097864f2decc8418e07e44db0547d
- 新增Dockerfile，便于便捷部署业务

## Version 1.01 2018-10-07

- 修复高负载下频繁启动Imagick进程造成性能瓶颈（服务器无法及时返回数据）的BUG-d8b91a43970f8cac4dc8ac326e2dbd6310cbadc6
- 解决性能较强+多核心环境下线程不同步导致主进程Crash的BUG-9203beb49be756c5f8fd1b13298258934a4ad0eb

## Version 1.00 2018-09-02

- 基础功能实现