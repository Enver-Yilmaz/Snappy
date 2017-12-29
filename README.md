# Snappy
A minimalistic media center focused on RPI, written in Go

![main menu of Snappy](https://i.imgur.com/Rvni8mL.png)
![some Snappy generated results](https://i.imgur.com/PF8J217.png)

## How to use it
1) download the latest release or compile the source yourself
2) create a config, see the config section, file with your real debrid and alluc keys (more options coming soon)
3) run it in a linux distro of your choice (I recommend dietpi for the rpi, and I'll have some automation soon)

## Creating a config
1) make a file in the same folder as the executable named config.cfg
2) give it the following parameters (note that the tmdb key isn't actually implemented yet, and therefore not required)
```
AllucKey="youralluckeyhere"
RealDebridKey="yourRDkeyhere"
TmdbKey="yourtmdbkeyhere"
```
[find alluc key here](https://accounts.alluc.ee/)

[find RD key here](https://real-debrid.com/)

## Controlling Snappy
You can either use the [web remote](http://remote.squared.technology), or a keyboard (simply arrow keys and enter)
