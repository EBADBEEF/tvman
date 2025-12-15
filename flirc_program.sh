#!/bin/sh
#from https://kiljan.org/2019/03/30/controlling-android-tv-and-kodi-with-logitech-harmony-and-flirc/
flirc_util() {
  ./Flirc-3.27.16/flirc_util "$@" || exit
}

flirc_util wait
flirc_util format

# TODO: use ir-ctl + keymap file to automate the whole thing...

echo -n "PowerOn: " && flirc_util record_api 1 102
echo -n "PowerOff: " && flirc_util record_api 50 102
echo -n "Up: " && flirc_util record_api 66 102
echo -n "Down: " && flirc_util record_api 67 102
echo -n "Left: " && flirc_util record_api 68 102
echo -n "Right: " && flirc_util record_api 69 102
echo -n "Channel up: " && flirc_util record_api 156 102
echo -n "Channel down: " && flirc_util record_api 157 102
echo -n "OK: " && flirc_util record_api 65 102
echo -n "Back: " && flirc_util record_api 70 102
echo -n "Home: " && flirc_util record_api 8 40
echo -n "Menu: " && flirc_util record_api 64 102
echo -n "Play: " && flirc_util record_api 176 102
echo -n "Pause: " && flirc_util record_api 176 102
echo -n "Stop: " && flirc_util record_api 183 102
echo -n "Skip back: " && flirc_util record_api 182 102
echo -n "Skip forward: " && flirc_util record_api 181 102
echo -n "Volume up: " && flirc_util record_api 233 102
echo -n "Volume down: " && flirc_util record_api 234 102
echo -n "1: " && flirc_util record_api 0 30
echo -n "2: " && flirc_util record_api 0 31
echo -n "3: " && flirc_util record_api 0 32
echo -n "4: " && flirc_util record_api 0 33
echo -n "5: " && flirc_util record_api 0 34
echo -n "6: " && flirc_util record_api 0 35
echo -n "7: " && flirc_util record_api 0 36
echo -n "8: " && flirc_util record_api 0 37
echo -n "9: " && flirc_util record_api 0 38
echo -n "0: " && flirc_util record_api 0 39

echo ""
echo "Onkyo replacement remote STRM BOX keys..."
echo ""
echo -n "PowerOn (STRMBOX): " && flirc_util record_api 1 102
#STRMBOX = necx, 0xd26d0c
echo -n "Up: " && flirc_util record_api 66 102
#UP = nec32, 0x87ee000b
echo -n "Down: " && flirc_util record_api 67 102
#DOWN = nec32, 0x87ee000d
echo -n "Left: " && flirc_util record_api 68 102
#LEFT = nec32, 0x87ee0008
echo -n "Right: " && flirc_util record_api 69 102
#RIGHT = nec32, 0x87ee0007
echo -n "OK: " && flirc_util record_api 65 102
#OK (ENTER and SOURCE) = nec32, 0x87ee0004
echo -n "Back: " && flirc_util record_api 70 102
#BACK (SETUP) = nec32, 0x87ee0002

#echo -n "Keypad Enter (text input): " && flirc_util record_api 0 40
#echo -n "Backspace: " && flirc_util record_api 0 42
#echo -n "Rewind: " && flirc_util record_api 180 102
#echo -n "Fast forward: " && flirc_util record_api 179 102
#echo -n "Kodi Info (I): " && flirc_util record_api 0 12
#echo -n "Kodi Fullscreen/GUI switch (Tab)" flirc_util record_api 0 43
#echo -n "Kodi Show Subtitles (T)" flirc_util record_api 0 23
#echo -n "Kodi Next Subtitles (L)" flirc_util record_api 0 15
