FROM golang:alpine as builder
#RUN apk add --no-cache curl jq
RUN mkdir -p /assets
WORKDIR /assets
COPY . /go/src/github.com/lorands/maven-stage-resource
ENV CGO_ENABLED 0
RUN go build -o /assets/in github.com/lorands/maven-stage-resource/in/cmd/in
RUN go build -o /assets/out github.com/lorands/maven-stage-resource/out/cmd/out
RUN go build -o /assets/check github.com/lorands/maven-stage-resource/check/cmd/check
WORKDIR /go/src/github.com/lorands/maven-stage-resource
RUN set -e; for pkg in $(go list ./... | grep -v "acceptance"); do \
		go test -o "/tests/$(basename $pkg).test" -c $pkg; \
	done

FROM alpine:edge AS resource
RUN apk add --no-cache bash tzdata ca-certificates
COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*
# RUN mv /opt/resource/cf /usr/bin/cf
# RUN mv /opt/resource/autopilot /usr/bin/autopilot
# RUN /usr/bin/cf install-plugin -f /usr/bin/autopilot

# FROM resource AS tests
# COPY --from=builder /tests /go-tests
# COPY out/assets /go-tests/assets
# WORKDIR /go-tests
# RUN set -e; for test in /go-tests/*.test; do \
# 		$test; \
# 	done

FROM resource
