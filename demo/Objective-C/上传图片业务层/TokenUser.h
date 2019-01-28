//
//  TokenUser.h
//  TokenUser
//
//  Created by 陈雄&武嘉晟 on 16/3/27.
//  Copyright © 2016年 com.feelings. All rights reserved.
//

#import <Foundation/Foundation.h>

@import UIKit;

typedef void(^TokenUserLoginSuccessBlock) (id object);
typedef void(^TokenUserLoginFailureBlock) (NSError *respnseError);

@interface TokenUser : NSObject

/**
 *  上传图片
 */
- (void)uploadAvatarImage:(UIImage *)image
                  success:(TokenUserLoginSuccessBlock)success
                  failure:(TokenUserLoginFailureBlock)failure;

@end
