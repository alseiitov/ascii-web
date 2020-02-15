FROM golang:1.13.8-alpine3.11
LABEL name="ASCII ART Generator"
LABEL description="alem school project"
LABEL authors="alseiitov; satsuls; bortico"
LABEL release-date="15.02.2020"
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]