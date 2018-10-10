# Token-Static-Center

> Token团队静态资源引擎
>
> LiuFuXin @ Token Team 2018 <loli@lurenjia.in>
>
> **该文档为运维部署文档，供运维部署业务（迁移业务）使用**

## 部署过程

### 安装必要依赖

安装jpeg、tiff、gif、png、bmp、webp的二进制安装包与编译头文件，此处以RHEL 7系列（Fedora、CentOS）为例，如果出现不存在的包，先安装epel：

```shell
yum install epel-release
yum makecache
yum install -y libjpeg-turbo libjpeg-turbo-devel libjpeg-turbo-utils libpng libpng-devel libpng-static libtiff libtiff-devel libtiff-static giflib giflib-devel giflib-utils libwebp libwebp-devel libwebp-tools make gcc gcc-g++ libtool pkgconfig zip unzip tar gzip freetype freetype-devel
```

### 安装go

> 安装之前先卸载原有的Go（yum的go版本太老了）

```shell
yum remove -y go
```

1. 在[Golang官方网站](https://golang.org/)下载最新的Golang安装包

2. 解压安装包，将安装包中的go文件夹复制到`/usr/local/`目录

3. 编辑`/etc/profile`文件，在最后追加以下内容

   ```shell
   export PATH=$PATH:/usr/local/go/bin
   ```

4. 假如GoPath是`/home/gopath/`，在`/etc/profile/`文件后面继续追加内容：

   ```shell
   export GOPATH=/home/gopath
   export PATH=$PATH:$GOPATH/bin
   ```

5. 保存文档，执行`source /etc/profile`，加载当前配置
6. 输入`go version`命令，查看go版本，如果版本与下载的版本一致，即安装成功

### 安装Glide

> Glide是Go的包管理器

1. 执行`curl https://glide.sh/get | sh`
2. 输入`glide -v`命令，查看glide版本，如果有响应，即安装成功

### 安装ImageMagick

> 该项目仅依赖ImageMagick6.9，最新的Imagick7将无法兼容

1. 在[ImageMagick官方GitHub仓库](https://github.com/ImageMagick/ImageMagick6/releases)下载ImageMagick6的最新版本

2. 解压，进入ImageMagick开头的目录

3. 执行以下命令

   ```shell
   ./configure --with-webp
   make -j 4
   make install
   ```

4. 执行`magick -version`查看版本是否正确，如果版本正确，即安装成功

5. 在`/etc/profile`文件后面追加以下内容并保存

   ```shell
   export PKG_CONFIG_PATH="/usr/local/lib/pkgconfig"
   ```

6. 执行`source /etc/profile`命令，保存当前配置

### 部署业务代码

> 由于某些原因，此处部署业务代码过程可能会面临被墙的风险（Golang是Google的产品）
>
> 建议在本地，配合VPN部署完成后，上传到线上

1. 进入`$GOPATH/src/github.com/TokenTeam/Token-Static-Center`（如果没有，则新建一个）
2. 执行`git clone https://github.com/TokenTeam/Token-Static-Center/ .`，等待数分钟
3. 执行`glide install`
4. 执行`glide up`

### 编译业务代码

> 不推荐在本机离线编译，建议在线编译

1. 执行`go build init.go`
2. 将生成的`init`可执行文件修改名称为`token-static-center`

> 注意：根据压力测试，go1.11编译出来的代码会遇到大量内存越界的错误，使用go1.10则不会出现，在编译的时候要灵活选择，并在编译结束后进行压力测试

### 配置业务代码

1. 将以下文件（文件夹）从业务代码的文件夹内复制到指定目录（可以不放在一起，以下均以`/home/htdocs/static-img.wutnews.net/`为例）：`token-static-center` `config.yaml` `static/` `template/`，其中`template/` 需与 `token-static-center` 在同一目录下

2. 示例结构：

   ```
   |- home/
   |  |- htdocs/
   |  |  |- static-img.wutnews.net/
   |  |  |  |- token-static-center
   |  |  |  |- config.yaml
   |  |  |  |- template/
   |  |  |  |- static/
   ```

3. 编辑`config.yaml`文件，修改以下配置项（如果是路径，一定要带结尾的`/`符号）

   ```
   storage-dir: /home/htdocs/static-img.wutnews.net/static/storage/
   log-dir: /home/htdocs/static-img.wutnews.net/static/log/
   # db-resource根据实际情况修改，空的sqlite数据库就只是一个空白的.db文件而已，touch一下就能生成
   db-resource: /home/htdocs/static-img.wutnews.net/static/db/example.db
   cache-dir: /home/htdocs/static-img.wutnews.net/static/cache/
   ```

4. 根据实际情况（例如安全要求，具体配置需求，数据库模式等）修改其他的配置项（均有完备注释）

5. 执行`/home/htdocs/static-img.wutnews.net/token-static-center --config=/home/htdocs/static-img.wutnews.net/config.yaml`（初次启动建议打开配置文件中的调试模式）

6. 根据接入文档检查功能是否正常（上传&下载）

### 配置服务

> 此处仅以CentOS 7 作为示例

1. 编辑`/usr/lib/systemd/system/token-static-center.service`，输入以下内容（其中的ExecStart根据具体情况而定）

   ```ini
   [Unit]
   Description=Token Static Center
   After=syslog.target network.target nginx.service
   
   [Service]
   Type=simple
   ExecStart=/home/gopath/src/github.com/TokenTeam/Token-Static-Center/token-static-center --config=/home/gopath/src/github.com/TokenTeam/Token-Static-Center/config.yaml
   
   [Install]
   WantedBy=multi-user.target
   ```

2. 执行`systemctl daemon-reload`

3. 执行`systemctl start token-static-center`

4. 执行`netstat -nlp | grep token-static-center`，查看是否有监听网络的行为

5. 如果没有，说明此前配置错误

6. 如果有监听网络行为，执行`systemctl enable token-static-center`，设置为开机自启

### 反向代理服务

> 此处不再赘述，略去配置部分。

需要注意的两点：

- 注意缓存配置（Response Header中设定Cache-Control）
- 注意关闭Nginx本身的反向代理缓存（防止静态资源引擎记录日志不完整，统计与日志功能失效）