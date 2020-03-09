FROM golang AS build-env

ADD . /app

RUN cd /app && \
    go build -o test

FROM tiny
COPY --from=build-env /app/test /test
ENTRYPOINT ["/test"]