FROM golang:1.19-alpine as builder 
RUN mkdir /app
WORKDIR /app
COPY ./ ./
RUN cd /app/server/ && go mod download && CGO_ENABLED=0 GOOS=linux go build -a -o quality-trace .

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /app/sample/test-descriptor.yaml /app/server/quality-trace ./
CMD [ "/quality-trace" ]
CMD "S"
CMD "o"
CMD "p"
CMD "j
CMD "L"
CMD "k"
CMD "l"
CMD "k"
CMD "j"
CMD "Y"
CMD "l"







