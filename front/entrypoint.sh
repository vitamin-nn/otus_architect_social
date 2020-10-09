#!/bin/sh
text="VUE_APP_HTTP_SERVER_URL=${VUE_APP_HTTP_SERVER_URL}"
echo $text > /usr/share/nginx/html/.env
echo "Starting Nginx"
nginx -g 'daemon off;'