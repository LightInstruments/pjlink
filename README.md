# pjlink
Go implementation of PJLink. Library only. Based on: github.com/byuoitav/pjlink-microservice.  
[![Apache 2 License](https://img.shields.io/hexpm/l/plug.svg)](https://raw.githubusercontent.com/byuoitav/pjlink-microservice/master/LICENSE)

## Setup
// TODO: This need to be configurable per projector:  
You need to set the `PJLINK_PORT` and `PJLINK_PASS` environment variables before running.

## Usage
Send a `POST` request to the `/command` endpoint with a body similar to the following:
```
{
    "address": "10.66.9.14",
    "port": "4352",
    "password": "sekret",
    "class": "1",
    "command": "INST",
    "parameter": "?"
}
```