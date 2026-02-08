FROM alpine:latest
COPY cultural .
RUN mkdir -p storage
ENV SERVER_ADDR=:3000
ENV STORAGE_LOCAL_PATH=/storage
EXPOSE 3000
CMD [ "cultural", "serve" ]
