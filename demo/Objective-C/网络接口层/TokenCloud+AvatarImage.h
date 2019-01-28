//
//  TokenCloud+AvatarImage.h
//  MyWHUT
//
//  Created by 武嘉晟 on 2019/1/14.
//  Copyright © 2019 com.feelings. All rights reserved.
//

#import "TokenCloud.h"

NS_ASSUME_NONNULL_BEGIN

typedef void(^TokenCloudAvatarImageSuccessBlock)(id obj);
typedef void(^TokenCloudAvatarImageFailureBlock)(NSError *error);

@interface TokenCloud (AvatarImage)

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

@end

NS_ASSUME_NONNULL_END
