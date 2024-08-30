# Replay - RF433T/RF433R

This script is used to send message to a RF433R device from a RF433T device. It's the environment for a replay attack by a Flipper Zero 

## Prerequisites
<p align="left">
<a href="https://go.dev/" target="_blank" rel="noreferrer"> <img src="https://cdn.icon-icons.com/icons2/2107/PNG/512/file_type_go_gopher_icon_130571.png" alt="go" width="60" height="60"/> </a>
<a href="https://flipperzero.one/" target="_blank" rel="noreferrer"> <img src="https://user-images.githubusercontent.com/29007647/182851959-afaa1367-9f4d-46c8-92af-aa5ff70fca64.png" alt="flipper zero" width="90" height="60"/> </a>
</p>

## Usage
1. Choose if you want to use rolling code or not
2. Run the go scrypt on sudo mode
   ```bash
   sudo go run .
   ```

## How it works
A message is sent in Serial to the RF433T device. The RF433R device receive the message and sent it by Serial.
The Flipper Zero device can record the message and replay it.

Additionally, the script can use a rolling code to make the replay attack impractical.

## Authors
- **Rémi Heredero** - _Initial work_ - [Klagarge](https://github.com/Klagarge)