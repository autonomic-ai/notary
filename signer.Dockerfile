FROM golang:1.14.1-alpine

RUN apk add --update git gcc libc-dev

RUN addgroup --system notary && adduser --system --ingroup notary --uid 999 notary
USER notary

ENV GO111MODULE=on

ARG MIGRATE_VER=v4.6.2
RUN go get -tags 'mysql postgres file' github.com/golang-migrate/migrate/v4/cli@${MIGRATE_VER} && mv /go/bin/cli /go/bin/migrate

ENV GOFLAGS=-mod=vendor
ENV NOTARYPKG github.com/autonomic-ai/notary

# Copy the local repo to the expected go path
COPY . /go/src/${NOTARYPKG}

WORKDIR /go/src/${NOTARYPKG}

ENV SERVICE_NAME=notary_signer
ENV NOTARY_SIGNER_DEFAULT_ALIAS="timestamp_1"
ENV NOTARY_SIGNER_TIMESTAMP_1="testpassword"

# Install notary-signer
RUN go install \
    -tags pkcs11 \
    -ldflags "-w -X ${NOTARYPKG}/version.GitCommit=`git rev-parse --short HEAD` -X ${NOTARYPKG}/version.NotaryVersion=`cat NOTARY_VERSION`" \
    ${NOTARYPKG}/cmd/notary-signer

USER root
RUN apk del git gcc libc-dev && rm -rf /var/cache/apk/*

USER notary
ENTRYPOINT [ "notary-signer" ]
CMD [ "-config=fixtures/signer-config-local.json" ]
