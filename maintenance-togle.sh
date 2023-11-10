#!/bin/bash
# maintenance-toggle.sh

CONTAINER_NAME=reverse_proxy
NGINX_ROOT=/usr/share/nginx/html

# メンテナンスフラグファイルのパス
MAINTENANCE_FLAG="$NGINX_ROOT/file/maintenance.flag"

# コンテナ内でメンテナンスフラグファイルの存在をチェック
if docker exec $CONTAINER_NAME [ -f $MAINTENANCE_FLAG ]; then
    # フラグファイルが存在する場合は削除してメンテナンスモードを無効にする
    echo "Disabling maintenance mode..."
    docker exec $CONTAINER_NAME rm -f $MAINTENANCE_FLAG
else
    # フラグファイルが存在しない場合は作成してメンテナンスモードを有効にする
    echo "Enabling maintenance mode..."
    docker exec $CONTAINER_NAME touch $MAINTENANCE_FLAG
fi