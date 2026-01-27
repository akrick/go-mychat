#!/bin/bash

echo "========================================"
echo "   MyChat 管理后台启动脚本"
echo "========================================"
echo

# 检查后端依赖
echo "[1/3] 检查后端依赖..."
cd backend
if [ ! -f go.mod ]; then
    echo "[错误] go.mod 文件不存在！"
    exit 1
fi
go mod tidy
if [ $? -ne 0 ]; then
    echo "[错误] Go 依赖安装失败！"
    exit 1
fi
echo "[√] 后端依赖检查完成"
cd ..
echo

# 检查前端依赖
echo "[2/3] 检查前端依赖..."
cd frontend
if [ ! -d node_modules ]; then
    echo "[提示] 正在安装前端依赖，请稍候..."
    npm install
    if [ $? -ne 0 ]; then
        echo "[错误] 前端依赖安装失败！"
        exit 1
    fi
fi
echo "[√] 前端依赖检查完成"
cd ..
echo

# 启动服务
echo "[3/3] 启动服务..."
echo
echo "========================================"
echo "   选择启动模式"
echo "========================================"
echo "[1] 仅启动后端服务 (端口 8081)"
echo "[2] 仅启动前端服务 (端口 3000)"
echo "[3] 同时启动后端和前端"
echo "[0] 退出"
echo "========================================"
read -p "请选择 (0-3): " choice

case $choice in
    1)
        echo
        echo "[启动] 后端服务..."
        cd backend
        go run main.go &
        echo "[√] 后端服务已启动: http://localhost:8081"
        ;;
    2)
        echo
        echo "[启动] 前端服务..."
        cd frontend
        npm run dev &
        echo "[√] 前端服务已启动: http://localhost:3000"
        ;;
    3)
        echo
        echo "[启动] 后端服务..."
        cd backend
        go run main.go &
        BACKEND_PID=$!
        echo "[√] 后端服务已启动: http://localhost:8081"
        cd ..
        sleep 2
        echo
        echo "[启动] 前端服务..."
        cd frontend
        npm run dev &
        FRONTEND_PID=$!
        echo "[√] 前端服务已启动: http://localhost:3000"
        
        # 捕获退出信号
        trap "kill $BACKEND_PID $FRONTEND_PID 2>/dev/null" EXIT
        ;;
    0)
        echo
        echo "[退出] 已取消启动"
        exit 0
        ;;
    *)
        echo
        echo "[错误] 无效的选择！"
        exit 1
        ;;
esac

echo
echo "========================================"
echo "   启动完成"
echo "========================================"
echo "管理后台前端: http://localhost:3000"
echo "管理后台后端: http://localhost:8081"
echo "========================================"
echo
echo "默认账号: admin"
echo "默认密码: admin123"
echo

# 如果是模式3，等待进程
if [ "$choice" == "3" ]; then
    wait
fi
