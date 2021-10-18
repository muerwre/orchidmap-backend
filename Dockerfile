FROM golang:1.14.1

WORKDIR /usr/src/app
COPY . .
RUN make build

WORKDIR /usr/src/app/bin
COPY ./wait-for-it.sh .
COPY ./views/* views/

EXPOSE 8080
HEALTHCHECK --interval=1m --timeout=30s CMD curl -f http://app:$EXPOSE/api/route/list/all || kill -s 2 1
CMD ["./orchidgo", "serve"]
