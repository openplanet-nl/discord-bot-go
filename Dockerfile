FROM golang:alpine AS build
COPY . /src
WORKDIR /src
RUN go build

FROM alpine:edge
COPY --from=build /src/openplanet-bot /bot

WORKDIR /
ENTRYPOINT [ "/bot" ]
