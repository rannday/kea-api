# ISC Kea Setup
## Installation
## Configuration

# Configuring
Add your username to these groups:
```bash
sudo usermod -aG adm,mysql,_kea,systemd-journal $(whoami)
exit
```
Log back in
## Firewall
Example nftables firewall:
```bash
/etc/nftables.conf
#!/usr/sbin/nft -f

flush ruleset

table inet filter {
  set mgmt_allow {
    type ipv4_addr; flags interval;
    elements = {
      192.168.10.0/24,  # Management range
    }
  }
  set api_allow {
    type ipv4_addr; flags interval;
    elements = {
      192.168.10.0/24,
      172.17.1.10       # Other HA member
    }
  }

  chain input {
    type filter hook input priority 0;
    policy drop;

    iif lo accept                                           # Allow localhost
    ct state established,related accept                     # Early tcp est,rel

    ip protocol icmp counter accept                         # IPv4 ICMP
    ip6 nexthdr icmpv6 counter accept                       # IPv6 ICMP

    ip protocol udp udp dport { 67, 68 } counter accept     # IPv4 DHCP
    ip6 nexthdr udp udp dport { 546, 547 } counter accept   # IPv6 DHCP
    tcp dport 80 counter accept                             # Certbot
    tcp dport 443 ip saddr @api_allow counter accept        # Caddy
    tcp dport 3306 ip saddr @api_allow counter accept       # MySQL
    tcp dport 8000 ip saddr @api_allow counter accept       # Kea Control Agent API
    tcp dport 22 ip saddr @mgmt_allow counter accept        # SSH

    log prefix "nftables INPUT drop: " flags all counter drop
  }

  chain forward {
    type filter hook forward priority 0;
    policy drop;
  }

  chain output {
    type filter hook output priority 0;
    policy accept;
  }
}
```
## MySQL
Make backups of the original configuration files:  
```bash
sudo find /etc/mysql -type f -name '*.cnf' -exec bash -c 'for f; do rel="${f#/etc/mysql/}"; dest="/root/backups/mysql/$rel"; sudo mkdir -p "$(dirname "$dest")"; sudo cp "$f" "${dest}.orig"; done' _ {} +
```
Backups will be in `/root/backups/mysql`  

Run this to add `log_bin_trust_function_creators = 1` to the bottom of `/etc/mysql/my.cnf`:  
```bash
echo -e '\n[mysqld]\nlog_bin_trust_function_creators = 1' | sudo tee -a /etc/mysql/my.cnf > /dev/null
```
You can verify this gets enabled by running:
```bash
mysql -u root -p -e "SHOW VARIABLES LIKE 'log_bin_trust_function_creators';"
```
Run this to add `innodb_flush_log_at_trx_commit=2` to the bottom of `/etc/mysql/mysql.conf.d/mysqld.cnf`  
```bash
echo -e "\ninnodb_flush_log_at_trx_commit=2" | sudo tee -a /etc/mysql/mysql.conf.d/mysqld.cnf > /dev/null
```
Then restart mysql:  
```bash
sudo systemctl restart mysql.service
```
Create the kea database and user:  
```bash
mysql -u root -p <<EOF
CREATE DATABASE kea;
CREATE USER 'kea'@'localhost' IDENTIFIED BY 'kea';
GRANT ALL PRIVILEGES ON kea.* TO 'kea'@'localhost';
FLUSH PRIVILEGES;
EOF
```
Then setup the kea database:  
```bash
sudo kea-admin db-init mysql -u kea -p kea -n kea -h localhost
```
You'll most likely run into some weird issues with getting MySQL setup. Just work through them. ChatGPT makes it easy to fix anything that comes up.

Add this so you don't have to keep entering Kea's MySQL password:  
```bash
echo -e "[client]\nuser=kea\npassword=kea" > ~/.my.cnf && chmod 600 ~/.my.cnf
```
### Verify Kea Tables
This will work if you setup `~/.my.cnf`  
```bash
mysql kea -e "SHOW TABLES;"
+-------------------------------+
| Tables_in_kea                 |
+-------------------------------+
| dhcp4_audit                   |
| dhcp4_audit_revision          |
| dhcp4_client_class            |
| dhcp4_client_class_dependency |
| dhcp4_client_class_order      |
| dhcp4_client_class_server     |
| dhcp4_global_parameter        |
| dhcp4_global_parameter_server |
| dhcp4_option_def              |
| dhcp4_option_def_server       |
| dhcp4_options                 |
| dhcp4_options_server          |
| dhcp4_pool                    |
| dhcp4_server                  |
| dhcp4_shared_network          |
| dhcp4_shared_network_server   |
| dhcp4_subnet                  |
| dhcp4_subnet_server           |
| dhcp6_audit                   |
| dhcp6_audit_revision          |
| dhcp6_client_class            |
| dhcp6_client_class_dependency |
| dhcp6_client_class_order      |
| dhcp6_client_class_server     |
| dhcp6_global_parameter        |
| dhcp6_global_parameter_server |
| dhcp6_option_def              |
| dhcp6_option_def_server       |
| dhcp6_options                 |
| dhcp6_options_server          |
| dhcp6_pd_pool                 |
| dhcp6_pool                    |
| dhcp6_server                  |
| dhcp6_shared_network          |
| dhcp6_shared_network_server   |
| dhcp6_subnet                  |
| dhcp6_subnet_server           |
| dhcp_option_scope             |
| host_identifier_type          |
| hosts                         |
| ipv6_reservations             |
| lease4                        |
| lease4_pool_stat              |
| lease4_stat                   |
| lease4_stat_by_client_class   |
| lease6                        |
| lease6_pool_stat              |
| lease6_relay_id               |
| lease6_remote_id              |
| lease6_stat                   |
| lease6_stat_by_client_class   |
| lease6_types                  |
| lease_hwaddr_source           |
| lease_state                   |
| logs                          |
| modification                  |
| option_def_data_type          |
| parameter_data_type           |
| schema_version                |
+-------------------------------+
```
## Kea
Make backups of the original configuration files:  
```bash
sudo find /etc/kea -type f -name '*.conf' -exec bash -c 'for f; do rel="${f#/etc/kea/}"; dest="/root/backups/kea/$rel"; sudo mkdir -p "$(dirname "$dest")"; sudo cp "$f" "${dest}.orig"; done' _ {} +
```
Backups will be in `/root/backups/kea`  
### DHCPv4
Base configuration needed to get things kick-started for configuring via SQL.  
```bash
sudo chmod 755 /etc/kea

sudo mkdir -p /var/log/kea
sudo chown -R _kea:_kea /var/log/kea

sudo vim /etc/kea/kea-dhcp4.conf
```
```bash
{
  "Dhcp4": {
    "interfaces-config": {
      "interfaces": [ "eth0" ]
    },
    "control-socket": {
      "socket-type": "unix",
      "socket-name": "kea4-ctrl-socket"
    },
    "lease-database": {
      "type": "mysql",
      "name": "kea",
      "user": "kea",
      "password": "kea",
      "host": "127.0.0.1",
      "port": 3306
    },
    "hosts-database": {
      "type": "mysql",
      "name": "kea",
      "user": "kea",
      "password": "kea",
      "host": "127.0.0.1",
      "port": 3306
    },
    "config-control": {
      "config-databases": [
        {
          "type": "mysql",
          "name": "kea",
          "user": "kea",
          "password": "kea",
          "host": "127.0.0.1",
          "port": 3306
        }
      ]
    },
    "hooks-libraries": [
      { "library": "/usr/lib/x86_64-linux-gnu/kea/hooks/libdhcp_bootp.so" },
      { "library": "/usr/lib/x86_64-linux-gnu/kea/hooks/libdhcp_flex_option.so" },
      {
      "library": "/usr/lib/x86_64-linux-gnu/kea/hooks/libdhcp_ha.so"
      },
      {
      "library": "/usr/lib/x86_64-linux-gnu/kea/hooks/libdhcp_lease_cmds.so"
      },
      {
      "library": "/usr/lib/x86_64-linux-gnu/kea/hooks/libdhcp_mysql_cb.so"
      },
      {
      "library": "/usr/lib/x86_64-linux-gnu/kea/hooks/libdhcp_perfmon.so"
      },
      {
      "library": "/usr/lib/x86_64-linux-gnu/kea/hooks/libdhcp_run_script.so"
      },
      {
      "library": "/usr/lib/x86_64-linux-gnu/kea/hooks/libdhcp_stat_cmds.so"
      }
    ],
    "loggers": [
      {
      "name": "kea-dhcp4",
      "output-options": [
        {
        "output": "/var/log/kea/dhcp4-dhcp4.log",
        "pattern": "%d{%Y-%m-%d %H:%M:%S.%q} %c %p: %m\n"
        }
      ],
      "severity": "INFO"
      },
      {
      "name": "kea-dhcp4.leases",
      "output-options": [
        {
        "output": "/var/log/kea/dhcp4-leases.log",
        "pattern": "%d{%Y-%m-%d %H:%M:%S.%q} %c %p: %m\n"
        }
      ],
      "severity": "INFO"
      },
      {
      "name": "kea-dhcp4.packets",
      "output-options": [
        {
        "output": "/var/log/kea/dhcp4-packets.log",
        "pattern": "%d{%Y-%m-%d %H:%M:%S.%q} %c %p: %m\n"
        }
      ],
      "severity": "INFO"
      },
      {
      "name": "kea-dhcp4.dhcp4",
      "output-options": [
        {
        "output": "/var/log/kea/dhcp4-subsys.log",
        "pattern": "%d{%Y-%m-%d %H:%M:%S.%q} %c %p: %m\n"
        }
      ],
      "severity": "INFO"
      },
      {
      "name": "kea-dhcp4.options",
      "output-options": [
        {
        "output": "/var/log/kea/dhcp4-options.log",
        "pattern": "%d{%Y-%m-%d %H:%M:%S.%q} %c %p: %m\n"
        }
      ],
      "severity": "INFO"
      },
      {
      "name": "kea-dhcp4.config",
      "output-options": [
        {
        "output": "/var/log/kea/dhcp4-config.log",
        "pattern": "%d{%Y-%m-%d %H:%M:%S.%q} %c %p: %m\n"
        }
      ],
      "severity": "INFO"
      },
      {
      "name": "kea-dhcp4.commands",
      "output-options": [
        {
        "output": "/var/log/kea/dhcp4-commands.log",
        "pattern": "%d{%Y-%m-%d %H:%M:%S.%q} %c %p: %m\n"
        }
      ],
      "severity": "INFO"
      },
      {
      "name": "kea-dhcp4.hooks",
      "output-options": [
        {
        "output": "/var/log/kea/dhcp4-hooks.log",
        "pattern": "%d{%Y-%m-%d %H:%M:%S.%q} %c %p: %m\n"
        }
      ],
      "severity": "INFO"
      }
    ]
  }
}

```
