ARG GO_VERSION=1.26.10
ARG ALPINE_VERSION=3.23

ARG BASE_IMAGE_PREFIX=docker.io
ARG GOPROXY
ARG HTTP_PROXY
ARG HTTPS_PROXY
ARG LD_FLAGS
ARG NODE_IMAGE_VERSION=22

##########################################

FROM ${BASE_IMAGE_PREFIX}/library/golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

WORKDIR /workspace

ARG GOPROXY
ENV GOPROXY=${GOPROXY}
# HTTP_PROXY/HTTPS_PROXY are predefined Docker build args: declaring them as
# ARG makes them available as env vars during RUN, without being persisted
# into the image config.
ARG HTTP_PROXY
ARG HTTPS_PROXY

RUN --mount=type=cache,target=/var/cache/apk \
    apk add gcc musl-dev

COPY . .
