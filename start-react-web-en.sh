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

# Check and start AKTools service
echo
echo "Checking AKTools service..."
if curl -s http://127.0.0.1:8080/health > /dev/null 2>&1; then
    echo "OK: AKTools service is already running (http://127.0.0.1:8080)"
else
    echo "AKTools service is not running. Starting it now..."
    
    # Check if Python is available
    if ! command -v python3 &> /dev/null && ! command -v python &> /dev/null; then
        echo "WARNING: Python not found. Please install Python 3.8+"
        echo "Continuing without AKTools (will use fallback data source)"
    else
        # Use python3 if available, otherwise python
        PYTHON_CMD="python3"
        if ! command -v python3 &> /dev/null; then
            PYTHON_CMD="python"
        fi
        
        # Check if AKTools is installed
        if ! $PYTHON_CMD -c "import aktools" > /dev/null 2>&1; then
            echo "WARNING: AKTools not installed. Installing now..."
            pip install aktools
            if [ $? -ne 0 ]; then
                echo "ERROR: Failed to install AKTools"
                echo "Continuing without AKTools (will use fallback data source)"
            else
                echo "AKTools installed successfully"
            fi
        fi
        
        # Try to start AKTools if it's available
        if $PYTHON_CMD -c "import aktools" > /dev/null 2>&1; then
            echo "Starting AKTools service in background..."
            nohup $PYTHON_CMD -m aktools --port 8080 > /dev/null 2>&1 &
            AKTOOLS_PID=$!
            
            # Wait for AKTools to start and become ready
            echo "Waiting for AKTools to start..."
            sleep 8
            
            # Check AKTools API with retry logic (same endpoint as Go backend tests)
            AKTOOLS_READY=0
            for i in {1..10}; do
                echo "Testing AKTools API endpoint... (attempt $i/10)"
                if curl -s "http://127.0.0.1:8080/api/public/stock_zh_a_hist" > /dev/null 2>&1; then
                    echo "OK: AKTools service is fully ready (http://127.0.0.1:8080)"
                    AKTOOLS_READY=1
                    break
                fi
                sleep 6
            done
            
            if [ $AKTOOLS_READY -eq 0 ]; then
                echo "WARNING: AKTools service failed to start properly"
                echo "The Go backend may fail to start if configured to use AKTools"
                echo "Please check AKTools installation: pip install aktools"
            fi
        fi
    fi
fi

# Check and start Go backend
echo
echo "Checking Go backend service..."
if curl -s http://localhost:8081/api/v1/health > /dev/null 2>&1; then
    echo "OK: Go backend is already running (http://localhost:8081)"
else
    echo "Go backend is not running. Starting it now..."
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        echo "ERROR: Go not found. Please install Go 1.22+"
        echo "Download: https://golang.org/dl/"
        read -p "Press any key to exit..."
        exit 1
    fi
    
    echo "Starting Go backend server in background..."
    nohup go run cmd/server/main.go > /dev/null 2>&1 &
    GO_PID=$!
    
    # Wait for server to start
    echo "Waiting for server to start..."
    sleep 8
    
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
echo "1. React Web App (browser-based)"
echo "2. React Native Mobile App (mobile simulator)"
echo "3. Both Web and Mobile"
echo
read -p "Enter your choice (1-3): " choice

case $choice in
    1)
        echo
        echo "Starting React Web application..."
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
        npm start
        ;;
    2)
        echo
        echo "Starting React Native mobile application..."
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
        echo "Use 'npm run android' or 'npm run ios' in another terminal to run on device/simulator"
        npm start
        ;;
    3)
        echo
        echo "Starting both Web and Mobile applications..."
        
        # Start Web App in background
        echo "Starting React Web App..."
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
        echo "Starting React Native Mobile App..."
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
        echo "- Web App: http://localhost:3000"
        echo "- Mobile Metro: http://localhost:8081 (Metro bundler)"
        echo "- Go Backend: http://localhost:8081 (API)"
        echo "- AKTools Service: http://127.0.0.1:8080 (Data provider)"
        echo
        echo "To run on mobile device/simulator, use:"
        echo "  npm run android  (Android)"
        echo "  npm run ios      (iOS)"
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
        echo "Invalid choice, starting React Web App..."
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
        npm start
        ;;
esac
