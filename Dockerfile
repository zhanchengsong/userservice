FROM golang:alpine
RUN mkdir /userservice
COPY ./ /userservice
WORKDIR /userservice
RUN go build -o userservice .
RUN adduser -S -D -H -h /userservice appuser
USER appuser
CMD ["./userservice"]
EXPOSE 3005/tcp