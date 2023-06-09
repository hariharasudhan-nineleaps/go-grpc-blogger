##########
# BUILD ##
##########

FROM golang:1.20.5-alpine3.17 as builder

RUN apk update && apk add --no-cache git

WORKDIR /home/app/auth

COPY . .

RUN ls -la

RUN go mod tidy

RUN go build -v -o .

#########
# FINAL #
#########

FROM golang:1.20.5-alpine3.17

ARG DEFAULT_PORT=80

RUN apk update

WORKDIR /home/app/auth

# create the app user
RUN addgroup -S app && adduser -S app -G app

RUN ls -lh

COPY --chown=app:app --from=builder /home/app/auth/auth .
COPY --chown=app:app --from=builder /home/app/auth/.env .

RUN ls -la

RUN echo "hello"

USER app

ENV PORT $DEFAULT_PORT

EXPOSE $PORT

RUN ls -lh

CMD ["./auth"]






