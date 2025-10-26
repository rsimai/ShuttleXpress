# ShuttleXpress Background Service

This application runs in the background and starts automatically when a ShuttleXpress USB device is plugged in.

## Installation

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd ShuttleXpress
    ```

2.  **Install the udev rule:**
    ```bash
    sudo cp udev/99-shuttlexpress.rules /etc/udev/rules.d/
    ```

3.  **Reload the udev rules:**
    ```bash
    sudo udevadm control --reload-rules
    sudo udevadm trigger
    ```

## Usage

Plug in your ShuttleXpress device. A log file named `shuttlexpress.log` will be created in your home directory, and a new line will be added to it each time the device is connected.
