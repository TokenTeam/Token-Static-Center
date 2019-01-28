//
//  TokenSession.h
//  MyWHUT
//
//  Created by ly on 2019/1/21.
//  Copyright © 2019 com.feelings. All rights reserved.
//

#ifndef TokenSession_h
#define TokenSession_h

#import <Foundation/Foundation.h>

typedef void(^TokenSessionSuccessBlock)(id obj);
typedef void(^TokenSessionFailureBlock)(NSError *error);

@interface TokenSession : NSObject


+ (void)uploadDataTaskUseFormData:(NSString *)method
                              url:(NSString *)urlString
                             data:(NSDictionary *)dataDictionary
                          success:(TokenSessionSuccessBlock)success
                          faliure:(TokenSessionFailureBlock)faliure ;

@end

#endif /* TokenSession_h */
