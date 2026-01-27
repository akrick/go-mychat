#!/bin/bash

echo "=========================================="
echo "测试管理后台登录接口"
echo "=========================================="
echo ""

# 测试登录接口
echo "发送登录请求..."
curl -X POST http://localhost:8081/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  -v 2>&1 | grep -E "(< HTTP|< {\"code\"|token)"

echo ""
echo "=========================================="
echo "如果看到 HTTP 200 和 token，说明接口正常"
echo "如果看到 HTTP 401，说明认证失败"
echo "=========================================="
