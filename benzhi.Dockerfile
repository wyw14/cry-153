FROM golang:1.26
ENV GOPROXY=off GOSUMDB=off
WORKDIR /src
COPY go.mod go.sum* ./
COPY vendor ./vendor
COPY . .
RUN go build -mod=vendor -o /riversentinel ./cmd/riversentinel
CMD ["/riversentinel"]
