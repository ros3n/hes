FROM ubuntu:16.04
RUN apt-get update -qq

RUN apt-get install -y ca-certificates

EXPOSE 5556
EXPOSE 8080

COPY hes-mailer /
CMD ["/hes-mailer"]