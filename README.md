# Service TTL

![build](https://github.com/bertrandmartel/service-ttl/workflows/build/badge.svg) [![License](http://img.shields.io/:license-mit-blue.svg)](LICENSE.md)

API server with the ability to launch a background service (or a set of commands) for a fixed period of time.

For instance a `POST /start` will :

* launch the service (or the set of commands) if not already running
* reset the expiration period for this service

This is useful for service you don't want to be up all the time to save CPU and memory usage. For example, a specific usecase is a camera streaming service that should be up only when at least one client is active (or connected). Successive call to the API reset the expiration time.

## Usage

* download latest release

* run

```bash
./service-ttl -port=6005 -config=$(pwd)/config.json -timeoutMinutes=30
```

## Configuration file

Example for the json configuration file :

```json
{
  "version": "0.1",
  "port": 6005,
  "serverPath": "http://localhost",
  "timeoutMinutes": 30,
  "commands": [{
    "binary": "v4l2-ctl",
    "params": ["--set-fmt-video=width=1920,height=1080,pixelformat=H264,field=4"]
  },{
    "binary":"v4l2-ctl",
    "params": ["--set-ctrl","video_bitrate=15000000"]
  },{
    "binary":"v4l2-ctl",
    "params": ["--set-ctrl","rotate=270"]
  },{
    "binary":"gst-launch-1.0",
    "params": [
      "v4l2src","device=/dev/video0", "!",
      "video/x-raw,width=1280,height=720,framerate=30/1", "!",
      "omxh264enc", "!",
      "rtph264pay", "config-interval=1", "pt=96", "!",
      "udpsink", "sync=false", "host=10.8.1.1", "port=8004"
    ]
  }]
}
```

On /start the set of commands is executed. Program waits for execution of each command.

## Install with systemd service

service-ttl.service

```
[Unit]
Description=TTL Service
After=network-online.target
 
[Service]
ExecStart=/bin/service-ttl -port=6005 -config=/etc/service-ttl/config.json -timeoutMinutes=30
StandardOutput=inherit
StandardError=inherit
Restart=always
User=pi
 
[Install]
WantedBy=multi-user.target
```

Then : 
```
cp service-ttl.service /lib/systemd/system/
sudo systemctl enable service-ttl.service
```

## Development

```bash
git clone git@github.com:bertrandmartel/service-ttl.git
make run
```

## Open Source components

* [echo](https://echo.labstack.com/)