<a id="readme-top"></a>

<!-- PROJECT SHIELDS -->
<div align="center">
  
  [![Tag][tag-shield]][tag-url]
  [![Contributors][contributors-shield]][contributors-url]
  [![Forks][forks-shield]][forks-url]
  [![Stargazers][stars-shield]][stars-url]
  [![Issues][issues-shield]][issues-url]
  [![MIT License][license-shield]][license-url]
  
</div>

<!-- PROJECT LOGO -->
<div align="center">
  <a href="https://github.com/markmybytes/driver-box">
    <img src="https://github.com/user-attachments/assets/ea47a738-6f1e-4e8d-bde0-4f12118ff103" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">driver-box</h3>

  <p align="center">
    Program/Driver Installation Tool
    <br>
    <a href="https://github.com/markmybytes/driver-box/issues/new?labels=bug&template=bug-report---.md">Report Bug</a>
    ·
    <a href="https://github.com/markmybytes/driver-box/issues/new?labels=enhancement&template=feature-request---.md">Request Feature</a>
  </p>

  <p align="center">
    <a href="https://github.com/markmybytes/driver-box//README.md">English</a>
    ·
    <a href="https://github.com/markmybytes/driver-box/readme/README.zh_Hant.md">繁體中文</a>
  </p>
</div>

<!-- ABOUT THE PROJECT -->
## About The Project

<p align="center">
  <img src="https://github.com/user-attachments/assets/35606055-7ce6-4e97-8152-a7042d7fe001" width="754" height="569">
</p>
<p align="right">

driver-box aims to speed up the installation of hardware drivers. With just a few clicks, you can install various types of drivers.<br />
Beyond drivers, driver-box can also execute system commands and install other programs, see [Usage](#usage) section for more.

</p>
<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Built With

[<img src="https://img.shields.io/badge/font%20awesome-538cd7?style=for-the-badge&logo=fontawesome&logoColor=white">](https://fontawesome.com/)
[<img src="https://img.shields.io/badge/go-01add8?style=for-the-badge&logo=go&logoColor=white">](https://go.dev/)
[<img src="https://img.shields.io/badge/tailwindcss-38bdf8?style=for-the-badge&logo=tailwindcss&logoColor=white">](https://tailwindcss.com/)
[<img src="https://img.shields.io/badge/vue.js-41b883?style=for-the-badge&logo=vue.js&logoColor=white">](https://vuejs.org/)
[<img src="https://img.shields.io/badge/wails-d32a2d?style=for-the-badge&logo=wails&logoColor=white">](https://wails.io/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

### Prerequisites

- Go ≥1.23 https://go.dev/doc/install
- Node 22 https://nodejs.org/en/download/package-manager

### Setup

#### Install dependencies

- Wails
  ```sh
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```
- NPM Dependencies
  ```sh
  cd ./frontend
  npm install
  ```

#### Commands

- Debug run

  ```sh
  wails dev
  ```

- Build Executable
  ```sh
  wails build

  # build with version number set
  wails build -ldflags "-X main.buildVersion=<version number>"
  ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->
## Usage

### Managing Drivers

<img src="https://github.com/user-attachments/assets/909dcbbd-9b02-4c06-941e-a77035e1250f" width="754" height="569">

When you add a driver to driver-box, there are three types you can choose from: `network`, `display`, and `miscellaneous`. Only `miscellaneous` allow multiple selection.

#### Path

The path to the driver installer.

driver-box are designed with portability in mind, therefore 
  1. It is recommended to place all your driver installer under the `driver/<category>` folder to facilitate management and transfer (e.g., copying the program to a USB drive).
  
  2. Paths selected via the "Select File" button will be a relative paths to driver-box's executable path. However, nothing prevents you from entering an absolute paths manually.

Other than installing drivers, driver-box can do more than that.

##### Execute commands

You can execute commands available in the OS `PATH` variable.

For CMD, you can execute commands by entering `cmd` in the path field, and `/c,<command>` in the install option field. Then it is equivlent to:

```batch
cmd /c <command>
```

For Powershell, you can execute commands by entering `powershell` in the path field, and `-Command,<command>` in the install option field. Then it is equivlent to:

```batch
powershell -Command <command>
```

##### Install programs other than drivers

Program installer usually provides a slient install options like driver installers. For example, Steam support silent install by supplying `/S` option when you executing `SteamSetup.exe`. Explore yourself and turn driver-box to your PC setup toolbox :).

#### Install Option

[CLI Option](https://en.wikipedia.org/wiki/Command-line_interface#Arguments) is to control the execution behaviour, or input data into the program.

Many installation executables support silent installation, where the program will be installed automatically without any interactions.<br>
It is highly recommended to enter the appropriate execution parameters to allow drivers to be installed in silent mode to maximise functionality of driver-box.

driver-box provides installation parameters preset for common drivers:<br>

| Option         | Applicable Program                                                                                                                                               |
| -------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Intel LAN      | [Intel® Ethernet Adapter Complete Driver Pack](https://www.intel.com/content/www/us/en/download/15084/intel-ethernet-adapter-complete-driver-pack.html)          |
| Realtek LAN    | [Realtek PCIe FE / GBE / 2.5G / 5G Ethernet Family Controller Software](https://www.realtek.com/Download/List?cate_id=584)                                       |
| Nvidia Display | [GeForce Game Ready Driver/Nvidia Studio Driver](https://www.nvidia.com/en-us/drivers/)                                                                          |
| AMD Display    | [AMD Software: Adrenalin Edition](https://www.amd.com/en/support/download/drivers.html)                                                                          |
| Intel Display  | [Intel® Arc™ & Iris® Xe Graphics/7th-10th Gen Processor Graphics](https://www.intel.com/content/www/us/en/support/articles/000090440/graphics.html)              |
| Intel WiFi     | [Intel® Wireless Wi-Fi Drivers](https://www.intel.com/content/www/us/en/download/19351/intel-wireless-wi-fi-drivers-for-windows-10-and-windows-11.html)          |
| Intel BT       | [Intel® Wireless Bluetooth® Drivers](https://www.intel.com/content/www/us/en/download/18649/intel-wireless-bluetooth-drivers-for-windows-10-and-windows-11.html) |
| Intel Chipset  | [Chipset INF Utility](https://www.intel.com/content/www/us/en/support/products/1145/software/chipset-software/intel-chipset-software-installation-utility.html)  |
| AMD Chipset    | [AMD Chipset Drivers](https://www.amd.com/en/support/download/drivers.html)                                                                                      |

For drivers that are not in the preset, you can try searching online for `driver name` + `silent`/`unattended`/`command line install`.

### Installation

Select all the suitable driver and click `Execute`. A popup will be displayed for execution status.

> Note: driver-box uses the exit status code to determine success. Some programs may return 0 (success) even if the installation failed.

#### Shutdown Option

Shutdown option will only be applied after all executions are successfully completed.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[tag-url]: https://github.com/markmybytes/driver-box/releases
[tag-shield]: https://img.shields.io/github/v/tag/markmybytes/driver-box?style=for-the-badge&label=LATEST&color=%23B1B1B1
[contributors-shield]: https://img.shields.io/github/contributors/markmybytes/driver-box.svg?style=for-the-badge
[contributors-url]: https://github.com/markmybytes/driver-box/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/markmybytes/driver-box.svg?style=for-the-badge
[forks-url]: https://github.com/markmybytes/driver-box/network/members
[stars-shield]: https://img.shields.io/github/stars/markmybytes/driver-box.svg?style=for-the-badge
[stars-url]: https://github.com/markmybytes/driver-box/stargazers
[issues-shield]: https://img.shields.io/github/issues/markmybytes/driver-box.svg?style=for-the-badge
[issues-url]: https://github.com/markmybytes/driver-box/issues
[license-shield]: https://img.shields.io/github/license/markmybytes/driver-box.svg?style=for-the-badge
[license-url]: https://github.com/markmybytes/driver-box/blob/master/LICENSE.txt

