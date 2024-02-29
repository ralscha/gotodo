FROM golang:1.21.7
RUN go install github.com/volatiletech/sqlboiler/v4@latest
RUN go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
WORKDIR /src/sqlboiler
CMD ["sqlboiler", "mysql"]
