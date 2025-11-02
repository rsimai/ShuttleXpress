# ShuttleXpress for Linux

This is a user-space driver for the Contour Design ShuttleXpress multimedia controller on Linux. It allows you to map the buttons, jog wheel, and shuttle ring to keyboard shortcuts.

## Features

*   Map each of the 5 buttons to a custom keyboard shortcut.
*   Map the jog wheel (left/right) to keyboard shortcuts.
*   Map the shuttle ring to different keyboard shortcuts for each of its 7 positions (each direction).
*   Use uinput to create a shuttlexpress-virtual-keyboard /dev/input device.

## Installation

### 1. Dependencies

You need to have Go installed on your system. You can find installation instructions at [https://golang.org/doc/install](https://golang.org/doc/install).

### 2. Build the application

Feel free to only use the `shuttlex` binary, it should have no dependencies. To build yourself, clone the repository and build the application:

```bash
git clone https://github.com/rsimai/ShuttleXpress.git
cd ShuttleXpress
go build -o shuttlex
```

### 3. udev rules

To not run as root but allow the application to access the ShuttleXpress device and create a virtual keyboard, you need to set up some `udev` rules.

1.  Copy the provided `udev` rules to `/etc/udev/rules.d/`:

    ```bash
    sudo cp udev/*.rules /etc/udev/rules.d/
    ```

2.  Reload the `udev` rules:

    ```bash
    sudo udevadm control --reload-rules && sudo udevadm trigger
    ```

3.  Create and add your user to the `shuttle` group:

    ```bash
    sudo groupadd shuttle
    sudo usermod -aG shuttle $USER
    ```

    You will need to log out and log back in for this change to take effect.

## Usage

1.  Connect your ShuttleXpress device.

2.  Run the application:

    ```bash
    ./shuttlex
    ```

The application will read the `shuttlex.json` file and start listening for events from the ShuttleXpress. See `shuttlex --help` for more information.


## Configuration

The `shuttlex.json` file allows you to customize the actions for each button, the jog wheel, and the shuttle ring. They are sent as keyboard events. Note the rate is miliseconds, +1/-1 doesn't repeat.

Here is an example `shuttlex.json`:

```json
{
  "buttons": {
    "260": "ctrl+c",
    "261": "ctrl+v",
    "262": "super+t",
    "263": "alt+f4",
    "264": " "
  },
  "jog": {
    "1": "up",
    "-1": "down"
  },
  "ring": {
    "1": { "action": "shift+up", "rate": 0 },
    "2": { "action": "shift+up", "rate": 500 },
    "3": { "action": "shift+up", "rate": 250 },
    "4": { "action": "shift+up", "rate": 130 },
    "5": { "action": "shift+up", "rate": 70 },
    "6": { "action": "shift+up", "rate": 40 },
    "7": { "action": "shift+up", "rate": 20 },
    "-1": { "action": "shift+down", "rate": 0 },
    "-2": { "action": "shift+down", "rate": 500 },
    "-3": { "action": "shift+down", "rate": 250 },
    "-4": { "action": "shift+down", "rate": 130 },
    "-5": { "action": "shift+down", "rate": 70 },
    "-6": { "action": "shift+down", "rate": 40 },
    "-7": { "action": "shift+down", "rate": 20 }
  }
}

```

### Finding Button Codes

To find the codes for the buttons on your ShuttleXpress, you can use the `evtest` utility.