# build stage
FROM ubi8/go-toolset AS builder
RUN mkdir -p /opt/app-root/src/ACPodNaming
WORKDIR /opt/app-root/src/ACPodNaming
ENV GOPATH=/opt/app-root/
#ENV GOFLAGS="-mod=vendor"
ENV PATH="${PATH}:/opt/app-root/src/go/bin/"
COPY  . .
RUN set -x && \
    go get -u github.com/golang/dep/cmd/dep && \
    dep init && \
    dep ensure -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o acpodnaming 
# copy to running image 
FROM ubi8/ubi-minimal
WORKDIR /opt/app-root/
USER 1001
COPY --from=builder  /opt/app-root/src/ACPodNaming/acpodnaming .
WORKDIR /opt/app-root/acpodnaming
ENTRYPOINT ["/opt/app-root/acpodnaming"]
