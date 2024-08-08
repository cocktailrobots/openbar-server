# Openbar Rasberry Pi Setup

## OS Installation
Use the [Raspberry Pi Imager](https://www.raspberrypi.com/software/) to install the Raspberry Pi Os Lite (64-bit). 

First click the "Choose OS" button and select "Raspberry Pi OS (other)" and then "Raspberry Pi OS Lite (64-bit)".

Then click the "Choose SD Card" button and select the SD card you want to install the OS on.

Once the OS and SD card are selected, click the settings button in order to configure the OS installation. Enable SSH
configuring an ssh key or username/password authentication. Also check the "Configure Wireless LAN" box and enter the
SSID and password for your wireless network.

Finally, click the "Write" button to start the installation. Once the installation is complete, eject the SD card and
insert it into the Raspberry Pi and power it on. Use the SSH key or username/password to log into the Raspberry Pi.

## Prerequisites

```
sudo bash
cd /root
apt update
apt upgrade
apt install vim netcat-traditional python3-pip libncurses5-dev mariadb-client python3-pip python3-twisted git -y
curl -o https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash
exit
```

```
sudo bash
cd /root
nvm install v18.18.0
```

## Dolt Installation and Setup


### Install the Dolt binary

`sudo bash -c 'curl -L https://github.com/dolthub/dolt/releases/latest/download/install.sh | sudo bash'`

### Clone your fork of the Cocktails Repo

```bash
openbar@openbar:~ $ sudo bash
root@openbar:/home/username# cd /var
root@openbar:/var# mkdir dbs
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
openbar@openbar:~ $ sudo systemctl start dolt
```


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

#### Install go

Remove any old go installations
`rm -rf /usr/local/go`

Download the latest arm64 go installation
`wget https://go.dev/dl/go1.21.1.linux-arm64.tar.gz`

Untar the downloaded file
`sudo tar -C /usr/local -xzf go1.21.1.linux-arm64.tar.gz`

Modify the ~/.bashrc file using vi or whatever editor you like. Add the following line to the end of the file. This will add the go binary to the PATH variable
`PATH=$PATH:/usr/local/go/bin:~/go/bin`

#### Configure git

```
  git config --global user.email "you@example.com"
  git config --global user.name "Your Name"
```

#### Clone the openbar-server repo and install
```
git clone https://github.com/cocktailrobots/openbar-server.git
cd openbar-server/cmd/openbar-server/
go install .
```

#### Run openbar-server to test installation succeeded
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
vim /etc/systemd/system/openbar-server.service
```

```
[Unit]
Description=openbar server
After=dolt.service

[Service]
Type=simple
User=root
Group=root
AmbientCapabilities=CAP_NET_BIND_SERVICE

Environment=DOLT_ROOT_PATH=/var/dbs/
WorkingDirectory=/etc/openbar-server
ExecStart=/home/openbar/go/bin/openbar-server config.yaml

LimitNOFILE=100000

Restart=always
RestartSec=5

MemoryAccounting=true
MemoryMax=90%

[Install]
WantedBy=multi-user.target
```

```yaml
hardware:
  debug:
    num-pumps: 8
    out-file: "/var/log/openbar-server/debug.log"
  gpio:
    pins: [...]
  sequent:
    expected-board-count: 1
reverse-pin:
  pin: 4
  forwand-high: true
buttons:
  gpio:
    pins: [...]
    debounce-duration: 10
    active-low: false
    pull-up: true
db:
  host: 127.0.0.1
  port: 3306
  user: openbar
  pass: password
cocktails-api:
  port: 8675
  host: 0.0.0.0
openbar-api:
  port: 3099
  host: 0.0.0.0
migration-dir: "/home/openbar/openbar-server/schema/openbardb"
```

```bash
root@raspberrypi:/home/brian/openbar-server/cmd/openbar-server# exit
openbar@openbar:~ $ sudo systemctl enable openbar-server.service
openbar@openbar:~ $ sudo systemctl start openbar-server
```

### Enable I2C for Sequent relay hat

To enable I2C:

```
1. Run: sudo raspi-config.
2. Select Interfacing Options > I2C.
3. Select Yes when prompted to enable the I2C interface.
4. Select Yes when prompted to automatically load the I2C kernel module.
5. Select Finish.
6. Select Yes when prompted to reboot or run sudo reboot
````

### Sequent 8 relay hat testing

```
~$ git clone https://github.com/SequentMicrosystems/8relind-rpi.git
~$ cd 8relind-rpi/
~/8relind-rpi$ sudo make install
~/8relind-rpi$ 8relind 0 test
```

## Openbar-client

`~$ git clone git@github.com:cocktailrobots/openbar-client.git`
`git clone https://github.com/cocktailrobots/openbar-client.git`

```
cd ~/openbar-client
sudo bash
npm install
npm install -g serve
npx update-browserslist-db@latest
exit
```

Create a `.env` file in the root of the openbar-client directory with the following structure:

```
REACT_APP_OPENBAR_API_HOST_DEV=127.0.0.1:3099
REACT_APP_COCKTAILS_API_HOST_DEV=127.0.0.1:8675

REACT_APP_OPENBAR_API_HOST_PROD=10.42.0.1:3099
REACT_APP_COCKTAILS_API_HOST_PROD=10.42.0.1:8675

PORT=80
```

```
cd ~/openbar-client
sudo bash
npm run build
mkdir /etc/openbar-client
cp -r build/* /etc/openbar-client/
vim /etc/systemd/system/openbar-client.service
```

```
[Unit]
Description=openbar client service
After=network.target

[Service]
Type=simple
User=root
Group=root

WorkingDirectory=/etc/openbar-client
Environment="PATH=/root/.nvm/versions/node/v18.18.0/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
ExecStart=/root/.nvm/versions/node/v18.18.0/bin/serve /etc/openbar-client

LimitNOFILE=100000

Restart=always
RestartSec=1

MemoryAccounting=true
MemoryMax=90%

[Install]
WantedBy=multi-user.target
```

```
export WIFI_SSID='openbar-net-'$(uuid|cut -d '-' -f 1)
export WIFI_PASS='password1234'

nmcli con add type wifi ifname wlan0 con-name Hotspot autoconnect yes ssid "$WIFI_SSID"
nmcli con modify Hotspot 802-11-wireless.mode ap 802-11-wireless.band bg ipv4.method shared
nmcli con modify Hotspot wifi-sec.key-mgmt wpa-psk
nmcli con modify Hotspot wifi-sec.psk "$WIFI_PASS"

export CNX_UUID=$(nmcli connection show |grep -i hotspot | tr -s ' ' | cut -d ' ' -f 2)
apt purge openresolv dhcpcd5
nmcli con up uuid $CNX_UUID
```

