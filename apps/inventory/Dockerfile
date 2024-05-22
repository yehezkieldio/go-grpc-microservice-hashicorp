FROM golang:1.22.3

WORKDIR /app

COPY dist/apps/inventory /app/inventory

EXPOSE 50051

CMD ["/app/inventory"]