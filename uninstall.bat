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

echo [1/2] Removing wsudo.exe from System32...
del /f /q "C:\Windows\System32\wsudo.exe" >nul 2>&1
if exist "C:\Windows\System32\wsudo.exe" (
    echo [-] Gagal menghapus wsudo.exe
) else (
    echo [+] wsudo.exe berhasil dihapus
)

echo [2/2] Removing wsudo.bat wrapper...
del /f /q "C:\Windows\System32\wsudo.bat" >nul 2>&1
if exist "C:\Windows\System32\wsudo.bat" (
    echo [-] Gagal menghapus wsudo.bat
) else (
    echo [+] wsudo.bat berhasil dihapus
)

echo.
echo ========================================
echo  Uninstalasi selesai!
echo ========================================
echo.
pause
