FROM alpine:latest
WORKDIR /app
COPY cultural .
RUN apk add tzdata
RUN mkdir -p storage
RUN chmod +x ./cultural
ENV SERVER_ADDR=:3000
ENV STORAGE_LOCAL_PATH=/app/storage
EXPOSE 3000
CMD [ "./cultural", "serve" ]
