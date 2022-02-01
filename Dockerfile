FROM golang:1.16-buster@sha256:3cb83def925efd866e72dbc724e1f7f0175d7d5559fe45a41585fc33386a4b6e AS build
RUN groupadd auth0 && useradd -m developer -g auth0
USER developer
RUN mkdir /home/developer/app
WORKDIR /home/developer/app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY cmd cmd
COPY pkg pkg
RUN go build -o api-server ./cmd/api-server/main.go

FROM gcr.io/distroless/base-debian10@sha256:586e10ceb097684dcd3e455dbb6d4141f3dd28986719de487d76d4c7c9da1a35
COPY --from=build /home/developer/app/api-server /api-server
ENV APP_ENV=production
USER 1000
EXPOSE 6060
CMD ["/api-server"]
