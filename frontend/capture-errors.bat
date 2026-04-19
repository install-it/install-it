@echo off
echo Running type check and saving errors...
cd /d "%~dp0"
call npm run type-check > build-errors.txt 2>&1
echo.
echo Errors saved to build-errors.txt
echo.
echo Showing errors:
echo ================================
type build-errors.txt
echo ================================
pause
