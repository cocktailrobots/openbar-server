#

## OS Installation
Use the [Raspberry Pi Imager](https://www.raspberrypi.com/software/) to install the Raspberry Pi Os Lite (64-bit). 

First click the "Choose OS" button and select "Raspberry Pi OS (other)" and then "Raspberry Pi OS Lite (64-bit)".

Then click the "Choose SD Card" button and select the SD card you want to install the OS on.

Once the OS and SD card are selected, click the settings button in order to configure the OS installation. Enable SSH
configuring an ssh key or username/password authentication. Also check the "Configure Wireless LAN" box and enter the
SSID and password for your wireless network.

Finally, click the "Write" button to start the installation. Once the installation is complete, eject the SD card and
insert it into the Raspberry Pi and power it on. Use the SSH key or username/password to log into the Raspberry Pi.

## Dolt Installation and Setup


### Install the Dolt binary

`sudo bash -c 'curl -L https://github.com/dolthub/dolt/releases/latest/download/install.sh | sudo bash'`

### Clone your fork of the Cocktails Repo

```
openbar@openbar:~ $ sudo bash
root@openbar:/var# cd dbs
root@openbar:/var/dbs# dolt clone openbar/cocktails
cloning https://doltremoteapi.dolthub.com/openbar/cocktails
root@openbar:/var/dbs# echo 'log_level: "debug"
user:
  name: "openbar"
listener:
  host: "0.0.0.0"
  port: 3306
  max_connections: 5' > config.yaml
root@openbar:/var/dbs# mkdir /var/log/dolt
root@openbar:/var/dbs# echo '[Unit]
Description=dolt service
After=network.target

[Service]
Type=simple
User=root
Group=root
AmbientCapabilities=CAP_NET_BIND_SERVICE

Environment=DOLT_ROOT_PATH=/var/dbs/
WorkingDirectory=/var/dbs/
ExecStart=dolt --out-and-err "/var/log/dolt/sqlserver.log" --ignore-lock-file sql-server --config config.yaml

LimitNOFILE=100000

Restart=always
RestartSec=1

MemoryAccounting=true
MemoryMax=90%

[Install]
WantedBy=multi-user.target' > /etc/systemd/system/dolt.service
root@openbar:/var/dbs# exit
openbar@openbar:~ $ sudo systemctl daemon-reload
openbar@openbar:~ $ sudo systemctl enable dolt.service
```

`sudo apt install mariadb-client`


```bash
openbar@openbar:~ $ mysql -h127.0.0.1 -uopenbar
Welcome to the MariaDB monitor.  Commands end with ; or \g.
Your MySQL connection id is 1
Server version: 5.7.9-Vitess Dolt

Copyright (c) 2000, 2018, Oracle, MariaDB Corporation Ab and others.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

MySQL [(none)]> show databases;
+--------------------+
| Database           |
+--------------------+
| cocktails          |
| information_schema |
| mysql              |
+--------------------+
3 rows in set (0.010 sec)

MySQL [(none)]> exit
Bye
```

## Install openbar-server

### Installing openbar-server from Source

`wget -o go1.21.1.linux-arm64.tar.gz https://go.dev/dl/go1.21.1.linux-arm64.tar.gz`
`rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.1.linux-arm64.tar.gz`
`vi ~/.bashrc`
`PATH=$PATH:/usr/local/go/bin:~/go/bin`

`sudo apt install git`
`sudo apt install libx11-dev`
`sudo apt install -y xvfb`
`sudo apt install libncurses5-dev`

`git clone https://github.com/cocktailrobots/openbar-server.git`

`cd openbar-server/cmd/openbar-server/`
`go install .`

```bash
$ openbar-server
2023/09/16 11:03:46 Usage: openbar-server [-migration-dir=<migration_file_dir>] <config file>
```

### Installing openbar-server from a Release

## Configuring openbar-server



```bash
sudo bash
mkdir /etc/openbar-server
cd /etc/openbar-server
vim config.yaml
```



## Enabling Android USB Tethering

Add to `/etc/dhcpcd.conf`

```
interface usb0
static ip_address=192.168.42.42/24
static routers=192.168.42.129
static domain_name_servers=192.168.42.129
```

Restart dhcpcd

`sudo systemctl restart dhcpcd`

Open Android Settings -> Network & Internet -> Hotspot & tethering and enable USB tethering.


