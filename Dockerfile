FROM golang:1.19-alpine as builder 
RUN mkdir /app
WORKDIR /app
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o quality-trace . 

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /app/test-descriptor.yaml .
COPY --from=builder /app/quality-trace .
CMD [ "/quality-trace" ]
