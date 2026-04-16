FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache tzdata ffmpeg vips vips-tools

COPY cultural .
RUN chmod +x ./cultural

RUN mkdir -p storage
RUN mkdir -p data

ENV SERVER_ADDR=:3000
ENV STORAGE_LOCAL_PATH=/app/storage

EXPOSE 3000

CMD [ "./cultural", "serve" ]
