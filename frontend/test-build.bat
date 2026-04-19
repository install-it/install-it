@echo off
echo ================================
echo NUXT UI Migration - Build Test
echo ================================
echo.
cd /d "%~dp0"
echo Current directory: %CD%
echo.
echo Running type check...
call npm run type-check
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo [ERROR] Type check failed!
    pause
    exit /b 1
)
echo.
echo Type check passed!
echo.
echo Running build...
call npm run build
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo [ERROR] Build failed!
    pause
    exit /b 1
)
echo.
echo ================================
echo Build completed successfully!
echo ================================
pause
