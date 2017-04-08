FROM rayyildiz/base

MAINTAINER "Ramazan AYYILDIZ <rayyildiz@gmail.com>"

ADD twitterStreaming /app/tws
ADD public           /app/public
ADD config.json      /app/
RUN chmod 755 -Rf    /app

EXPOSE 3000

CMD ["/app/tws","-config=/app/config.json"]

