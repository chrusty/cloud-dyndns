FROM golang:1.8.0-alpine
ADD cloud-dyndns.linux-amd64 /cloud-dyndns.linux-amd64

ENV DNS_FREQUENCY=60m
ENV DNS_ZONEID=XXXXXXXXXXXXX
ENV DNS_HOSTNAME=host.domain.com.
ENV DNS_TTL=900

CMD ["/cloud-dyndns", "-zoneid=$DNS_ZONEID", "-frequency=$DNS_FREQUENCY", "-hostname=$DNS_HOSTNAME", "-ttl=$DNS_TTL"]
