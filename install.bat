@echo off
echo ========================================
echo  WinSudo Installer
echo ========================================
echo.

:: Check admin
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo [!] Jalankan script ini sebagai Administrator!
    echo     Klik kanan - Run as administrator
    pause
    exit /b 1
)

echo [1/2] Copying sudo.exe to System32...
copy /Y "%~dp0sudo.exe" "C:\Windows\System32\winsudo.exe" >nul
if %errorLevel% equ 0 (
    echo [+] sudo.exe berhasil diinstall ke System32
) else (
    echo [-] Gagal copy ke System32
    pause
    exit /b 1
)

echo [2/2] Creating sudo.bat wrapper...
echo @echo off > "C:\Windows\System32\sudo.bat"
echo winsudo.exe %%* >> "C:\Windows\System32\sudo.bat"
echo [+] sudo.bat wrapper dibuat

echo.
echo ========================================
echo  Instalasi selesai!
echo ========================================
echo.
echo Sekarang kamu bisa menjalankan:
echo   sudo cmd
echo   sudo powershell
echo   sudo --help
echo.
pause
