FROM golang:1.10 as golang

WORKDIR /go/src/github.com/IBM/portieris
RUN mkdir -p /go/src/github.com/IBM/portieris
COPY . ./
# RUN make build-deps
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/trust ./cmd/trust

FROM scratch
COPY --from=golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Create /tmp for log files
WORKDIR /tmp 
WORKDIR /
COPY --from=golang /go/src/github.com/IBM/portieris/bin/trust .
CMD ["./trust","--alsologtostderr","-v=4","2>&1"]