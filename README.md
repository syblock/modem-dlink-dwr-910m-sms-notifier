# Modem SMS Notifier for D-Link DWR-910M

`modem-dlink-dwr-910m-sms-notifier` is a simple Go script that checks for new SMS messages from a D-Link DWR-910M modem and sends desktop notifications. It prevents duplicate notifications and automatically cleans up old messages.

-----

## 🚀 Features

  * **Desktop Notifications**: Sends notifications for new SMS.
  * **Reads Messages**: Tracks the last read SMS to avoid duplicates.
  * **Automatic Cleanup**: Deletes old SMS messages from the modem.

-----

## 🛠️ Requirements

  * **Go** (v1.16 or newer)
  * **D-Link DWR-910M modem**

-----

## ⚙️ Usage

1.  **Clone the repository**:

    ```bash
    git clone https://github.com/syblock/modem-dlink-dwr-910m-sms-notifier.git
    cd modem-dlink-dwr-910m-sms-notifier
    ```

2.  **Set up environment variables**:
    Create a `.env` file or set variables directly.

      * **MODEM\_API\_BASE\_URL**: `http://192.168.0.1` (the modem's IP address)
      * **MODEM\_STORAGE\_DIR**: A directory path to store the last check timestamp.

    **Example `.env` file:**

    ```ini
    MODEM_API_BASE_URL="http://192.168.0.1"
    MODEM_STORAGE_DIR="/home/user/modem-data"
    ```

    **Example command-line usage:**

    ```bash
    MODEM_API_BASE_URL="http://192.168.0.1" MODEM_STORAGE_DIR="/home/user/modem-data" go run main.go
    ```

-----

## 📦 Building the Executable

Build a standalone executable for easier deployment without the Go runtime.

1.  **Build**: `go build -o modem-notifier main.go`
2.  **Run**:
    ```bash
    MODEM_API_BASE_URL="http://192.168.0.1" MODEM_STORAGE_DIR="/home/user/modem-data" ./modem-notifier
    ```

## ⏰ Scheduled Task

For continuous monitoring, use a **cron job**. This example runs the script every 5 minutes.

```bash
*/5 * * * * MODEM_API_BASE_URL="http://192.168.0.1" MODEM_STORAGE_DIR="/home/user/modem-data" /path/to/your/executable/modem-notifier
```

## ⚖️ License

This project is licensed under the [MIT License](LICENSE).