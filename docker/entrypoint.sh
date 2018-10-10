#!/bin/bash
# Init Scripts
rm -rf /etc/localtime && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
mv -f /home/gopath/src/github.com/TokenTeam/Token-Static-Center/token-static-center /home/htdocs
mv -f /home/gopath/src/github.com/TokenTeam/Token-Static-Center/template /home/htdocs
mv -n /home/gopath/src/github.com/TokenTeam/Token-Static-Center/static /home/htdocs
chmod +x /home/htdocs/token-static-center
cd /home/htdocs/ && ./token-static-center --config=/etc/token-static-center/config.yaml
