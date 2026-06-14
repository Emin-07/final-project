FROM ubuntu:latest

ENV TODO_PORT=:7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=1234

WORKDIR /usr/src/app

COPY web/ ./web/
COPY scheduler .

EXPOSE 7540

CMD ["./scheduler"]
