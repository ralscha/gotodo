FROM golang:1.17.5
RUN go install github.com/volatiletech/sqlboiler/v4@latest
RUN go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
WORKDIR /src/
CMD ["make", "db/run-sqlboiler"]
