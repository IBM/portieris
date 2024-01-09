# This first stage of the build uses go-toolset to build the portieris binary creates
# a simplified operating system image that satisfies vulnerability scanning requirements
FROM --platform=$BUILDPLATFORM registry.access.redhat.com/ubi8/go-toolset:1.20.10-3 as builder
ARG PORTIERIS_VERSION=undefined

# switch to root user as we need to run yum and rpm to ensure packages are up to date
USER root
RUN yum update -y

# Work within the /opt/app-root/src working directory of the UBI go-toolset image
WORKDIR /opt/app-root/src/github.com/IBM/portieris
RUN mkdir -p /opt/app-root/src/github.com/IBM/portieris
# Create directory to store the built binary
RUN mkdir -p /opt/app-root/bin
COPY . ./

RUN go mod download

ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-X github.com/IBM/portieris/internal/info.Version=$PORTIERIS_VERSION" -a \
    -tags containers_image_openpgp -o /opt/app-root/bin/portieris ./cmd/portieris
RUN go version -m -v /opt/app-root/bin/portieris | (grep dep || true) | awk '{print "{\"Path\": \""$2 "\", \"Version\": \"" $3 "\"}"}' > /deps.jsonl

FROM registry.access.redhat.com/ubi8/go-toolset:1.20.10-3 as installer
ARG TARGETOS TARGETARCH
USER root
RUN yum update -y
# prep target rootfs for scratch container
WORKDIR /
RUN mkdir /image && \
    ln -s usr/bin /image/bin && \
	ln -s usr/sbin /image/sbin && \
	ln -s usr/lib64 /image/lib64 && \
	ln -s usr/lib /image/lib && \
	mkdir -p /image/{usr/bin,usr/lib64,usr/lib,root,home,proc,etc,sys,var,dev}
# see files-{amd64,s390x}.txt for a list of needed files from the UBI image to copy into our
# final "FROM scratch" image; this would need to be modified if any additional
# content was required from UBI for the Portieris binary to function.
COPY files-${TARGETARCH}.txt /tmp
RUN tar cf /tmp/files.tar -T /tmp/files-${TARGETARCH}.txt && tar xf /tmp/files.tar -C /image/ \
  && strip --strip-unneeded /image/usr/lib64/*[0-9].so
RUN rpm --root /image --initdb \
  && PACKAGES=$(rpm -qf $(cat /tmp/files-${TARGETARCH}.txt) | grep -v "is not owned by any package" | sort -u) \
  && echo dnf install -y 'dnf-command(download)' \
  && dnf download --destdir / ${PACKAGES} \
  && rpm --root /image -ivh --justdb --nodeps `for i in ${PACKAGES}; do echo $i.rpm; done`


# Check dependencies for vulnerabilities
FROM --platform=$BUILDPLATFORM sonatypecommunity/nancy:alpine as nancy
COPY --from=builder /deps.jsonl /
COPY /.nancy-ignore /
RUN cat /deps.jsonl | nancy --skip-update-check --loud sleuth --no-color
RUN echo true> /nancy-checked

#################################################################################
# Finally, copy the minimal image contents and the built binary into the scratch image
FROM scratch
COPY --from=installer /image/ /
COPY --from=builder /opt/app-root/bin/portieris /portieris
# buildkit skips stages which dont contribute to the final image
COPY --from=nancy /nancy-checked /nancy-checked 
# Create /tmp for logs and /run for working directory
RUN [ "/portieris", "--mkdir",  "/tmp,/run" ]
WORKDIR /run
# quiet image config checkers, this is the default runAsUser in the deployment
USER 1000060001
CMD ["/portieris","--alsologtostderr","-v=4","2>&1"]
