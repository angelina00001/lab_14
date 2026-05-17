@echo off
REM Запускайте ЭТОТ файл (или команду ниже), а не двойной клик по .ps1
cd /d "%~dp0.."
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0run-tests.ps1"
exit /b %ERRORLEVEL%
