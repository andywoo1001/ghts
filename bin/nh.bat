@echo off

SET OLDPATH=%PATH%
SET GOARCH=386
SET GHTS_DIR=%GOPATH%\src\github.com\ghts\ghts
SET GCC_DIR=%GHTS_DIR%\3rd_party\ruby_devkit_32
SET BUILD_DEP_DIR=%GHTS_DIR%\3rd_party\build_dep
SET C_INCLUDE_PATH=%BUILD_DEP_DIR%\include
SET LIBRARY_PATH=%BUILD_DEP_DIR%\lib
SET NH_OpenAPI_DIR=%GHTS_DIR%\3rd_party\NH_OpenAPI
SET PATH=GHTS_DIR\bin;%GCC_DIR%\bin;%GCC_DIR%\mingw\bin;%BUILD_DEP_DIR%\bin;%NH_OpenAPI_DIR%;%PATH%

D:
cd %GHTS_DIR%\api\nh\internal

cls
go test

SET PATH=%OLDPATH%