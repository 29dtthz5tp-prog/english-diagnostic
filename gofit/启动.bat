@echo off
chcp 65001 >nul
title GoFit · 英语诊断工具
cd /d "%~dp0"

echo.
echo   ╔══════════════════════════════════════════╗
echo   ║       GoFit · 英语诊断工具            ║
echo   ╚══════════════════════════════════════════╝
echo.
echo   启动本地服务 + Cloudflare 公网隧道...
echo.

REM 先杀旧进程
taskkill /f /im gofit.exe >nul 2>&1
taskkill /f /im cloudflared.exe >nul 2>&1

REM 启动 gofit
start "GoFit-Server" /MIN gofit.exe

REM 等服务器就绪
timeout /t 2 /nobreak >nul

REM 启动 cloudflared
start "Cloudflare-Tunnel" /MIN cloudflared.exe tunnel --url http://localhost:8877

REM 等隧道建立
timeout /t 4 /nobreak >nul

REM 从 cloudflared 日志提取公网链接
for /f "tokens=2 delims= " %%u in ('findstr /c:"trycloudflare.com" "%TEMP%\cf_output.txt" 2^>nul ^| findstr "https"') do set "CFURL=%%u"

echo.
echo   ╔══════════════════════════════════════════════╗
echo   ║  🌐 公网链接（发给学生）：                ║
echo   ║                                            ║
echo   ║  首页: 从 cloudflared 窗口复制链接        ║
echo   ║  学生端: 链接/学生端.html                   ║
echo   ║  教师端: 链接/教师端.html                   ║
echo   ║                                            ║
echo   ║  链接格式: https://xxx.trycloudflare.com   ║
echo   ║  ⚠️ 每次重启链接会变                     ║
echo   ║  ⚠️ 关闭本窗口 = 链接失效               ║
echo   ╚══════════════════════════════════════════════╝
echo.
pause
