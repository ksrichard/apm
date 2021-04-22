apm (Arduino Package Manager)
---

A package manager for Arduino projects to make them portable (no more manual download etc...).

This tool fully includes the official `arduino-cli` (it won't be installed nor need to be installed)
so it is automatically compatible with any kind of Arduino projects/sketches.

### Features
```bash
A package manager for Arduino projects.
The official arduino-cli packages are used to perform actions.

Usage:
  apm [command]

Available Commands:
  add         Adding new libraries to the project
  help        Help about any command
  init        Init APM project
  install     Install dependencies of project
  remove      Remove library from the project

Flags:
  -h, --help                 help for apm
  -p, --project-dir string   Project directory to use (default "/Users/klavorar/Documents/Arduino/temp_sensor")

Use "apm [command] --help" for more information about a command.
```

### Showcase
TODO

### Installation
#### Install on Linux/MacOS
Run the following in your terminal to install the latest version of `apm`
```shell script
curl -s https://raw.githubusercontent.com/ksrichard/apm/master/install.sh | sh
```
#### Install on Windows
Download the latest version of `apm` from https://github.com/ksrichard/apm/releases/latest for Windows (`apm_windows_amd64.exe`), rename it to apm.exe and put it on path.

### How it works
Every `apm` based project must have a file called `apm.json` in the project root (it can be create by running `apm init`)
This configuration file is containing all the information that an Arduino project needs.

`NOTE on versioning` - if you would like to use always the latest version, please use `latest` in any package version and always latest will be used!  

`apm.json` structure:
- `board` - (Optional) you can select here the package/architecture of the board you will use, it will be automatically installed
    - `package` - Arduino core package name
    - `architecture` -  Architecture of Arduino core package
    - `version` - Version of core package (`latest` for always latest version)
    - `board_manager_url` - (Optional) Additional Board Manager URL if needed for the board core package to be installed
- `dependencies` - (Optional, if empty, no dependencies will be installed of course)
contains all Arduino Library dependencies that the actual project needs (if any Version mismatch will be in place, process will be stopped) 
    - `library` - Arduino Library name
    - `version` - Arduino Library version
    - `git` - (Optional - if it's set, do not set `library` and `version`) install library from git repository
    - `zip` - (Optional - if it's set, do not set `library` and `version`) install library from local zip file
    
Example `apm.json`:
```json
{
    "board": {
        "package": "esp8266",
        "architecture": "esp8266",
        "version": "latest",
        "board_manager_url": "https://arduino.esp8266.com/stable/package_esp8266com_index.json"
    },
    "dependencies": [
        {
            "library": "HomeKit-ESP8266",
            "version": "1.2.0"
        },
        {
            "library": "OneWire",
            "version": "latest"
        },
        {
            "library": "DallasTemperature",
            "version": "latest"
        },
        {
            "git": "https://github.com/jandrassy/ArduinoOTA"
        },
        {
            "zip": "ESP8266NetBIOS.zip"
        }
    ]
}
```
 
