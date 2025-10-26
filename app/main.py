import os
import json
import evdev
from evdev import UInput, ecodes

def find_device():
    devices = [evdev.InputDevice(path) for path in evdev.list_devices()]
    for device in devices:
        if "ShuttleXpress" in device.name:
            return device
    return None

def main():
    log_file = os.path.expanduser('~/shuttlexpress.log')
    with open(log_file, 'a') as f:
        f.write(f'ShuttleXpress service started\n')

    device = find_device()
    if not device:
        with open(log_file, 'a') as f:
            f.write("ShuttleXpress not found\n")
        return

    with open('config.json', 'r') as f:
        config = json.load(f)

    ui = UInput()

    for event in device.read_loop():
        if event.type == ecodes.EV_KEY and event.value == 1:  # Key press
            key_name = ecodes.KEY[event.code]
            if key_name in config:
                keys_to_press = config[key_name]
                for key in keys_to_press:
                    ui.write(ecodes.EV_KEY, ecodes.ecodes[key], 1)  # Press
                for key in reversed(keys_to_press):
                    ui.write(ecodes.EV_KEY, ecodes.ecodes[key], 0)  # Release
                ui.syn()

if __name__ == '__main__':
    main()