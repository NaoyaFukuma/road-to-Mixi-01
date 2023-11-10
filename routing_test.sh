#!/bin/bash

# 実行例: 
# > chmod +x routing_test.sh
# > ./routing_test.sh 2> test_log.txt

# ベースURL
BASE_URL="http://localhost:8080"

# ルートのテスト
echo "Testing root..."
curl -I -v "$BASE_URL/"

echo "Testing non-existent file..."
curl -I -v "$BASE_URL/nonexistent"

echo "Testing img/image1.png..."
curl -I -v "$BASE_URL/img/image1.png"

echo "Testing img/image2.png..."
curl -I -v "$BASE_URL/img/image2.png"

echo "Testing /test..."
curl -I -v "$BASE_URL/test"
