FROM golang:1.22 as build

WORKDIR /home/sda/app/

COPY  . .

#RUN ls
#RUN echo $(pwd)
#RUN cd src/tender
#RUN ls src/tender
COPY src/tender/go.mod ./
COPY src/tender/go.sum ./
COPY src/tender ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd/tender/main.go

FROM alpine

RUN mkdir /app

COPY --from=build /home/sda/app/myapp /app/

RUN chmod +x ./app/myapp

EXPOSE 8080
CMD ["/app/myapp"]

