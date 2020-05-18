FROM golang:1.13.2 as golang

WORKDIR /go/src/github.com/IBM/portieris
RUN mkdir -p /go/src/github.com/IBM/portieris
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags containers_image_openpgp -o ./bin/trust ./cmd/trust

FROM scratch
COPY --from=golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=golang /go/src/github.com/IBM/portieris/bin/trust /trust
# Create /tmp for logs and /run for working directory
RUN [ "/trust", "--mkdir",  "/tmp,/run" ]
WORKDIR /run
CMD ["/trust","--alsologtostderr","-v=4","2>&1"]
