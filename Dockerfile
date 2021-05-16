FROM golang:alpine AS builder
RUN apk --no-cache add ca-certificates
#RUN apk add --no-cache ca-certificates git

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

# Import the code from the context.
COPY ./ ./

# Build the executable to `/app`. Mark the build as statically linked.
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app/ .
ADD config.json /app/

# Final stage: the running container.
FROM scratch  AS final

# Import the compiled executable from the first stage.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app

# Expose both 443 and 80 to our application
#EXPOSE 443
EXPOSE 80

WORKDIR /app
# Run the compiled binary.
ENTRYPOINT ["./WordsBot"]