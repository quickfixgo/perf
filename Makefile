all: inbound outbound

clean:
	rm inbound/inbound
	rm outbound/outbound

inbound: 
	cd $@; go build .

outbound: 
	cd $@; go build .


.PHONY: inbound outbound
