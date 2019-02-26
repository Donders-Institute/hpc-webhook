@ECHO OFF

SET CONVERT="C:\Program Files\Git\usr\bin\dos2unix.exe"

ECHO Change Windows to Unix line endings
%CONVERT% "..\init\01-initialize-database.sh"
%CONVERT% "..\cmd\client\wait-for-it.sh"
%CONVERT% "..\cmd\server\wait-for-it.sh"

docker-compose -f ..\deployments\docker-compose.yml build
docker-compose -f ..\deployments\docker-compose.yml up