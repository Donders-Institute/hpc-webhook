FROM golang:1.12-rc-stretch

COPY ../ /go/src/github.com/Donders-Institute/hpc-qaas
WORKDIR /go/src/github.com/Donders-Institute/hpc-qaas/cmd/server

RUN go get ./
RUN go build

EXPOSE 4444

CMD wait-for-it.sh db:5432 && echo "Started" && ./qaas-server