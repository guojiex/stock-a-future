#!/bin/bash

# Stock-A-Future Mobile App 安装脚本

echo "🚀 开始设置 Stock-A-Future Mobile App..."

# 检查Node.js版本
node_version=$(node -v | cut -d'v' -f2)
required_version="18.0.0"

if [ "$(printf '%s\n' "$required_version" "$node_version" | sort -V | head -n1)" = "$required_version" ]; then 
    echo "✅ Node.js版本检查通过: $node_version"
else
    echo "❌ 需要Node.js 18+，当前版本: $node_version"
    exit 1
fi

# 安装依赖
echo "📦 安装依赖包..."
npm install

if [ $? -eq 0 ]; then
    echo "✅ 依赖安装成功"
else
    echo "❌ 依赖安装失败"
    exit 1
fi

# iOS设置 (仅在macOS上)
if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "🍎 检测到macOS，设置iOS环境..."
    if command -v pod &> /dev/null; then
        cd ios
        pod install
        cd ..
        echo "✅ iOS Pods安装成功"
    else
        echo "⚠️  未找到CocoaPods，请手动安装: sudo gem install cocoapods"
    fi
fi

# 检查Android环境
if [ -d "$ANDROID_HOME" ]; then
    echo "✅ 检测到Android SDK: $ANDROID_HOME"
else
    echo "⚠️  未检测到ANDROID_HOME环境变量"
    echo "   请确保已安装Android Studio并设置环境变量"
fi

echo ""
echo "🎉 设置完成！"
echo ""
echo "📱 运行应用："
echo "   npm run android  # 运行Android版本"
if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "   npm run ios      # 运行iOS版本"
fi
echo "   npm start        # 启动Metro bundler"
echo ""
echo "🔧 开发工具："
echo "   npm run type-check  # TypeScript类型检查"
echo "   npm run lint        # ESLint代码检查"
echo ""
echo "📖 更多信息请查看 README.md"
