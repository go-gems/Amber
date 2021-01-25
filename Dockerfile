FROM golang:alpine as builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o Amber AmberServer.go

WORKDIR /dist
RUN cp -R /build/public .
RUN cp /build/config.json .
RUN cp /build/Amber .

FROM scratch
COPY --from=builder /dist /app/

WORKDIR /app
EXPOSE 8080

CMD ["./Amber"]