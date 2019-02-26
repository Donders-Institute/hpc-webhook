@ECHO OFF

SET CONVERT="C:\Program Files\Git\usr\bin\dos2unix.exe"

ECHO Change Windows to Unix line endings
%CONVERT% "..\init\01-initialize-database.sh"
%CONVERT% "wait-for-it.sh"

docker-compose -f ..\docker-compose.yml build
docker-compose -f ..\docker-compose.yml up