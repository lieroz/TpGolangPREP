FROM golang:onbuild

RUN apt-get update
RUN apt-get install -y git
RUN mkdir -p /var/www/html &&\
 git clone https://github.com/init/http-test-suite.git &&\
 mv ./http-test-suite/httptest /var/www/html

EXPOSE 80