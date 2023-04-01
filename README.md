# Stay Focused

Stay Focused is a background process for staying focused during your working day. It helps you eliminate distractions by automatically killing specified applications during your working hours.

## Features

- Cross-platform support (macOS, Windows, Linux)
- Customizable configuration file
- Automatically kills distracting applications during working hours
- Displays a notification before killing a process (Windows only)

## Installation

1. Install [Go](https://golang.org/doc/install) if you haven't already.
2. Clone this repository:
```
git clone https://github.com/yourusername/stay-focused.git
```
3. Change to the project directory:
```
cd stay-focused
```
4. Build the project:
```
go build
```

## Usage

1. Run the `stay-focused` executable:
```
./stay-focused
```
2. On the first run, the application will create a configuration file at `~/.config/.stay_focused.json` if it doesn't exist. You can edit this file to customize the working hours and the list of distracting applications.

## Configuration

The configuration file is located at `~/.config/.stay_focused.json`. It has the following structure:

```json
{
    "time": {
        "start": "09:00",
        "end": "18:00"
    },
    "applications": [
        "Skype",
        ...,
        "Telegram"
    ]
}
```

- `time.start`: The start of your working hours (24-hour format).
- `time.end`: The end of your working hours (24-hour format).
- `applications`: An array of application names or paths that you want to be killed during your working hours. The application will use regular expressions to match the process names.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License.