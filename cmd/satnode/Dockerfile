FROM golang:1.25-alpine
WORKDIR /app
COPY . .
RUN go build -o satnode ./cmd/satnode
EXPOSE 8080
CMD ["./satnode"]
