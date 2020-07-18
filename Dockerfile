# build stage
FROM ubi8/go-toolset AS builder
WORKDIR /opt/app-root/
ENV GOPATH=/opt/app-root/acpodnaming
RUN mkdir acpodnaming
COPY  . acpodnaming/
USER 1001
WORKDIR /opt/app-root/acpodnaming
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o acpodnaming 
# copy to running image 
FROM ubi8/ubi-minimal
WORKDIR /opt/app-root/
USER 1001
COPY --from=builder  /opt/app-root/acpodnaming/acpodnaming .
WORKDIR /opt/app-root/acpodnaming
ENTRYPOINT ["/acpodnaming"]
