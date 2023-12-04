FROM golang AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o match-display main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
# Expose port for accessing the app
EXPOSE 8080
# This allows Heroku bind its PORT 
# Heroku needs to use its own PORT before the App can be made accessible
EXPOSE $PORT 
WORKDIR /root/
COPY --from=builder /app/match-display ./
ENTRYPOINT [ "./match-display" ]
