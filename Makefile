.PHONY: help default clean gbserver

default: gbserver

clean:
	rm -f ./objs/gbserver

gbserver: ./objs/gbserver

./objs/gbserver:
	go build -o objs/gbserver .

help:
	@echo "Usage: make [gbserver]"
	@echo "     gbserver       Make the gbserver to ./objs/gbserver"