#!/bin/bash
# Stock-A-Future Full Stack Startup Script (English Version - macOS/Linux)
# Starts Go backend and React Native/Web frontend

echo "Starting Stock-A-Future Full Stack Application..."
echo
echo "Service startup order:"
echo "1. AKTools Service (Data Provider) - Port 8080"
echo "2. Go Backend API - Port 8081"
echo "3. Frontend Application - Port 3000"
echo

# Check Node.js
if ! command -v node &> /dev/null; then
    echo "ERROR: Node.js not found. Please install Node.js 18+"
    echo "Download: https://nodejs.org/"
    read -p "Press any key to exit..."
    exit 1
else
    NODE_VERSION=$(node --version)
    echo "OK: Node.js version: $NODE_VERSION"
fi

# Start AKTools service
echo
echo "=========================================="
echo "STEP 1: Starting AKTools service..."
echo "=========================================="
echo "Starting AKTools on port 8080..."

# Use python3 if available, otherwise python
PYTHON_CMD="python3"
if ! command -v python3 &> /dev/null; then
    PYTHON_CMD="python"
fi

nohup $PYTHON_CMD -m aktools --port 8080 > /dev/null 2>&1 &
AKTOOLS_PID=$!

echo "Waiting 8 seconds for AKTools to initialize..."
sleep 8
echo "AKTools should now be running on http://127.0.0.1:8080"

# Check Go installation
echo
echo "=========================================="
echo "STEP 2: Checking Go installation..."
echo "=========================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "ERROR: Go not found. Please install Go 1.22+"
    echo "Download: https://golang.org/dl/"
    read -p "Press any key to exit..."
    exit 1
else
    echo "OK: Go is available"
fi

# Start Go backend directly with go run
echo
echo "=========================================="
echo "STEP 3: Starting Go backend service..."
echo "=========================================="
if curl -s http://localhost:8081/api/v1/health > /dev/null 2>&1; then
    echo "OK: Go backend is already running (http://localhost:8081)"
else
    echo "Go backend is not running. Starting it now..."
    
    echo "Starting Go backend server in background..."
    echo "Command: SERVER_PORT=8081 go run cmd/server/main.go"
    nohup env SERVER_PORT=8081 go run cmd/server/main.go > /dev/null 2>&1 &
    GO_PID=$!
    
    # Wait for server to start
    echo "Waiting for server to start..."
    sleep 10
    
    # Check if server started successfully
    if curl -s http://localhost:8081/api/v1/health > /dev/null 2>&1; then
        echo "OK: Go backend started successfully (http://localhost:8081)"
    else
        echo "WARNING: Go backend may still be starting..."
        echo "If the web app fails to load data, please wait a moment and refresh."
    fi
fi

echo
echo "Choose frontend to start:"
echo "1. React Web App (Development Mode - Hot Reload)"
echo "2. React Web App (Production Mode - Faster Startup)"
echo "3. React Native Mobile App (Development Mode)"
echo "4. React Native Mobile App (Production Mode)"
echo "5. Both Web and Mobile (Development Mode)"
echo "6. Both Web and Mobile (Production Mode)"
echo
read -p "Enter your choice (1-6): " choice

case $choice in
    1)
        echo
        echo "Starting React Web application (Development Mode)..."
        cd web-react
        if [ ! -d "node_modules" ]; then
            echo "Installing web dependencies..."
            npm install
            if [ $? -ne 0 ]; then
                echo "ERROR: Failed to install web dependencies"
                read -p "Press any key to exit..."
                exit 1
            fi
        fi
        echo "Web app will open in browser: http://localhost:3000"
        echo "Development mode: Hot reload enabled, slower startup"
        npm start
        ;;
    2)
        echo
        echo "Starting React Web application (Production Mode)..."
        cd web-react
        if [ ! -d "node_modules" ]; then
            echo "Installing web dependencies..."
            npm install
            if [ $? -ne 0 ]; then
                echo "ERROR: Failed to install web dependencies"
                read -p "Press any key to exit..."
                exit 1
            fi
        fi
        echo "Building production version..."
        npm run build
        if [ $? -ne 0 ]; then
            echo "ERROR: Failed to build web application"
            read -p "Press any key to exit..."
            exit 1
        fi
        echo "Web app will open in browser: http://localhost:3000"
        echo "Production mode: Optimized build, faster startup, no hot reload"
        npm run serve
        ;;
    3)
        echo
        echo "Starting React Native mobile application (Development Mode)..."
        cd mobile
        if [ ! -d "node_modules" ]; then
            echo "Installing mobile dependencies..."
            npm install
            if [ $? -ne 0 ]; then
                echo "ERROR: Failed to install mobile dependencies"
                read -p "Press any key to exit..."
                exit 1
            fi
        fi
        echo "Starting React Native Metro bundler..."
        echo "Development mode: Hot reload enabled, slower startup"
        echo "Use 'npm run android' or 'npm run ios' in another terminal to run on device/simulator"
        npm start
        ;;
    4)
        echo
        echo "Starting React Native mobile application (Production Mode)..."
        cd mobile
        if [ ! -d "node_modules" ]; then
            echo "Installing mobile dependencies..."
            npm install
            if [ $? -ne 0 ]; then
                echo "ERROR: Failed to install mobile dependencies"
                read -p "Press any key to exit..."
                exit 1
            fi
        fi
        echo "Building production bundles..."
        echo "Creating Android bundle..."
        npm run bundle:android
        echo "Creating iOS bundle..."
        npm run bundle:ios
        echo "Production mode: Optimized bundles, faster app performance"
        echo "Use 'npm run android:release' or 'npm run ios:release' to run production builds"
        echo "Starting Metro bundler for additional development..."
        npm start
        ;;
    5)
        echo
        echo "Starting both Web and Mobile applications (Development Mode)..."
        
        # Start Web App in background
        echo "Starting React Web App (Development)..."
        cd web-react
        if [ ! -d "node_modules" ]; then
            echo "Installing web dependencies..."
            npm install
            if [ $? -ne 0 ]; then
                echo "ERROR: Failed to install web dependencies"
                read -p "Press any key to exit..."
                exit 1
            fi
        fi
        
        # Start web app in background
        npm start > /dev/null 2>&1 &
        WEB_PID=$!
        cd ..
        
        # Start Mobile App
        echo "Starting React Native Mobile App (Development)..."
        cd mobile
        if [ ! -d "node_modules" ]; then
            echo "Installing mobile dependencies..."
            npm install
            if [ $? -ne 0 ]; then
                echo "ERROR: Failed to install mobile dependencies"
                read -p "Press any key to exit..."
                exit 1
            fi
        fi
        
        echo
        echo "All services are starting:"
        echo "- Web App: http://localhost:3000 (Development Mode)"
        echo "- Mobile Metro: http://localhost:8081 (Metro bundler - Development)"
        echo "- Go Backend: http://localhost:8081 (API)"
        echo "- AKTools Service: http://127.0.0.1:8080 (Data provider)"
        echo
        echo "To run on mobile device/simulator, use:"
        echo "  npm run android  (Android - Development)"
        echo "  npm run ios      (iOS - Development)"
        echo
        echo "Press Ctrl+C to stop all services"
        echo
        
        # Function to cleanup background processes
        cleanup() {
            echo
            echo "Stopping services..."
            if [ ! -z "$WEB_PID" ]; then
                kill $WEB_PID 2>/dev/null
            fi
            if [ ! -z "$GO_PID" ]; then
                kill $GO_PID 2>/dev/null
            fi
            if [ ! -z "$AKTOOLS_PID" ]; then
                kill $AKTOOLS_PID 2>/dev/null
            fi
            exit 0
        }
        
        # Set trap to cleanup on script exit
        trap cleanup INT TERM
        
        npm start
        ;;
    6)
        echo
        echo "Starting both Web and Mobile applications (Production Mode)..."
        
        # Start Web App in background
        echo "Building and starting React Web App (Production)..."
        cd web-react
        if [ ! -d "node_modules" ]; then
            echo "Installing web dependencies..."
            npm install
            if [ $? -ne 0 ]; then
                echo "ERROR: Failed to install web dependencies"
                read -p "Press any key to exit..."
                exit 1
            fi
        fi
        
        echo "Building production web app..."
        npm run build
        if [ $? -ne 0 ]; then
            echo "ERROR: Failed to build web application"
            read -p "Press any key to exit..."
            exit 1
        fi
        
        # Start web app in background
        npm run serve > /dev/null 2>&1 &
        WEB_PID=$!
        cd ..
        
        # Start Mobile App
        echo "Building React Native Mobile App (Production)..."
        cd mobile
        if [ ! -d "node_modules" ]; then
            echo "Installing mobile dependencies..."
            npm install
            if [ $? -ne 0 ]; then
                echo "ERROR: Failed to install mobile dependencies"
                read -p "Press any key to exit..."
                exit 1
            fi
        fi
        
        echo "Building production mobile bundles..."
        npm run bundle:android
        npm run bundle:ios
        
        echo
        echo "All services are starting:"
        echo "- Web App: http://localhost:3000 (Production Mode)"
        echo "- Mobile Metro: http://localhost:8081 (Metro bundler)"
        echo "- Go Backend: http://localhost:8081 (API)"
        echo "- AKTools Service: http://127.0.0.1:8080 (Data provider)"
        echo
        echo "To run on mobile device/simulator, use:"
        echo "  npm run android:release  (Android - Production)"
        echo "  npm run ios:release      (iOS - Production)"
        echo
        echo "Press Ctrl+C to stop all services"
        echo
        
        # Function to cleanup background processes
        cleanup() {
            echo
            echo "Stopping services..."
            if [ ! -z "$WEB_PID" ]; then
                kill $WEB_PID 2>/dev/null
            fi
            if [ ! -z "$GO_PID" ]; then
                kill $GO_PID 2>/dev/null
            fi
            if [ ! -z "$AKTOOLS_PID" ]; then
                kill $AKTOOLS_PID 2>/dev/null
            fi
            exit 0
        }
        
        # Set trap to cleanup on script exit
        trap cleanup INT TERM
        
        npm start
        ;;
    *)
        echo "Invalid choice, starting React Web App (Development Mode)..."
        cd web-react
        if [ ! -d "node_modules" ]; then
            echo "Installing web dependencies..."
            npm install
            if [ $? -ne 0 ]; then
                echo "ERROR: Failed to install web dependencies"
                read -p "Press any key to exit..."
                exit 1
            fi
        fi
        echo "Web app will open in browser: http://localhost:3000"
        echo "Development mode: Hot reload enabled, slower startup"
        npm start
        ;;
esac
