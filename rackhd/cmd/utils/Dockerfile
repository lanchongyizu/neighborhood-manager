FROM rackhd/golang:1.7.0-wheezy

ADD ./bin/endpoint /go/bin/endpoint

CMD ["/go/bin/endpoint","-endpoint-address=http://localhost:10002","-backend-address=consulclient:8500"]
