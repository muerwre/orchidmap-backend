FROM golang:1.14.1

WORKDIR /usr/src/app
COPY . .
RUN make build

WORKDIR /usr/src/app/bin
COPY ./wait-for-it.sh .
COPY ./config.yaml .
COPY ./views/* views/

EXPOSE 7777
CMD ["./orchidgo", "serve"]
