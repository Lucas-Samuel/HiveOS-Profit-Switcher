# HiveOS Profit Switcher
This is a app inspired by [TheRetroMike's HiveOSProfitSwitcher](https://github.com/TheRetroMike/HiveOSProfitSwitcher) that can connect to your Hive OS account and profit switch your mining rigs based on WhatToMine calculations and using a configuration that you setup.

Make sure you configure overclock defaults for each algorithm to optimize your cards and prevent crashes. This Profit switcher doesn't apply overclocks itself, but overclocks will get reapplied by hive os if you configure algo default ones when this tool applies the flightsheet changes.

## How To Use
```bash
# Download the app
$ wget https://github.com/Lucas-Samuel/HiveOS-Profit-Switcher/releases/latest/download/HiveOS-Profit-Switcher.zip

# Extract and remove the zip file
$ unzip HiveOS-Profit-Switcher.zip -d /usr/profit-switcher && rm HiveOS-Profit-Switcher.zip

# Give permission to execute app
$ chmod +x /usr/profit-switcher/switcher

# Edit the configs.json whit yours keys and save
$ nano /usr/profit-switcher/configs.json

# Add the following line to the end of the crontab file and save
# */5 * * * * /usr/profit-switcher/switcher >> /usr/profit-switcher/switcher.log
$ crontab -e

# You can execute the app manually by running
$ /usr/profit-switcher/switcher
```

## Example Config

- api_key = Hive OS API Key (You can get this by logging into your account, clicking your username at the top right, selecting the Sessions tab, and then creating a new API token)
- farm_id = Hive OS Farm ID (You can get this from your farm URL)
- coin_diference = This is the profit percentage difference that is needed in order for the coin to switch.
- wtm_endpoint = WhatToMine JSON (Go to WhatToMine, fill in your hashrates, hit calculate, then click JSON at the top, and copy the URL from the address bar)

Note: Set "Average for Revenue" to "Current Values" in WhatToMine before copying the URL if you want real-time profitability that auto-adjusts for real-time difficulty changes.

```json
{
  "api_key": "xxxx",
  "farm_id": "xxxx",
  "coin_diference": "5",
  "workers": [
    {
      "name": "RIG_01",
      "wtm_endpoint": "https://whattomine.com/coins.json?...",
      "coins": [
        {
          "tag": "ETH",
          "fs": "ETH_RIG_01"
        },
        {
          "tag": "NEOX",
          "fs": "NEOX_RIG_01"
        },
        {
          "tag": "RVN",
          "fs": "RVN_RIG_01"
        }
      ]
    },
    {
      "name": "RIG_02",
      "wtm_endpoint": "https://whattomine.com/coins.json?...",
      "coins": [
        {
          "tag": "NEOX",
          "fs": "NEOX_RIG_02"
        },
        {
          "tag": "RVN",
          "fs": "RVN_RIG_02"
        }
      ]
    }
  ]
}
```

## Buy me a coffee
If this project help you, you can give me a cup of coffee :)

ETH Mainnet, Matic and BSC

```
0xEEdb9541a2E9551677e22037736BEfb7ddE9c467
```