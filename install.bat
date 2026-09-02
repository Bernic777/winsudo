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

echo [1/2] Copying wsudo.exe to System32...
copy /Y "%~dp0wsudo.exe" "C:\Windows\System32\wsudo.exe" >nul
if %errorLevel% equ 0 (
    echo [+] wsudo.exe berhasil diinstall ke System32
) else (
    echo [-] Gagal copy ke System32
    pause
    exit /b 1
)

echo [2/2] Creating wsudo.bat wrapper...
echo @echo off > "C:\Windows\System32\wsudo.bat"
echo wsudo.exe %%* >> "C:\Windows\System32\wsudo.bat"
echo [+] wsudo.bat wrapper dibuat

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
