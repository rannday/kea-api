# Kea Debian Installation
https://cloudsmith.io/~isc/repos/stork/setup/#formats-deb

Install some prerequisites:  
```bash
sudo apt install -y curl gnupg apt-transport-https
```
## MySQL Installation
Download the MySQL APT Config Package
```bash
wget https://dev.mysql.com/get/mysql-apt-config_0.8.29-1_all.deb
sudo dpkg -i mysql-apt-config_0.8.29-1_all.deb
```
Just select the latest version of MySQL. Nothing else.  
Go here for the latest: https://dev.mysql.com/downloads/repo/apt/  

Then install:  
```bash
sudo apt update && sudo apt install -y mysql-server
```
MySQL will have you setup the root user password, here.  

Enable and run:  
```bash
sudo systemctl enable --now mysql
```
### Verify MySQL Installation
```bash
sudo mysql --version
mysql  Ver 8.0.42 for Linux on x86_64 (MySQL Community Server - GPL)
```
```bash
sudo ss -tulpn|grep 3306
tcp   LISTEN 0      151                *:3306             *:*    users:(("mysqld",pid=4496,fd=23))
tcp   LISTEN 0      70                 *:33060            *:*    users:(("mysqld",pid=4496,fd=21))
```
33060 is the MySQL Shell API.
## Kea Installation
Use the cloudsmith repo. ISC maintains it. Their binaries have SQL support packaged in, already.  
https://kb.isc.org/docs/isc-kea-packages  
https://www.isc.org/blogs/cloudsmith-repos    
```bash
curl -1sLf 'https://dl.cloudsmith.io/public/isc/kea-2-6/setup.deb.sh' | sudo -E bash
```
Then install:  
```bash
sudo apt update && sudo apt install -y isc-kea
```
### Verify Kea Installation
```bash
sudo kea-dhcp4 -V
2.6.3 (isc20250522135511 deb)
premium: yes (isc20250522135511 deb)
linked with:
- log4cplus 2.0.8
- OpenSSL 3.0.16 11 Feb 2025
backends:
- MySQL backend 22.2, library 3.3.14
- PostgreSQL backend 22.2, library 150013
- Memfile backend 3.0
```
```bash 
sudo kea-dhcp6 -V
2.6.3 (isc20250522135511 deb)
premium: yes (isc20250522135511 deb)
linked with:
- log4cplus 2.0.8
- OpenSSL 3.0.16 11 Feb 2025
backends:
- MySQL backend 22.2, library 3.3.14
- PostgreSQL backend 22.2, library 150013
- Memfile backend 5.0
```
## Verify Services
Kea gets enabled and started after installation. MySQL we started and enabled.   
```bash
systemctl list-units --type=service --state=running

  isc-kea-dhcp-ddns-server.service loaded active running Kea DHCP-DDNS Service
  isc-kea-dhcp4-server.service     loaded active running Kea DHCPv4 Service
  isc-kea-dhcp6-server.service     loaded active running Kea DHCPv6 Service
  mysql.service                    loaded active running MySQL Community Server
```
Stop the DDNS Service:  
```bash
sudo systemctl disable --now isc-kea-dhcp-ddns-server.service
```
Check the other two:  
```bash
sudo systemctl status isc-kea-dhcp4-server.service isc-kea-dhcp6-server.service
```
Check on the control agent:  
```bash
sudo systemctl status isc-kea-ctrl-agent.service
```
It probably failed. You need to make a password file for API users:  
```bash
echo "kea" | sudo tee /etc/kea/kea-api-password > /dev/null
```
Comment out the control sockets you're not using - most likely just DDNS.
```bash
sudo vim /etc/kea/kea-ctrl-agent.conf
```
Restart the service:
```bash
sudo systemctl restart isc-kea-ctrl-agent.service
```
Verify:  
```bash
sudo systemctl status isc-kea-ctrl-agent.service
● isc-kea-ctrl-agent.service - Kea Control Agent
	...
```
```bash
sudo ss -tulp|grep 8000
tcp   LISTEN 0      4096       127.0.0.1:8000        0.0.0.0:*    users:(("kea-ctrl-agent",pid=5740,fd=7))
```
## Updating
### Kea
### MySQL
```bash
sudo apt update && sudo apt upgrade -y
```