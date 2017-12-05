FROM golang:alpine

ARG TOKEN
ENV TOKEN ${TOKEN}

ARG URL
ENV TOKEN ${URL}

ARG PORT
ENV PORT ${PORT}

COPY . /go/src/taskbot

RUN go install taskbot

EXPOSE 80

CMD [ "taskbot" ]