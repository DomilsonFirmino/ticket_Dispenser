FROM golang:1.27

WORKDIR /app

COPY go.mod ./

#COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

EXPOSE 8080

CMD ["/docker-gs-ping"]

# docker build --tag docker-gs-ping .
# docker build -t docker-gs-ping:multistage -f Dockerfile.multistage .