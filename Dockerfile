FROM golang:1.19-alpine as builder 
RUN mkdir /app
WORKDIR /app
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o test . 

FROM alpine
WORKDIR /
COPY --from=builder /app/test-descriptor.yaml .
COPY --from=builder /app/test .
CMD [ "/test" ]
