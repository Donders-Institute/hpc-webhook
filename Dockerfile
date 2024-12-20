# stage 0: compile go program
FROM golang:1.23-bullseye
RUN mkdir -p /tmp/hpc-webhook
WORKDIR /tmp/hpc-webhook
ADD cmd ./cmd
ADD internal ./internal
ADD pkg ./pkg
ADD go.mod .
ADD go.sum .
RUN ls -l /tmp/hpc-webhook && GOOS=linux go build -a -installsuffix cgo -o bin/server cmd/server/server.go

# stage 1
FROM almalinux:8 as server
RUN ulimit -n 1024 && yum install -y which nfs4-acl-tools sssd-client attr acl && yum clean all && rm -rf /var/cache/yum/*
WORKDIR /root
EXPOSE 5111
COPY --from=0 /tmp/hpc-webhook/bin/server .
COPY scripts/wait-for-it.sh .
RUN chmod +x wait-for-it.sh
# Wait and sleep for 30 sec before starting server
CMD ./wait-for-it.sh db:5432 --timeout=0 -- sleep 30 && echo "Started" && ./server
