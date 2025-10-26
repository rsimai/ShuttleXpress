# ShuttleXpress Background Service

This application runs in the background and starts automatically when a ShuttleXpress USB device is plugged in.

## Installation

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd ShuttleXpress
    ```

2.  **Install dependencies:**
    ```bash
    pip install -r requirements.txt
    ```

3.  **Install the udev rule:**
    ```bash
    sudo cp udev/99-shuttlexpress.rules /etc/udev/rules.d/
    ```

4.  **Reload the udev rules:**
    ```bash
    sudo udevadm control --reload-rules
    sudo udevadm trigger
    ```

## Usage

1.  **Permissions:**
    This application requires access to the input device. You can either run it as root or add your user to the `input` group:
    ```bash
    sudo adduser $USER input
    ```
    You will need to log out and log back in for this change to take effect.

2.  **Configuration:**
    The button mappings are defined in `config.json`. You can edit this file to customize the keyboard shortcuts. The key names correspond to the `ecodes` in the `evdev` library.

3.  **Running the application:**
    Plug in your ShuttleXpress device. The application will start automatically. A log file named `shuttlexpress.log` will be created in your home directory.
