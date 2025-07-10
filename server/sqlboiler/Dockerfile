FROM golang:1.22.6
RUN go install github.com/aarondl/sqlboiler/v4@latest
RUN go install github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-mysql@latest
WORKDIR /src/sqlboiler
CMD ["sqlboiler", "mysql"]
