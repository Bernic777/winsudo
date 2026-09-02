@echo off
echo ========================================
echo  WinSudo Uninstaller
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

echo [1/2] Removing winsudo.exe from System32...
del /f /q "C:\Windows\System32\winsudo.exe" >nul 2>&1
if exist "C:\Windows\System32\winsudo.exe" (
    echo [-] Gagal menghapus winsudo.exe
) else (
    echo [+] winsudo.exe berhasil dihapus
)

echo [2/2] Removing sudo.bat wrapper...
del /f /q "C:\Windows\System32\sudo.bat" >nul 2>&1
if exist "C:\Windows\System32\sudo.bat" (
    echo [-] Gagal menghapus sudo.bat
) else (
    echo [+] sudo.bat berhasil dihapus
)

echo.
echo ========================================
echo  Uninstalasi selesai!
echo ========================================
echo.
pause
