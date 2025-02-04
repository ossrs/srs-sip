#!/bin/bash

BINARY_NAME="objs/srs-sip"
MAIN_PATH="main/main.go"
VUE_DIR="html/NextGB"
CONFIG_FILE="conf/config.yaml"

# 检测操作系统类型
case "$(uname -s)" in
    Darwin*)    
        echo "Mac OS X detected"
        ;;
    Linux*)     
        echo "Linux detected"
        ;;
    *)          
        echo "Unknown operating system"
        exit 1
        ;;
esac

build() {
    echo "Building Go binary..."
    mkdir -p objs
    go build -o ${BINARY_NAME} ${MAIN_PATH}
    
    echo "Copying config file..."
    if [ -f "${CONFIG_FILE}" ]; then
        mkdir -p "objs/$(dirname ${CONFIG_FILE})"
        cp -a "${CONFIG_FILE}" "objs/$(dirname ${CONFIG_FILE})/"
        echo "Config file copied to objs/$(dirname ${CONFIG_FILE})/"
    else
        echo "Warning: ${CONFIG_FILE} not found"
    fi
}

clean() {
    echo "Cleaning..."
    rm -rf ${BINARY_NAME}
    rm -rf ${VUE_DIR}/dist
    rm -rf ${VUE_DIR}/node_modules
    rm -rf objs/html
    rm -rf objs/${CONFIG_FILE}
}

run() {
    echo "Running application..."
    go build -o ${BINARY_NAME} ${MAIN_PATH}
    ./${BINARY_NAME}
}

vue_install() {
    echo "Installing Vue dependencies..."
    cd ${VUE_DIR}
    npm install
    cd ../..
}

vue_build() {
    echo "Building Vue project..."
    if [ ! -d "${VUE_DIR}" ]; then
        echo "Error: Vue directory not found at ${VUE_DIR}"
        return 1
    fi

    # Check Node.js version
    if ! command -v node &> /dev/null; then
        echo "Error: Node.js is not installed"
        return 1
    fi
    
    NODE_VERSION=$(node -v | cut -d "v" -f 2)
    NODE_MAJOR_VERSION=$(echo $NODE_VERSION | cut -d "." -f 1)
    if [ "$NODE_MAJOR_VERSION" -lt 17 ]; then
        echo "Error: Node.js version 17 or higher is required (current version: $NODE_VERSION)"
        echo "Please upgrade Node.js using your package manager or nvm"
        return 1
    fi

    pushd ${VUE_DIR} > /dev/null
    echo "Current directory: $(pwd)"
    
    if [ ! -f "package.json" ]; then
        echo "Error: package.json not found in ${VUE_DIR}"
        popd > /dev/null
        return 1
    fi

    # Check if node_modules exists and install dependencies if needed
    if [ ! -d "node_modules" ] || [ ! -f "node_modules/.package-lock.json" ]; then
        echo "Node modules not found or incomplete, installing dependencies..."
        npm install
        if [ $? -ne 0 ]; then
            echo "Error: Failed to install dependencies"
            popd > /dev/null
            return 1
        fi
    fi

    echo "Running npm run build..."
    npm run build
    if [ $? -ne 0 ]; then
        echo "Error: Vue build failed"
        popd > /dev/null
        return 1
    fi
    popd > /dev/null
    echo "Vue build completed successfully"

    echo "Copying dist files to objs directory..."
    rm -rf objs/html
    mkdir -p objs
    if [ ! -d "${VUE_DIR}/dist" ]; then
        echo "Error: Vue dist directory not found at ${VUE_DIR}/dist"
        return 1
    fi
    cp -r "${VUE_DIR}/dist" "objs/html"
    if [ $? -eq 0 ]; then
        echo "Vue dist files successfully copied to objs/html"
    else
        echo "Error copying files"
        return 1
    fi
}

vue_dev() {
    echo "Starting Vue development server..."
    cd ${VUE_DIR}
    npm run dev
    cd ../..
}

build_all() {
    clean
    build
    vue_build
}

# 根据命令行参数执行相应的功能
case "$1" in
    "build")
        build
        ;;
    "clean")
        clean
        ;;
    "run")
        run
        ;;
    "vue-install")
        vue_install
        ;;
    "vue-build")
        vue_build
        ;;
    "vue-dev")
        vue_dev
        ;;
    "all"|"")
        build_all
        ;;
    *)
        echo "Unknown command: $1"
        echo "Usage: $0 {build|clean|run|vue-install|vue-build|vue-dev|all}"
        exit 1
        ;;
esac 