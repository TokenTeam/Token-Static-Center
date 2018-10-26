#!/bin/bash
# Init Scripts
rm -rf /etc/localtime && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
mv -f /home/gopath/src/github.com/TokenTeam/Token-Static-Center/token-static-center /bin
mv -f /home/gopath/src/github.com/TokenTeam/Token-Static-Center/template /home/htdocs
mv -n /home/gopath/src/github.com/TokenTeam/Token-Static-Center/static /home/htdocs
chmod +x /bin/token-static-center
cd /bin/token-static-center --config=/etc/token-static-center/config.yaml
