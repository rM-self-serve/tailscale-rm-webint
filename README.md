![Static Badge](https://img.shields.io/badge/reMarkable-v3.9.4-green)
[![rm1](https://img.shields.io/badge/rM1-supported-green)](https://remarkable.com/store/remarkable)
[![rm2](https://img.shields.io/badge/rM2-supported-green)](https://remarkable.com/store/remarkable-2)
[![Discord](https://img.shields.io/discord/385916768696139794.svg?label=reMarkable&logo=discord&logoColor=ffffff&color=7389D8&labelColor=6A7EC2)](https://discord.gg/ATqQGfu)

# Deprecated

Use webinterface-localhost from [Toltec](https://toltec-dev.org/) instead.


# tailscale-rm-webint

Make the ReMarkable Tablet's [web interface](https://remarkable.guide/tech/usb-web-interface.html) available over tailscale.


### Limitations 

Without additional programs, the web interface will only be available over tailscale while the device is plugged in and the web interface is enabled/reachable at 10.11.99.1.
To ensure the web interface is always available, use [webinterface-onboot](https://github.com/rM-self-serve/webinterface-onboot).

Drag and drop does not work well on mobile, though it is simple to add an [upload button](https://github.com/rM-self-serve/upload_button).


## Installation/Removal

**It is required to install via the [toltec package manager](https://toltec-dev.org/).** 

```
$ opkg update
$ opkg install tailscale-rm-webint
$ opkg remove tailscale-rm-webint
```

## Usage

### To use tailscale-rm-webint, run:

`$ systemctl enable --now tailscale-rm-webint`

To view the web interface, type the remarkable's tailscale address in the browser.

## Config

The default address this reverse proxy runs on is localhost:80. In order to run on a different port, such as 8080, a config file can be created at
`/home/root/.config/tailscale-rm-webint/config`
with the contents:

```toml
port=8080
```

Then run:

```bash
$ systemctl restart tailscale-rm-webint
```
