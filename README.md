# Token-Static-Center

> Token团队静态资源引擎
>
> LiuFuXin @ Token Team 2018 <loli@lurenjia.in>
>
> **该文档为接入文档，供各业务接入使用**
>
> - 如需运维部署文档，请参考同目录下OPS.md
> - 如需开发、维护辅助文档，请参考同目录下CONTRIBUTE.md
> - 如需查看开发历史、待实现功能，请参考同目录下CHANGELOG.md



## 接入指南

1. 找到配置文件中的`security`项
2. 将类似于`https://web.wutnews.net/`的完整域名（必须以/结尾，不包含子目录）加入`white-list`配置中
3. 使用任意方式（脚本、在线、滚键盘皆可）生成一段64位含大小写数字的AppCode，加入`app-code`配置中，类似于`6Jukw2pPvR0zWT3qJP3mKNYI1INfiQsqYkdGM9OPltW3JlRSBjPoFIwYAdq2XuKt`
4. 重启Token-Static-Center即可生效
5. 将AppCode引入需要接入的业务系统，并按照后文中的AccessToken计算方法，对AccessToken进行生成

## 请求格式

### 获取图片资源
#### 请求图片（带图片水印）
http://example.com/image/GUID-width-quality-watermarkName-watermarkPosition-watermarkOpacity-watermarkSize.fileExtension

共计八个可变参数：

- GUID 图片唯一资源识别码
- width 图片宽度（px）
- quality 图片质量（0-100，越高质量越好，通常只对jpg格式有效）
- watermarkName 水印名称（存储在静态资源目录/watermark/``watermarkName`.png）
- watermarkPosition 水印位置
  - 1 左上角
  - 2 右上角
  - 3 左下角
  - 4 右下角
  - 5 正中央
- watermarkOpacity 水印透明度（0-100，值越高越透明）
- watermarkSize 水印宽度（0-100，相对于所请求图片的宽度占百分比）
- fileExtension 图片格式，以配置文件内所支持的格式列表为准

> 例：http://static-img.wutnews.net/image/e44378ac-0237-4331-aaf2-63b8818e5c34-300-80-wutnews-1-30-15.jpg 即为请求GUID为 e44378ac-0237-4331-aaf2-63b8818e5c34，宽度为300，质量为80，水印名称为wutnews，水印位置为左上角，水印透明度为30%透明，水印大小为15%宽度（相对于图片宽度）的JPG格式图片资源
**注意，对于GIF动图，无法添加水印**

------

#### 请求图片（不带任何水印）~~（无印良品）~~
http://example.com/image/GUID-width-quality.fileExtension

共计四个可变参数：

- GUID 图片唯一资源识别码
- width 图片宽度（px）
- quality 图片质量（0-100，越高质量越好，通常只对jpg格式有效）
- fileExtension 图片格式，以配置文件内所支持的格式列表为准

> 例：http://static-img.wutnews.net/image/e44378ac-0237-4331-aaf2-63b8818e5c34-300-80.jpg 即为请求GUID为e44378ac-0237-4331-aaf2-63b8818e5c34，宽度为300，质量为80，不带水印的JPG格式图片资源

------

#### 请求图片（带文字水印）

http://example.com/image/GUID-width-quality-text-fontPosition-fontOpacity-fontSize-fontColor-fontStyle.fileExtension

共计十个可变参数：

- GUID 图片唯一资源识别码
- width 图片宽度（px）
- quality 图片质量（0-100，越高质量越好，通常只对jpg格式有效）
- text 水印文本
- fontPosition 水印位置（同请求图片-带图片水印）
- fontOpacity 水印透明度（0-100，值越高越透明）
- fontSize 水印字体大小（px）
- fontColor 水印颜色（16进制码，不含#符号，大小写不限，例如黑色：000000，白色：FFFFFF）
- fontStyle 水印字体样式
  - light 细字体
  - bold 粗体
  - regular 普通字体
- fileExtension 图片格式，以配置文件内所支持的格式列表为准

**切记，text参数内不允许包含特殊符号（只允许包含@符号），否则会造成转义错误**

> 例：http://static-img.wutnews.net/image/e44378ac-0237-4331-aaf2-63b8818e5c34-300-80-%40Token+Team-1-20-30-FFFFFF-regular.jpg 即为请求GUID为 e44378ac-0237-4331-aaf2-63b8818e5c34，宽度为300，质量为80，水印文本为@Token Team，水印位置为左上角，水印透明度为20%透明，水印字体大小为30px，水印颜色为FFFFFF，水印字体样式为普通字体样式的JPG格式图片资源



### 上传图片资源

http://example.com/upload/accessToken-Nonce.fileFormat

方法：POST表单，其中

- input type="file" name="image"

共计三个可变参数：

- accessToken 许可密钥
- Nonce 16位随机小写字母
- fileFormat 当前图片的格式（以后缀名为准，小写）

返回数据为JSON格式：

> 例：{"error_code": 0, "message": "GUID"}

如果存在错误：

> 例：{"error_code": 小于0的整数, "message": "错误相关提示"}

### AccessToken计算方法

md5(AppCode前32位数+时间戳去掉最后四位数+AppCode后32位数+随机Nonce+配置文件中设置的SaltString)

其中AppCode为64位小写字母&数字随机数，存储于静态资源引擎的配置文件中，每个业务一个AppCode

Nonce为16位随机小写字母

> *建议由服务端传递给前端，前端不要完成计算过程以保障安全性*
>
> 时间戳去掉最后四位的原因：该AccessToken在999s内均为有效，可以保证客户端在999s内上传文件均为合法（注意：服务器与客户端时间间隔不能过大）



## 错误页面类型及排查方案

### 404错误

- 资源不存在
- GUID格式有误
- 参数个数有误
- 参数数值范围有误

### 403错误

- http referrer有误，或不在白名单中
- AccessToken校验失败，请检查客户端与服务端的时间设置是否一致，检查AppCode是否正确

### 500错误

- 服务器内部错误，原因多样，请记录下出现错误的时间，结合日志进行错误排查

### 服务器不响应

- 服务中断（建议访问其他资源检查是否存在服务中断情况）
- 参数有误（请根据前面的请求格式进行）