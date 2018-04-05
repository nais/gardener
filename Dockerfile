FROM alpine:3.5
MAINTAINER Audun Strand <audun.fauchald.strand@nav.no>

WORKDIR /app

COPY gardener .

CMD /app/gardener  --logtostderr=true
