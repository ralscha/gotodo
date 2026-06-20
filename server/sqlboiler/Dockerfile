FROM golang:1.25.5
RUN go install github.com/aarondl/sqlboiler/v4@latest
RUN go install github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-psql@latest
WORKDIR /src/sqlboiler
CMD ["sqlboiler", "psql"]
