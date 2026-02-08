FROM alpine:latest
COPY cultural .
ENV SERVER_ADDR=:3000
EXPOSE 3000
CMD [ "cultural", "serve" ]
