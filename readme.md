The program I use to manage my various audio/video (AV) devices. Basically it's a huge mishmash of logic I've built up to manage my devices. Quite specific to my setup but there are some useful bits.

Useful code:
- cmd [cecmonitor](./cmd/cecmonitor/main.go): Listen for all CEC messages on the bus and print them in a human readable format. Build with: `go build ./cmd/cecmonitor/`.
- cmd [tvman](./cmd/tvman/main.go#L143): The program I use to manage my setup. Listens for input over http and infrared. Build with: `go build ./cmd/tvman/`. Must initialize the CEC device before running.
- pkg [cec](./cec): interface with Linux kernel CEC devices (ioctl/read/write)
- pkg [ioctl](./ioctl): do linux ioctls in golang
- pkg [lirc](./lirc): interface with Linux kernel Lirc devices (ioctl/read/write)
- pkg [lgtv](./lgtv): interface with LG TVs over RS232 serial port (read/write)

The command to initialize my CEC adapter: `cec-ctl --clear -e /sys/class/drm/card1-HDMI-A-1/edid --record --osd-name "PC" --no-rc-passthrough`
- `cec-ctl`: From [v4l-utils](https://git.linuxtv.org/v4l-utils.git)
- `--clear`: Remove any existing settings on the adapter
- `-e (--phys-addr-from-edid)`: Makes sure the cec driver knows the physical location of the port it is plugged into. I can't remember if this is actually important.
- `--record`: Claim a logical address for a recording device
- `--osd-name`: Automatically respond to give osd name queries
- `--no-rc-passthrough`: This program is managing the RC input and output, don't let the cec adapter automatically do anything with remote controls. I don't remember if this is important.

I choose to make it a "record" device because I don't normally have one of those and most other devices will accept commands from a recording device. Also there can only be three playback devices connected to one CEC bus so I do not want to waste one of those logical addresses.

I use the Pulse8 CEC dongle connected to a PC:
https://www.pulse-eight.com/p/104/usb-hdmi-cec-adapter

Big shout out to https://www.cec-o-matic.com/ for helping understand the CEC commands.
