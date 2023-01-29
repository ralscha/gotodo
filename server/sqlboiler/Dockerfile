FROM golang:1.19.5
RUN go install github.com/volatiletech/sqlboiler/v4@latest
RUN go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
WORKDIR /src/sqlboiler
CMD ["sqlboiler", "mysql"]
