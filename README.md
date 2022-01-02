# vpub

Simple message board software.

## Installation

In order to compile vpub, you need to have:
* go
* git
* make
* gcc

Here is how to build vpub:
1. `git clone git.sr.ht/~m15o/vpub`
2. `cd vpub`
3. `make`

You should now have `vpub` in `./bin/`. You can keep it there, or move it to `/usr/sbin` or anywhere else.

## systemd service

Here's an example service file. Create it on `/etc/systemd/system/vpub.service`.

```
[Install]
WantedBy=multi-user.target

[Unit]
Description=vpub

[Service]
ExecStart=/usr/sbin/vpub
```