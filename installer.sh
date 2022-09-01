wget https://github.com/Lucas-Samuel/HiveOS-Profit-Switcher/releases/latest/download/HiveOS-Profit-Switcher.zip && unzip HiveOS-Profit-Switcher.zip -d /usr/profit-switcher && rm HiveOS-Profit-Switcher.zip && cd /usr/profit-switcher && chmod +x switcher

printf "\n*/5 * * * * /usr/profit-switcher/switcher >> /usr/profit-switcher/switcher.log\n" >> /hive/etc/crontab.root