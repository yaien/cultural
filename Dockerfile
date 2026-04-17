FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache tzdata ffmpeg vips vips-tools

COPY cultural .
RUN chmod +x ./cultural

RUN mkdir -p data/certs data/db data/storage

ENV SERVER_ADDR=:3000
ENV STORAGE_PROVIDER=local
ENV STORAGE_LOCAL_PATH=/app/data/storage
ENV SQLITE_DSN=/app/data/db/data.db

EXPOSE 3000

CMD [ "./cultural", "serve" ]
