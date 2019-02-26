FROM golang:1.12-rc-stretch

RUN go get -u github.com/lib/pq
RUN adduser --disabled-password --gecos '' dccnuser

ENV SRC_DIR=/go/src/github.com/Donders-Institute/hpc-qaas/cmd/client
RUN mkdir -p $SRC_DIR

ADD . $SRC_DIR
RUN cd $SRC_DIR; go build -o qaas; cp qaas /usr/local/bin;

WORKDIR /home/dccnuser
USER dccnuser

CMD $SRC_DIR/wait-for-it.sh db:5432 && echo "Started" && tail -f /dev/null