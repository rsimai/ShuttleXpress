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

Clone the repository and build the application:

```bash
git clone https://github.com/rsimai/ShuttleXpress.git
cd ShuttleXpress
go build
```

### 3. udev rules

To allow the application to access the ShuttleXpress device and create a virtual keyboard, you need to set up some `udev` rules.

1.  Copy the provided `udev` rules to `/etc/udev/rules.d/`:

    ```bash
    sudo cp udev/*.rules /etc/udev/rules.d/
    ```

2.  Reload the `udev` rules:

    ```bash
    sudo udevadm control --reload-rules && sudo udevadm trigger
    ```

3.  Create a `shuttle` group and add your user to it:

    ```bash
    sudo groupadd shuttle
    sudo usermod -aG shuttle $USER
    ```

    You will need to log out and log back in for this change to take effect.

## Usage

1.  Connect your ShuttleXpress device.

2.  Run the application:

    ```bash
    ./ShuttleXpress
    ```

The application will read the `config.json` file and start listening for events from the ShuttleXpress.

## Configuration

The `config.json` file allows you to customize the actions for each button, the jog wheel, and the shuttle ring.

Here is an example `config.json`:

```json
{
  "buttons": {
    "256": "ctrl+c",
    "257": "ctrl+v",
    "258": "ctrl+t",
    "259": "ctrl+shift+t",
    "260": "super+f4"
  },
  "jog": {
    "1": "up",
    "-1": "down"
  },
  "ring": {
    "1": "shift+up",
    "2": "shift+up",
    "3": "shift+up",
    "4": "shift+up",
    "5": "shift+up",
    "6": "shift+up",
    "7": "shift+up",
    "-1": "shift+down",
    "-2": "shift+down",
    "-3": "shift+down",
    "-4": "shift+down",
    "-5": "shift+down",
    "-6": "shift+down",
    "-7": "shift+down"
  }
}
```

### Finding Button Codes

To find the codes for the buttons on your ShuttleXpress, you can use the `evtest` utility.