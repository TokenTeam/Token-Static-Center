# 上传头像的iOS示范代码

## 介绍

向团队静态资源平台上传图片并且让谢站服务器知道图片URL。

iOS这一块相对比较复杂，尤其是request的生成，本示范详细阐述了如何向平台上传图片。

## 上传图片业务层

这是一个对外的业务层API，所有想要上传头像的开发者可以通过如下入口上传

```objective-c
/**
 *  上传图片
 */
- (void)uploadAvatarImage:(UIImage *)image
                  success:(TokenUserLoginSuccessBlock)success
                  failure:(TokenUserLoginFailureBlock)failure;
```

参数给图片，成功回调和失败回调。两个代码块的定义如下

```objective-c
typedef void(^TokenUserLoginSuccessBlock) (id object);
typedef void(^TokenUserLoginFailureBlock) (NSError *respnseError);
```

该方法的实现如下

``` 
/**
 上传头像
 
 @param image 图片
 @param success 成功回调
 @param failure 失败回调
 */
- (void)uploadAvatarImage:(UIImage *)image
                  success:(TokenUserLoginSuccessBlock)success
                  failure:(TokenUserLoginFailureBlock)failure {
    //获取AccessToken和Nonce
    [TokenCloud fetchAccessTokenAndNonceWithsuccess:^(id  _Nonnull obj) {
        //获取成功
        //obj是个字典，有两个字段，分别是AccessToken还有Nonce
        NSDictionary *para = obj;
        //上传图片到团队静态资源平台
        [TokenCloud uploadImage:image withPara:para success:^(id  _Nonnull obj) {
            //上传成功
            //obj是个字符串，是静态资源平台给的图片URL
            NSString *imageURL = obj;
            //上传图片url到团队服务器
            [TokenCloud uploadAvatarImageUrlStr:imageURL success:^(id  _Nonnull obj) {
                //上传url成功
                //obj是URL字符串
                !success?:success(obj);
                //持久化
                self.avatarImage = image;
                [TokenImageAssistant saveImage:image forKey:obj];
            } failure:^(NSError * _Nonnull error) {
                //上传url失败
                !failure?:failure(error);
            }];
        } failure:^(NSError * _Nonnull error) {
            //上传到静态资源平台失败
            !failure?:failure(error);
        }];
    } failure:^(NSError * _Nonnull error) {
        //获取AccessToken和Nonce失败
        !failure?:failure(error);
    }];
}
```

第一步需要向谢站服务器申请AccessToken和Nonce两个参数，这两个参数用于第二步的工作。

第二步是向团队静态资源平台上传图片，平台会返回一个初步的URL

第三部，拼接URL，向谢站服务器上传这个URL。

## 网络接口层

网络接口层是业务代码和网络进行沟通的一层，AvatarImage分类主管头像接口，对外有三个API。

```objective-c
/**
 *  获取AccessToken和Nonce

 @param success 成功回调
 @param faliure 失败回调
 */
+ (void)fetchAccessTokenAndNonceWithsuccess:(TokenCloudAvatarImageSuccessBlock)success
                                    failure:(TokenCloudAvatarImageFailureBlock)faliure;
/**
 *  上传图片到团队静态资源平台
 
 @param para 参数字典
 @param success 成功回调
 @param faliure 失败回调
 */
+ (void)uploadImage:(UIImage *)image
           withPara:(NSDictionary *)para
            success:(TokenCloudAvatarImageSuccessBlock)success
            failure:(TokenCloudAvatarImageFailureBlock)faliure;
/**
 *  上传图片URL到服务器
 
 @param urlStr 图片URL
 @param success 成功回调
 @param faliure 失败回调
 */
+ (void)uploadAvatarImageUrlStr:(NSString *)urlStr
                        success:(TokenCloudAvatarImageSuccessBlock)success
                        failure:(TokenCloudAvatarImageFailureBlock)faliure;
```

介绍的很详细了，block的定义如下。

```objective-c
typedef void(^TokenCloudAvatarImageSuccessBlock)(id obj);
typedef void(^TokenCloudAvatarImageFailureBlock)(NSError *error);
```

## 网络基础层

这一块是基础的TokenNetworking和TokenSession。

TokenNetworking封装了iOS原生的网络请求代码，方便的进行网络请求。

TokenSession则可以以表单形式上传data，基于TokenNetworking。

表单形式上传数据的难度在于request的构造，需要加上合理的请求头，image需要转化为data，并且data前后需要拼接特定的data。

具体细节请直接看代码