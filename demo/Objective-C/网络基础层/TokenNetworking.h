//
//  TokenNetworking.h
//  NewHybrid
//
//  Created by 陈雄&武嘉晟 on 2018/6/11.
//  Copyright © 2018年 com.feelings. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "TokenNetworkingCategories.h"

/** 使用示范
 NSString *loginURL = @"https://www.xxx.com";
 NSDictionary *parameter = @{
 @"userName":@"xxx",
 @"password":@"xxx"
 };
 
 TokenNetworking.networking
 .postWithURL(loginURL, parameter)
 .responseJSON(^(NSURLSessionTask *task, NSError *jsonError,id responsedObj) {
 if (jsonError) {
 NSLog(@"json parse error %@",jsonError);
 }
 if (responsedObj) {
 NSLog(@"json = %@",responsedObj);
 }
 })
 // you can get Text at the same time
 .responseText(^(NSURLSessionTask *task, NSString *responsedText) {
 NSLog(@"responsedText = %@",responsedText);
 })
 //you can send another request,可以保证上一个请求处理结束才开始下一个请求
 .postWithURL(loginURL1, parameter1)
 .responseJSON(^(NSURLSessionTask *task, NSError *jsonError,id responsedObj) {
 if (jsonError) {
 NSLog(@"json parse error %@",jsonError);
 }
 if (responsedObj) {
 NSLog(@"json = %@",responsedObj);
 }
 })
 .failure(^(NSError *error) {
 NSLog(@"error = %@",error);
 });
 
 //custom request
 
 //creat a Request
 NSDictionary *parameter = @{
 @"userName":@"xxx",
 @"password":@"xxx"
 };
 NSURL *url = [NSURL URLWithString:@"http://www.xxx.com"];
 NSMutableURLRequest *request = NSMutableURLRequest.token_requestWithURL(url)
 .token_setMethod(@"POST")
 .token_setTimeout(30);
 
 //send the Request
 TokenNetworking.networking
 .request(^NSURLRequest *{
 return request;
 })
 .responseText(^(NSURLSessionTask *task, NSString *responsedText) {
 NSLog(@"%@",responsedText);
 })
 .responseJSON(^(NSURLSessionTask *task, NSError *jsonError,id responsedObj) {
 if (jsonError) {
 NSLog(@"json parse error %@",jsonError);
 }
 if (responsedObj) {
 NSLog(@"json = %@",responsedObj);
 }
 })
 .failure(^(NSError *error) {
 NSLog(@"error = %@",error);
 });
 */

@class TokenNetworking;

typedef NSString *(^TokenNetworkingGetStringBlock)(void);

//send
typedef NSURLRequest    *(^TokenRequestMakeBlock)(void);
typedef TokenNetworking *(^TokenSendRequestBlock)(TokenRequestMakeBlock make);
typedef TokenNetworking *(^TokenNetParametersBlock)(NSString *urlString,NSDictionary *parameters);

//redirect
typedef NSURLRequest    *(^TokenChainRedirectParameterBlock)(NSURLRequest *request,NSURLResponse *response);
typedef TokenNetworking *(^TokenChainRedirectBlock)(TokenChainRedirectParameterBlock redirectParameter);

//JSON TEXT FAILURE参数BLOCK
typedef void(^TokenNetSuccessJSONBlock)(NSURLSessionTask *task,NSError *jsonError,id responsedObj);
typedef void(^TokenNetSuccessTextBlock)(NSURLSessionTask *task,NSString *responsedText);
typedef void(^TokenNetFailureParameterBlock)(NSError *error);

//response
typedef TokenNetworking *(^TokenResponseJSONBlock)(TokenNetSuccessJSONBlock jsonBlock);
typedef TokenNetworking *(^TokenResponseTextBlock)(TokenNetSuccessTextBlock textBlock);

//willFailure
typedef TokenNetworking *(^TokenWillFailureBlock)(TokenNetFailureParameterBlock failureBlock);

//失败BLOCK
typedef TokenNetworking *(^TokenNetFailureBlock)(TokenNetFailureParameterBlock failure);

@interface TokenNetworking : NSObject

//初始化方法
+ (instancetype)networking;
@property (nonatomic, copy, readonly, class) TokenNetworkingGetStringBlock randomUA;
@property (nonatomic, copy, readonly, class) TokenNetworkingGetStringBlock defaultUA;
@end

@interface TokenNetworking(Chain)

//链式调用的基础
@property (nonatomic, copy, readonly) TokenNetParametersBlock getWithURL;
@property (nonatomic, copy, readonly) TokenNetParametersBlock postWithURL;
@property (nonatomic, copy, readonly) TokenSendRequestBlock   request;
@property (nonatomic, copy, readonly) TokenChainRedirectBlock willRedict;
@property (nonatomic, copy, readonly) TokenResponseJSONBlock  willResponseJSON;
@property (nonatomic, copy, readonly) TokenResponseTextBlock  willResponseText;
@property (nonatomic, copy, readonly) TokenResponseJSONBlock  responseJSON;
@property (nonatomic, copy, readonly) TokenResponseTextBlock  responseText;
@property (nonatomic, copy, readonly) TokenWillFailureBlock   willFailure;
@property (nonatomic, copy, readonly) TokenNetFailureBlock    failure;

@end
