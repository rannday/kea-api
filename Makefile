.PHONY: test

test:
	go test github.com/rannday/isc-kea/client \
					github.com/rannday/isc-kea/agent \
	        github.com/rannday/isc-kea/dhcp4 \
	        github.com/rannday/isc-kea/dhcp6
