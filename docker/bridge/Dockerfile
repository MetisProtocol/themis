FROM golang:latest

ARG THEMIS_DIR=/var/lib/themis
ENV THEMIS_DIR=$THEMIS_DIR

RUN apt-get update -y && apt-get upgrade -y \
    && apt install build-essential git -y \
    && mkdir -p $THEMIS_DIR

WORKDIR ${THEMIS_DIR}
COPY . .

RUN make install

COPY docker/entrypoint.sh /usr/local/bin/entrypoint.sh

ENV SHELL /bin/bash
EXPOSE 1317 26656 26657

ENTRYPOINT ["entrypoint.sh"]
