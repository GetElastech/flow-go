FROM golang:1.17

COPY flow-localnet.json /go
COPY flow-cli /go/flow-cli

WORKDIR /go/flow-cli
RUN go mod download
RUN make
RUN /go/flow-cli/cmd/flow/flow version

CMD /go/flow-cli/cmd/flow/flow -f /go/flow-localnet.json -n observer blocks get latest
