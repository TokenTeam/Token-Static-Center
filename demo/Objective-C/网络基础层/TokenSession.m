//
//  TokenSession.m
//  MyWHUT
//
//  Created by ly&武嘉晟 on 2019/1/21.
//  Copyright © 2019 com.feelings. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "TokenSession.h"
#import "TokenNetworking.h"

NSString * const NSCustomErrorDomain = @"TokenSessionErrorDomain";

FOUNDATION_EXPORT NSErrorDomain const NSCustomErrorDomain;

@interface TokenSession ()

+ (NSString *)getDemarcationString;

@end

@implementation TokenSession

+ (NSString *)getDemarcationString {
    NSString *demarcation = @"QWERTYUIOPASDFGHJKLZXCVBNM" ;
    return  demarcation;
}

/**
 *  以表单的格式上传数据
 
 @param method 请求方式
 @param urlString 服务器url
 @param dataDictionary 上传的数据字典
 @param success 上传成功的回调
 @param faliure 方法失败的回掉
 */
+ (void)uploadDataTaskUseFormData:(NSString *)method
                              url:(NSString *)urlString
                             data:(NSDictionary *)dataDictionary
                          success:(TokenSessionSuccessBlock)success
                          faliure:(TokenSessionFailureBlock)faliure {
    NSURL *url = [NSURL URLWithString:urlString];
    NSMutableURLRequest *request  = [NSMutableURLRequest requestWithURL:url cachePolicy:0 timeoutInterval:30];
    NSString *demarcation = [TokenSession getDemarcationString];
    //HTTP方法设置
    request.HTTPMethod = method ;
    //用于拼接上传的数据
    NSMutableData *data = [NSMutableData data];
    //把上传的数据字典转化为NSData
    [dataDictionary enumerateKeysAndObjectsUsingBlock:^(id  _Nonnull key, id  _Nonnull obj, BOOL * _Nonnull stop) {
        //headerString就是QWERTYUIOPASDFGHJKLZXCVBNM在拼接花里胡哨的，图片的data包含在头尾之间
        NSMutableString *headerString = [NSMutableString stringWithFormat:@"\r\n--%@\r\n",demarcation] ;
        //如果是图片类型
        if ([obj isKindOfClass:[UIImage class]]) {
            NSString *font;
            UIImage *image = obj;
            NSData *imageData;
            if (UIImagePNGRepresentation(image)) {
                //如果是png
                imageData = UIImagePNGRepresentation(image);
                font = @".png";
                [headerString appendFormat:@"Content-Disposition: form-data; name=\"%@\"; filename=\"%@%@\"\r\n",key,@"avatarImage",font];
                [headerString appendFormat:@"Content-Type: image/png\r\n\r\n"];
            } else if (UIImageJPEGRepresentation(image, 0.5)) {
                //如果是jpeg
                imageData = UIImageJPEGRepresentation(image, 0.5);
                font = @".jpeg" ;
                [headerString appendFormat:@"Content-Disposition: form-data; name=\"%@\"; filename=\"%@%@\"\r\n",key,@"avatarImage",font];
                [headerString appendFormat:@"Content-Type: image/jpeg\r\n\r\n"];
            } else {
                //图片格式错误
                NSDictionary *dict = @{NSLocalizedDescriptionKey:@"UIImage格式错误"};
                NSError *myError = [NSError errorWithDomain:NSCustomErrorDomain code:9994 userInfo:dict];
                *stop = true;
                return !faliure?:faliure(myError);
            }
            //data拼接
            [data appendData:[headerString dataUsingEncoding:NSUTF8StringEncoding]];
            //image转的imageData
            [data appendData:imageData];
        } else {
            //不是图片
            [headerString appendFormat:@"Content-Disposition: form-data; name=\"%@\"\r\n\r\n",key];
            NSData *dataHeaderString =  [headerString dataUsingEncoding:NSUTF8StringEncoding] ;
            [data appendData:dataHeaderString];
            if ([obj isKindOfClass:[NSData class]]) {
                [data appendData:obj];
            } else if ([obj isKindOfClass:[NSString class]]){
                [data appendData:[obj dataUsingEncoding:NSUTF8StringEncoding]] ;
            } else {
                NSDictionary *dict = @{NSLocalizedDescriptionKey:@"对象未知类型错误"};
                NSError *myError = [NSError errorWithDomain:NSCustomErrorDomain code:9995 userInfo:dict];
                *stop = true;
                return !faliure?:faliure(myError);
            }
        }
    }];
    //尾部srr
    NSMutableString *footerString = [NSMutableString stringWithFormat:@"\r\n--%@--",demarcation];
    [data appendData:[footerString dataUsingEncoding:NSUTF8StringEncoding]];
    //HTTPBody赋值
    request.HTTPBody = data;
    //请求头设置
    [request setValue:[NSString stringWithFormat:@"multipart/form-data; boundary=%@",demarcation] forHTTPHeaderField:@"Content-Type"];
    [request setValue:[NSString stringWithFormat:@"%ld",data.length] forHTTPHeaderField:@"Content-Length"];
    [request setValue:@"http://token.wutnews.net/tools/image" forHTTPHeaderField:@"Referer"];
    //网络请求
    TokenNetworking.networking
    .request(^NSURLRequest *{
        return request;
    })
    .responseJSON(^(NSURLSessionTask *task, NSError *jsonError, id responsedObj) {
        if (jsonError) {
            !faliure?:faliure(jsonError);
        }
        if (responsedObj) {
            NSDictionary *json = responsedObj;
            if ([json isKindOfClass:[NSDictionary class]]) {
                //json是字典
                NSString *errorCode = json[@"error_code"];
                NSString *message = json[@"message"];
                if ([errorCode integerValue] == 0) {
                    //数据正确
                    !success?:success(message);
                } else {
                    //errorCode不为0
                    NSDictionary *dict = @{NSLocalizedDescriptionKey:@"TokenSession请求数据errorCode错误",
                                           NSLocalizedFailureReasonErrorKey:message};
                    NSError *myError = [NSError errorWithDomain:NSCustomErrorDomain code:9996 userInfo:dict];
                    !faliure?:faliure(myError);
                }
            } else {
                //返回数据格式不对
                NSDictionary *dict = @{NSLocalizedDescriptionKey:@"TokenSession请求数据解析错误"};
                NSError *myError = [NSError errorWithDomain:NSCustomErrorDomain code:9997 userInfo:dict];
                !faliure?:faliure(myError);
            }
        } else {
            //返回的data为空
            NSDictionary *dict = @{NSLocalizedDescriptionKey:@"TokenSession请求没有数据返回"};
            NSError *myError = [NSError errorWithDomain:NSCustomErrorDomain code:9998 userInfo:dict];
            !faliure?:faliure(myError);
        }
    })
    .failure(^(NSError *error) {
        !faliure?:faliure(error);
    });
}

@end
