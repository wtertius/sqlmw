ARG GOVER
FROM golang:${GOVER}

ARG GOVER

RUN mkdir -p ${GOPATH}/src/github.com/wtertius/sqlmw
COPY . ${GOPATH}/src/github.com/wtertius/sqlmw
WORKDIR ${GOPATH}/src/github.com/wtertius/sqlmw

ENV MSSQL_HOST mssql
ENV PG_HOST pg
ENV BOUNCER_HOST bouncer

ENTRYPOINT ["make", "test"]
