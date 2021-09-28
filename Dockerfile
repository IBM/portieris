FROM registry.access.redhat.com/ubi8/go-toolset:1.15.14-10 as gobuild
ARG VERSION=undefined
# Work within the /opt/app-root/src working directory of the UBI go-toolset image
WORKDIR /opt/app-root/src/github.com/IBM/portieris
RUN mkdir -p /opt/app-root/src/github.com/IBM/portieris
# Create directory to store the built binary
RUN mkdir -p /opt/app-root/bin
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-X github.com/IBM/portieris/internal/info.Version=$VERSION" -a \
    -tags containers_image_openpgp -o /opt/app-root/bin/portieris ./cmd/portieris

FROM registry.access.redhat.com/ubi8/ubi-minimal
COPY --from=gobuild /opt/app-root/bin/portieris /portieris
# Create /tmp for logs and /run for working directory
RUN [ "/portieris", "--mkdir",  "/tmp,/run" ]
WORKDIR /run
CMD ["/portieris","--alsologtostderr","-v=4","2>&1"]
