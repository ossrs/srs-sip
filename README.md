## SRS-SIP

### Sequence
```mermaid
sequenceDiagram
    Device ->> SRS-SIP : Register(GB28181)
    SRS-SIP ->>  SRS : Publish Request(gb/v1/publish)
    SRS ->> SRS-SIP : Response(with port)
    SRS-SIP ->>  Device : Auto Invite(GB28181)
    Device -->> SRS : Media Stream
    SRS -->> Client : Media Stream
```

### Building from source

Pre-requisites:
- Go 1.17+ is installed
- GOPATH/bin is in your PATH

Then run
```
git clone https://github.com/ossrs/srs-sip
cd srs-sip
./bootstrap.sh
mage
```

If you are on a Unix-like system, you can also run the following command.
```
make
```

### Run

```
srs-sip -sip-port 5060 -media-addr 127.0.0.1:1985
```

- `sip-port` : the SIP port, this program listen on, for device register with gb28181
- `media-addr` : the API address for SRS, typically on port 1985, used to send HTTP requests to "gb/v1/publish"
