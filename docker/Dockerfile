FROM centos
MAINTAINER LiuFuXin <loli@lurenjia.in>

COPY go1.11.linux-amd64.tar.gz /tmp
COPY ImageMagick-6.9.10-12.tar.gz /tmp
COPY vendor.zip /tmp
COPY entrypoint.sh /

RUN yum update -y \
    && yum reinstall -y glibc-common \
    && yum install -y telnet net-tools \
    && localedef -c -f UTF-8 -i zh_CN zh_CN.UTF-8 \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ENV LANG=zh_CN.UTF-8 \
    LANGUAGE=zh_CN:zh \
    LC_ALL=zh_CN.UTF-8

ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin
ENV LD_LIBRARY_PATH /usr/local/lib
ENV GOPATH /home/gopath
ENV PKG_CONFIG_PATH /usr/local/lib/pkgconfig

RUN yum install -y epel-release \
    && yum makecache \
    && yum install -y libjpeg-turbo libjpeg-turbo-devel libjpeg-turbo-utils libpng libpng-devel libpng-static libtiff libtiff-devel libtiff-static giflib giflib-devel giflib-utils libwebp libwebp-devel libwebp-tools make gcc gcc-g++ libtool pkgconfig zip unzip tar gzip freetype freetype-devel

RUN cd /tmp \
    && tar -zxvf /tmp/go1.11.linux-amd64.tar.gz \
    && tar -zxvf /tmp/ImageMagick-6.9.10-12.tar.gz \
    && mkdir /tmp/vendor \
    && mv /tmp/vendor.zip /tmp/vendor/ \
    && cd /tmp/vendor \
    && unzip vendor.zip \
    && rm -rf vendor.zip

RUN mv /tmp/go /usr/local/ \
    && mkdir -p /home/gopath/src \
    && mkdir -p /home/gopath/bin \
    && mkdir -p /home/gopath/pkg \
    && mkdir -p /home/htdocs \
    && mkdir -p /home/log 

RUN cd /tmp/ImageMagick6-6.9.10-12 \
    && ./configure --with-webp \
    && make -j 24 \
    && make install

RUN mkdir -p /home/gopath/src/github.com/TokenTeam/Token-Static-Center \
    && cd /home/gopath/src/github.com/TokenTeam/Token-Static-Center \
    && yum install -y git \
    && git clone https://github.com/TokenTeam/Token-Static-Center/ . \
    && mv /tmp/vendor/vendor .

RUN cd /home/gopath/src/github.com/TokenTeam/Token-Static-Center \
    && ls -l \
    && go build /home/gopath/src/github.com/TokenTeam/Token-Static-Center/init.go \
    && mv init token-static-center 

RUN chmod +x /entrypoint.sh \
    && rm -rf /tmp \
    && yum clean all

COPY config.yaml /etc/token-static-center/config.yaml

CMD /entrypoint.sh
