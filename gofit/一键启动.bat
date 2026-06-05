@echo off
chcp 65001 >nul
title GoFit · 永久链接模式
cd /d "%~dp0"

echo.
echo   ╔══════════════════════════════════════════╗
echo   ║    GoFit · 英语诊断工具（永久链接）    ║
echo   ╚══════════════════════════════════════════╝
echo.

REM ========== 1. 清理旧进程 ==========
echo   [1/3] 清理旧进程...
taskkill /f /im gofit.exe >nul 2>&1
taskkill /f /im ssh.exe >nul 2>&1
timeout /t 1 /nobreak >nul

REM ========== 2. 启动 gofit ==========
echo   [2/3] 启动本地服务...
start "GoFit-Server" /MIN gofit.exe
timeout /t 2 /nobreak >nul

REM ========== 3. 启动 SSH 隧道 ==========
echo   [3/3] 建立永久隧道...
start "Serveo-Tunnel" /MIN ssh -o StrictHostKeyChecking=no -o ServerAliveInterval=60 -o ServerAliveCountMax=3 -R gaotu-english:80:localhost:8877 serveo.net

timeout /t 3 /nobreak >nul

REM ========== 完成 ==========
echo.
echo   ╔══════════════════════════════════════════════════════╗
echo   ║  ✅ 启动成功！永久链接如下：                       ║
echo   ║                                                    ║
echo   ║  🧑‍🎓 学生端:                                     ║
echo   ║  https://gaotu-english.serveousercontent.com/s     ║
echo   ║                                                    ║
echo   ║  👩‍🏫 教师端:                                     ║
echo   ║  https://gaotu-english.serveousercontent.com/t     ║
echo   ║                                                    ║
echo   ║  ⚠️ 不要关闭本窗口，否则链接失效                ║
echo   ╚══════════════════════════════════════════════════════╝
echo.
pause
