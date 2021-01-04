FROM golang:1.14-alpine

WORKDIR /src/
COPY . .

RUN apk update && \
    apk add git && \
    go get github.com/cespare/reflex && \
    go get github.com/gorilla/mux && \
    go get github.com/go-sql-driver/mysql && \
    go build -o ./go_todo_server

EXPOSE 9999
# CMD ["reflex", "-c" "reflex.conf"]
CMD ./go_todo_server -p 9999