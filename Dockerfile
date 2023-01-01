# Build stage
FROM golang:1.17 as build

WORKDIR /go/src/app
COPY . .

ENV CGO_ENABLED=0
ENV PORT=8080
ENV DB_PORT=3306
ENV DB_NAME=sql6587040
ENV DB_PASSWORD=YzKa1MWg8q
ENV DB_USER=sql6587040
ENV DB_HOST=sql6.freesqldatabase.com
ENV DB_DRIVER=mysql
ENV API_SECRET=chatgunawan
ENV TOKEN_HOUR_LIFESPAN=1
ENV CLOUDINARY_URL=cloudinary://672374291877475:Kfvgb4DjyW65yA7mHOYurxfVuyw@da9irqgak

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -v -o go-app

# Run stage
FROM alpine:3.11
COPY --from=build go/src/app/ app/
CMD ["./app/go-app"]