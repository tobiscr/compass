FROM golang:1.14.4-alpine3.12 as builder

ENV BASE_APP_DIR /go/src/github.com/kyma-incubator/compass/components/admiral-watcher
WORKDIR ${BASE_APP_DIR}

COPY . .

RUN go build -v -o /app/admiral-watcher ./cmd/main.go

FROM alpine:edge
LABEL source = git@github.com:kyma-incubator/compass.git
WORKDIR /app

RUN apk --no-cache add curl ca-certificates git

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.16.3/bin/linux/amd64/kubectl
RUN chmod u+x kubectl && mv kubectl /bin/kubectl

COPY --from=builder /app /app
COPY --from=builder /go/src/github.com/kyma-incubator/compass/components/admiral-watcher/resources /app/

CMD ["/app/admiral-watcher"]
