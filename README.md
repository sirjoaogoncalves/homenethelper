# HomeNetHelper

HomeNetHelper is a network monitoring tool written in Go that provides real-time information about devices connected to your home network. It captures and analyzes network traffic to give you insights into device usage, connection types, and data transfer rates.

## Features

- Monitor multiple network interfaces simultaneously (e.g., Ethernet and Wi-Fi)
- Display real-time data transfer rates for each device
- Show connection types (Wi-Fi or Ethernet) for each device
- Identify top protocols and ports used by each device
- Provide cumulative data transfer totals

## Prerequisites

Before you can run HomeNetHelper, you need to have the following installed on your system:

- Go (version 1.16 or later)
- libpcap development files

### Installing libpcap

On Ubuntu or Debian-based systems:
```bash
sudo apt-get update
sudo apt-get install libpcap-dev
```

On Fedora or RHEL-based systems:
```bash
sudo dnf install libpcap-devel
```

On macOS (using Homebrew):
```bash
brew install libpcap
```

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/sirjoaogoncalves/homenethelper.git
   cd homenethelper
   ```

2. Build the project:
   ```bash
   go build -o homenethelper cmd/homenethelper/main.go
   ```

## Usage

To run HomeNetHelper, you need to specify the network interfaces you want to monitor. You also need to run it with sudo privileges to capture network packets.

```bash
sudo ./homenethelper -interfaces=eth0,wlan0 -refresh=5
```

Replace `eth0,wlan0` with the names of the network interfaces you want to monitor. You can find your interface names by running `ip link show` or `ifconfig` on Linux, or `ipconfig` on Windows.

The `-refresh` flag sets the refresh rate in seconds. In this example, it's set to 5 seconds.

## Output

The program will display a table with the following columns:

- IP Address: The IP address of the device
- Device Name: The hostname of the device (if available)
- Connection: The type of connection (Wi-Fi or Ethernet)
- Interface: The network interface the device is connected to
- Down Speed: Current download speed
- Up Speed: Current upload speed
- Total Down: Total data downloaded
- Total Up: Total data uploaded
- Top Protocol: Most used protocol and its data usage
- Top Port: Most used port and its data usage

## Troubleshooting

If you encounter any issues:

1. Check the `homenethelper.log` file for error messages.
2. Ensure you're running the program with sudo privileges.
3. Verify that you've correctly specified your network interface names.

## Contributing

Contributions to HomeNetHelper are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Disclaimer

This tool is for educational and personal use only. Always respect privacy laws and obtain necessary permissions before monitoring network traffic.
