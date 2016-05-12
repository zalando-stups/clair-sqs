# use clair native image
FROM quay.io/coreos/clair:v1.2.0

# add skipper
ADD vendor/github.com/zalando /go/src/github.com/zalando
RUN go get -v github.com/zalando/skipper/cmd/skipper
RUN go install -v github.com/zalando/skipper/cmd/skipper
EXPOSE 8080

# add configurations
ADD clair.conf /etc/clair/config.yaml
ADD skipper.eskip /etc/skipper.eskip

# add supervisor for our multiprocess container
RUN apt-get update && \
    apt-get install -y supervisor && \
    apt-get autoremove -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
ADD supervisord.conf /etc/supervisord.conf
ADD run.sh /run.sh

ENTRYPOINT ["/run.sh"]
