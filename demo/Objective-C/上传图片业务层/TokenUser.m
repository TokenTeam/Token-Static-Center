//
//  TokenUser.m
//  TokenUser
//
//  Created by 陈雄&武嘉晟 on 16/3/27.
//  Copyright © 2016年 com.feelings. All rights reserved.
//

#import "TokenUser.h"
#import "TokenImageAssistant.h"

@interface TokenUser()

@end

@implementation TokenUser

#pragma mark - 上传头像

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

@end
