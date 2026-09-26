FROM golang:1.27-alpine

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

EXPOSE 3000

CMD ["go", "run", "."]
