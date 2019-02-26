@ECHO OFF

SET CONVERT="C:\Program Files\Git\usr\bin\dos2unix.exe"

ECHO Change Windows to Unix line endings
%CONVERT%  "..\init\01-initialize-database.sh"
%CONVERT%  "..\init\02-fill-database.sh"
%CONVERT%  "..\..\scripts\wait-for-it.sh"

docker-compose -f ..\..\docker-compose-test.yml build --no-cache
docker-compose -f ..\..\docker-compose-test.yml up