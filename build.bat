@echo off
setlocal

set BINARY_NAME=objs\srs-sip.exe
set MAIN_PATH=main\main.go
set VUE_DIR=html\NextGB
set CONFIG_FILE=conf\config.yaml

if "%1"=="" goto all
if "%1"=="build" goto build
if "%1"=="clean" goto clean
if "%1"=="run" goto run
if "%1"=="vue-install" goto vue-install
if "%1"=="vue-build" goto vue-build
if "%1"=="vue-dev" goto vue-dev
if "%1"=="all" goto all

:build
echo Building Go binary...
if not exist "objs" mkdir objs
go build -o %BINARY_NAME% %MAIN_PATH%

echo Copying config file...
if exist "%CONFIG_FILE%" (
    mkdir "objs\%~dp0%CONFIG_FILE%" 2>nul
    xcopy /s /i /y "%CONFIG_FILE%" "objs\%~dp0%CONFIG_FILE%\"
    echo Config file copied to objs\%~dp0%CONFIG_FILE%
) else (
    echo Warning: %CONFIG_FILE% not found
)
goto :eof

:clean
echo Cleaning...
if exist %BINARY_NAME% del /F /Q %BINARY_NAME%
if exist %VUE_DIR%\dist rd /S /Q %VUE_DIR%\dist
if exist %VUE_DIR%\node_modules rd /S /Q %VUE_DIR%\node_modules
if exist objs\html rd /S /Q objs\html
if exist objs\%CONFIG_FILE% del /F /Q objs\%CONFIG_FILE%
goto :eof

:run
echo Running application...
go build -o %BINARY_NAME% %MAIN_PATH%
%BINARY_NAME%
goto :eof

:vue-install
echo Installing Vue dependencies...
cd %VUE_DIR%
call npm install
cd ..\..
goto :eof

:vue-build
echo Building Vue project...
if not exist "%VUE_DIR%" (
    echo Error: Vue directory not found at %VUE_DIR%
    goto :eof
)

rem Check Node.js version
where node >nul 2>nul
if errorlevel 1 (
    echo Error: Node.js is not installed
    goto :eof
)

for /f "tokens=1,2,3 delims=." %%a in ('node -v') do (
    set NODE_MAJOR=%%a
)
set NODE_MAJOR=%NODE_MAJOR:~1%
if %NODE_MAJOR% LSS 17 (
    echo Error: Node.js version 17 or higher is required ^(current version: %NODE_MAJOR%^)
    echo Please upgrade Node.js using the official installer or nvm-windows
    goto :eof
)

pushd %VUE_DIR%
echo Current directory: %CD%
if not exist "package.json" (
    echo Error: package.json not found in %VUE_DIR%
    popd
    goto :eof
)

rem Check if node_modules exists and install dependencies if needed
if not exist "node_modules" (
    echo Node modules not found, installing dependencies...
    call npm install
    if errorlevel 1 (
        echo Error: Failed to install dependencies
        popd
        goto :eof
    )
)

echo Running npm run build...
call npm run build
if errorlevel 1 (
    echo Error: Vue build failed
    popd
    goto :eof
)
popd
echo Vue build completed successfully

echo Copying dist files to objs directory...
if exist objs\html rd /S /Q objs\html
if not exist objs mkdir objs
if not exist "%VUE_DIR%\dist" (
    echo Error: Vue dist directory not found at %VUE_DIR%\dist
    goto :eof
)
robocopy "%VUE_DIR%\dist" "objs\html" /E /NFL /NDL /NJH /NJS /nc /ns /np
if errorlevel 8 (
    echo Error copying files
) else (
    echo Vue dist files successfully copied to objs\html
)
goto :eof

:vue-dev
echo Starting Vue development server...
cd %VUE_DIR%
call npm run dev
cd ..\..
goto :eof

:all
echo Building entire project...
call :build
call :vue-build
echo.
echo Press any key to exit...
pause>nul
goto :eof 