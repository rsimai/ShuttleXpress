#!/bin/bash

# It's a good practice to use absolute paths in udev scripts
logger "ShuttleXpress plugged in"
/usr/bin/python3 /home/robert/git/ShuttleXpress/app/main.py
