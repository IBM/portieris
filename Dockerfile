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

# Use another intermediary step to identify and extract the minimum content required for the runtime image.
# The purpose of this is to keep the image size and attack surface as small as possible,
#  while providing enough information for vulnerability scanning tools to inspect it.
FROM registry.access.redhat.com/ubi8/s2i-base:latest as installer
RUN yum upgrade -y
# prep target rootfs for scratch container
WORKDIR /
RUN mkdir /image && \
    ln -s usr/bin /image/bin && \
	ln -s usr/sbin /image/sbin && \
	ln -s usr/lib64 /image/lib64 && \
	ln -s usr/lib /image/lib && \
	mkdir -p /image/{usr/bin,usr/lib64,usr/lib,root,home,proc,etc,sys,var,dev}
# see files.txt for a list of needed files from the UBI image to copy into our
# final "FROM scratch" image; this would need to be modified if any additional
# content was required from UBI for a specific binary.
COPY files.txt /tmp
RUN tar cf /tmp/files.tar -T /tmp/files.txt && tar xf /tmp/files.tar -C /image/ \
  && strip --strip-unneeded /image/usr/lib64/*[0-9].so
RUN rpm --root /image --initdb \
  && PACKAGES=$(rpm -qf $(cat /tmp/files.txt) | grep -v "is not owned by any package" | sort -u) \
  && echo dnf install -y 'dnf-command(download)' \
  && dnf download --destdir / ${PACKAGES} \
  && rpm --root /image -ivh --justdb --nodeps `for i in ${PACKAGES}; do echo $i.rpm; done`

# Finally copy the minimal image contents and the built binary into the scratch image
FROM scratch

COPY --from=installer /image/ /

COPY --from=gobuild /opt/app-root/bin/portieris /portieris

# Create /tmp for logs and /run for working directory
RUN [ "/portieris", "--mkdir",  "/tmp,/run" ]
WORKDIR /run
CMD ["/portieris","--alsologtostderr","-v=4","2>&1"]
