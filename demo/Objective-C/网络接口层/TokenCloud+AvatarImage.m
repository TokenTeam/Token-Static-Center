//
//  TokenCloud+AvatarImage.m
//  MyWHUT
//
//  Created by 武嘉晟 on 2019/1/14.
//  Copyright © 2019 com.feelings. All rights reserved.
//

#import "TokenCloud+AvatarImage.h"
#import "TokenNetworking.h"
#import "NSError+TokenCloud.h"
#import "TokenSession.h"

@implementation TokenCloud (AvatarImage)

+ (void)fetchAccessTokenAndNonceWithsuccess:(TokenCloudAvatarImageSuccessBlock)success
                                    failure:(TokenCloudAvatarImageFailureBlock)faliure {
    TokenNetworking.networking
    .request(^NSURLRequest *{
        NSString *api = @"https://test-api-iwut.wutnews.net/user/user/make_avatar_token";
        NSString *token = [[TokenUser sharedUser] newToken];
        NSDictionary *para = @{
                               @"token":token
                               };
        NSMutableURLRequest *request = NSMutableURLRequest.token_requestWithURL(api)
        .token_setMethod(@"POST")
        .token_setHTTPParameter(para);
        return request;
    })
    .responseJSON(^(NSURLSessionTask *task, NSError *jsonError, id responsedObj) {
        if (jsonError) {
            !faliure?:faliure(jsonError);
        }
        if (responsedObj) {
            NSDictionary *json = responsedObj;
            NSString *msg = [json objectForKey:@"msg"];
            DebugLog(@"%@",msg);
            if ([[json objectForKey:@"code"] integerValue] == 0) {
                //para有两个字段，access_token和nonce
                NSDictionary *para = json[@"data"];
                !success?:success(para);
            } else {
                //code非0则代表有问题
                NSString *data = [json objectForKey:@"data"];
                if ([[data class] isKindOfClass:[NSString class]]) {
                    DebugLog(@"失败原因是：%@",data);
                }
                NSError *error = [NSError token_errorWithLocalDescription:@"状态码不正确"];
                !faliure?:faliure(error);
            }
        }
    })
    .failure(^(NSError *error) {
        !faliure?:faliure(error);
    });
}

+ (void)uploadImage:(UIImage *)image
           withPara:(NSDictionary *)para
            success:(TokenCloudAvatarImageSuccessBlock)success
            failure:(TokenCloudAvatarImageFailureBlock)faliure {
    //两个必备参数
    NSString *accessToken = para[@"access_token"];
    NSString *nonce = para[@"nonce"];
    //URL的拼接，请参考Token静态资源平台
    //格式 http://example.com/upload/accessToken-Nonce.fileFormat
    NSMutableString *urlString = [NSMutableString stringWithFormat:@"%@",@"https://static-img.wutnews.net/upload/"];
    __block NSString *imageFormat = @"";
    //拼接accessToken-Nonce
    [urlString appendFormat:@"%@-%@",accessToken,nonce];
    //拼接.fileFormat
    if (UIImagePNGRepresentation(image)) {
        //如果是png图片
        [urlString appendFormat:@"%@",@".png"];
        imageFormat = @"png" ;
    } else if (UIImageJPEGRepresentation(image, 0.5)) {
        //如果是jpeg图片
        [urlString appendFormat:@"%@",@".jpeg"] ;
        imageFormat = @"jpeg" ;
    } else {
        //图片格式错误
        NSError *error = [NSError token_errorWithLocalDescription:@"图片格式错误"];
        return !faliure?:faliure(error);
    }
    void (^successBlock)(id) = ^(id uuid) {
        NSMutableString *targetUrlString = [NSMutableString stringWithFormat:@"%@",@"https://static-img.wutnews.net/image/"];
        NSInteger width = image.size.width;
        //静态资源平台只返回uuid，完整的URL还需要拼接部分东西
        [targetUrlString appendFormat:@"%@-%ld-100.%@",uuid,(long)width,imageFormat];
        !success?:success(targetUrlString) ;
    };
    NSDictionary *dict = @{@"image":image};
    [TokenSession uploadDataTaskUseFormData:@"POST"
                                        url:urlString
                                       data:dict
                                    success:successBlock
                                    faliure:faliure];
}

+ (void)uploadAvatarImageUrlStr:(NSString *)urlStr
                        success:(TokenCloudAvatarImageSuccessBlock)success
                        failure:(TokenCloudAvatarImageFailureBlock)faliure {
    TokenNetworking.networking
    .request(^NSURLRequest *{
        NSString *api = @"https://test-api-iwut.wutnews.net/user/user/upload_avatar";
        NSString *token = [[TokenUser sharedUser] newToken];
        NSString *cardno = [TokenUser sharedUser].cardno;
        NSDictionary *para = @{
                               @"token":token,
                               @"cardno":cardno,
                               @"avatar_path":urlStr,
                               @"Referer":@"http://token.wutnews.net/tools/image"
                               };
        NSMutableURLRequest *request = NSMutableURLRequest.token_requestWithURL(api)
        .token_setMethod(@"POST")
        .token_setHTTPParameter(para);
        return request;
    })
    .responseJSON(^(NSURLSessionTask *task, NSError *jsonError, id responsedObj) {
        if (jsonError) {
            !faliure?:faliure(jsonError);
        }
        if (responsedObj) {
            NSDictionary *json = responsedObj;
            DebugLog(@"msg:%@",json[@"msg"]);
            if ([[json objectForKey:@"code"] integerValue] == 0) {
                //这个最好不要用吧
                !success?:success(json[@"data"]);
            } else {
                //code非0则代表有问题
                NSString *data = [json objectForKey:@"data"];
                if ([[data class] isKindOfClass:[NSString class]]) {
                    DebugLog(@"失败原因是：%@",data);
                }
                NSError *error = [NSError token_errorWithLocalDescription:@"状态码不正确"];
                !faliure?:faliure(error);
            }
        }
    })
    .failure(^(NSError *error) {
        !faliure?:faliure(error);
    });
}

@end
