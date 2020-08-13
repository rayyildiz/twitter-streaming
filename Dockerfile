FROM golang as builder1
WORKDIR /src

ADD api .
ADD conf .
ADD twitter .
ADD util .
ADD main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o tws .


FROM node:lts as builder2
WORKDIR /src
ADD frontend/ .
RUN npm install
RUN npm run build



FROM rayyildiz/base

COPY --from=builder1 /src/tws /app/tws
COPY --from=builder2 /src/public /app/public
ADD config.json      /app/
RUN chmod 755 -Rf    /app

EXPOSE 3000

CMD ["/app/tws","-config=/app/config.json"]



