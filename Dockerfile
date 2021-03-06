FROM alpine:3.8
MAINTAINER Audun Strand <audun.fauchald.strand@nav.no>

RUN apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /app

COPY gardener .

CMD /app/gardener  --logtostderr=true --clustername=$clustername --slackUrl=$slackUrl
