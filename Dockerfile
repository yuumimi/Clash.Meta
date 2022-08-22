FROM golang:alpine as builder

RUN apk add --no-cache make git && \
    mkdir /clash-config && \
    wget -O /clash-config/Country.mmdb https://raw.githubusercontent.com/yuumimi/rules/release/clash/Country.mmdb && \
    wget -O /clash-config/geosite.dat https://raw.githubusercontent.com/yuumimi/rules/release/clash/geosite.dat && \
    wget -O /clash-config/geoip.dat https://raw.githubusercontent.com/yuumimi/rules/release/clash/geoip.dat


COPY . /clash-src
WORKDIR /clash-src
RUN go mod download &&\
    make docker &&\
    mv ./bin/Clash.Meta-docker /clash

FROM alpine:latest
LABEL org.opencontainers.image.source="https://github.com/MetaCubeX/Clash.Meta"

RUN apk add --no-cache ca-certificates tzdata

VOLUME ["/root/.config/clash/"]

COPY --from=builder /clash-config/ /root/.config/clash/
COPY --from=builder /clash /clash
RUN chmod +x /clash
ENTRYPOINT [ "/clash" ]