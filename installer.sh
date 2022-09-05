wget https://github.com/Lucas-Samuel/HiveOS-Profit-Switcher/releases/latest/download/HiveOS-Profit-Switcher.zip

unzip HiveOS-Profit-Switcher.zip -d /usr/profit-switcher 

rm HiveOS-Profit-Switcher.zip 

chmod +x /usr/profit-switcher/switcher

(crontab -l ; echo "*/5 * * * * /usr/profit-switcher/switcher >> /usr/profit-switcher/switcher.log") | crontab -