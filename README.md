## SRS-SIP

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

### Run

```
srs-sip -media-host 127.0.0.1 -media-api-port 1985 -sip-port 5060 
```