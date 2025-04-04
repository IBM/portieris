# syntax=docker/dockerfile:1

# This first stage of the build uses go-toolset from RHEL.
# We specify a higher version # than available, meaning that the  `toolchain` statement in `go.mod`
# will replace go-toolset with the specified version of go. The specified version is then used
# to build the portieris binary. go-toolset is used to create a simplified operating system image
# that satisfies vulnerability scanning requirements
ARG BASE_IMAGE=registry.access.redhat.com/ubi9/go-toolset:1.22.9
FROM $BASE_IMAGE AS builder
ARG PORTIERIS_VERSION=undefined
ARG TARGETOS TARGETARCH

# prep target rootfs for scratch container
USER root
WORKDIR /
RUN mkdir /image \
 && ln -s usr/bin /image/bin \
 && ln -s usr/sbin /image/sbin \
 && ln -s usr/lib64 /image/lib64 \
 && ln -s usr/lib /image/lib \
 && mkdir -p /image/{usr/bin,usr/lib64,usr/lib,root,home,proc,etc,sys,var,dev}
# see files-{amd64,s390x}.txt for a list of needed files from the UBI image to copy into our
# final "FROM scratch" image; this would need to be modified if any additional
# content was required from UBI for the Portieris binary to function.
COPY files-${TARGETARCH}.txt /tmp
RUN tar cf /tmp/files.tar -T /tmp/files-${TARGETARCH}.txt && tar xf /tmp/files.tar -C /image/ \
 && rpm --root /image --initdb \
 && PACKAGES=$(rpm -qf $(cat /tmp/files-${TARGETARCH}.txt) --queryformat "%{NAME}\n" | grep -v "is not owned by any package" | sort -u) \
 && dnf download --destdir /rpmcache ${PACKAGES} \
 && rpm --root /image -ivh --justdb --nodeps /rpmcache/*.rpm

# setup workdir and build binary
WORKDIR /go/github.com/IBM/portieris
COPY . ./
# override GOTOOLCHAIN because go-toolset has a value of local found in /usr/lib/golang/go.env
# this is because we want Go version upgrades to be applied based on the toolchain value in go.mod
# a mounted secret of GOPROXY allows setting of the GOPROXY env var
ENV GOTOOLCHAIN=auto
RUN --mount=type=secret,id=GOPROXY,env=GOPROXY go mod download \
    && CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-X github.com/IBM/portieris/internal/info.Version=$PORTIERIS_VERSION" -a \
    -tags containers_image_openpgp -o portieris ./cmd/portieris \
    && go version -m -v portieris | (grep dep || true) | awk '{print "{\"Path\": \""$2 "\", \"Version\": \"" $3 "\"}"}' > /deps.jsonl

# Check dependencies for vulnerabilities
FROM sonatypecommunity/nancy:alpine AS nancy
COPY --from=builder /deps.jsonl /
COPY /.nancy-ignore /
RUN cat /deps.jsonl | nancy --skip-update-check --loud sleuth --no-color
RUN echo true> /nancy-checked

#################################################################################
# Finally, copy the minimal image contents and the built binary into the scratch image
FROM scratch
COPY --from=builder /image/ /
COPY --from=builder /go/github.com/IBM/portieris/portieris /portieris
# buildkit skips stages which dont contribute to the final image
COPY --from=nancy /nancy-checked /nancy-checked
# Create /tmp for logs and /run for working directory
RUN [ "/portieris", "--mkdir",  "/tmp,/run" ]
WORKDIR /run
# quiet image config checkers, this is the default runAsUser in the deployment
USER 1000060001
CMD ["/portieris","--alsologtostderr","-v=4","2>&1"]
