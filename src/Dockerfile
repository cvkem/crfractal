# A named stage: "build"
FROM golang:1.15 AS build

WORKDIR /go/src/app/

COPY fractal /go/src/app/fractal

COPY http /go/src/app/http

ENV CGO_ENABLED=0

RUN cd http && go build -o /go/bin/app


# RUN go build -o /go/bin/app/
# A new stage: "run"
#FROM gcr.io/distroless/base-debian10:nonroot AS run
FROM scratch

# Copy the binary from stage build
COPY --from=build /go/bin/app /

ENTRYPOINT ["/app"]
