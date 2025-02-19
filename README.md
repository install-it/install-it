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
<br />
<div align="center">
  <a href="https://github.com/markmybytes/driver-box">
    <img src="https://github.com/user-attachments/assets/ea47a738-6f1e-4e8d-bde0-4f12118ff103" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">driver-box</h3>

  <p align="center">
    程式／軀動安裝工具
    <br />
    <br />
    <a href="https://github.com/markmybytes/driver-box/issues/new?labels=bug&template=bug-report---.md">Report Bug</a>
    ·
    <a href="https://github.com/markmybytes/driver-box/issues/new?labels=enhancement&template=feature-request---.md">Request Feature</a>
  </p>
</div>



<!-- ABOUT THE PROJECT -->
## Project 簡介

<p align="center">
  <img src="https://github.com/user-attachments/assets/8a29d13f-5058-4c4e-ada1-b0689add5675">
<p align="right">

driver-box 旨在加快安裝大量硬件軀動的時間。用家可以將不同類型的軀動程式加入到本程式中。之後每次只需選擇合適的軀動程式即可。<br>
除了安裝軀動程式，你亦可以利用 driver-box 來安裝其他程式或執行指令。

設計上，driver-box 是利用程式／指令執行完成後的狀態碼來判斷是否成功執行。一些程式會回應表示執行成功的狀態碼（例如 `0`），但實際上並非執行成功。

<p align="right">(<a href="#readme-top">回到最頂</a>)</p>

### 第三方工具使用
[<img src="https://img.shields.io/badge/bootstrap%20icons-7532fa?style=for-the-badge&logo=bootstrap&logoColor=white">](https://icons.getbootstrap.com/)
[<img src="https://img.shields.io/badge/go-01add8?style=for-the-badge&logo=go&logoColor=white">](https://go.dev/)
[<img src="https://img.shields.io/badge/tailwindcss-38bdf8?style=for-the-badge&logo=tailwindcss&logoColor=white">](https://tailwindcss.com/)
[<img src="https://img.shields.io/badge/vue.js-41b883?style=for-the-badge&logo=vue.js&logoColor=white">](https://vuejs.org/)
[<img src="https://img.shields.io/badge/wails-d32a2d?style=for-the-badge&logo=wails&logoColor=white">](https://wails.io/)

<p align="right">(<a href="#readme-top">回到最頂</a>)</p>



<!-- GETTING STARTED -->
## 開發

### 所需軟件

- Go https://go.dev/doc/install
- Node 22 https://nodejs.org/en/download/package-manager

### 安裝 Dependency

- Wails
  ```sh
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```
- NPM Dependencies
  ```sh
  cd ./frontend
  npm install
  ```

### 常用指令

- Debug run
  ```sh
  wails dev
  ```

- Build Executable
  ```sh
  wails build -ldflags "-X main.buildVersion=<version> -X main.binaryType=<binary type>"
  ```
  
<p align="right">(<a href="#readme-top">回到最頂</a>)</p>



<!-- USAGE EXAMPLES -->
## 使用

<img src="https://github.com/user-attachments/assets/8fb85b19-133e-4cbf-9ee4-21e5237c9089">

driver-box 最細的執行單位為執行檔。因此你可利用 driver-box 來進行安裝硬件的軀動程式外的工作，例如透過 `cmd` 或 `powershell` 來執行任何腳本（script）或殼層命令（Shell command）。

而你可以將多個軀動程式加至同一個「組合」中。所有加至同一個組合的軀動程式會一同執行。

### 加入、編輯軀動程式

<img src="https://github.com/user-attachments/assets/65d2b1fc-6138-4e81-95d6-605cecd14128">


#### 軀動路徑
軀動程式的路徑。你亦可以在此輸入 Shell command 或其他在 `PATH` 內的程式。

設計上，driver-box 是預設所有軀動程都放置 `drivers` 資料夾內，因為將軀動程式的檔案（執行檔 `.exe` 或資料夾）複製到程式的 `drivers/<分類>/` 資料夾內可以方便管理及轉移（例如複製程式到 USB 上）。
因此透過「選擇檔案」按鈕選擇的檔案將會以相對路徑表示。但你亦可以輸入絕對路徑。



#### 執行參數
[執行參數](https://en.wikipedia.org/wiki/Command-line_interface#Arguments)一般用於命令列介面（Command Line）上，以修改程式執行設定或輸入資料至程式中。

不少安裝程序都會支援以自動模式安裝（silent install），用戶無需進行任何輸入，相關程式便會自行安裝。<br>
我們十分建議輸入相關的執行參數，令加入的軀動能以自動模式安裝。

driver-box 已提供常見軀動的安裝參數。<br>

| 選項           	| 適用程式                                                                                                                                                         	|
|----------------	|------------------------------------------------------------------------------------------------------------------------------------------------------------------	|
| Intel LAN      	| [Intel® Ethernet Adapter Complete Driver Pack](https://www.intel.com/content/www/us/en/download/15084/intel-ethernet-adapter-complete-driver-pack.html)          	|
| Realtek LAN    	| [Realtek PCIe FE / GBE / 2.5G / 5G Ethernet Family Controller Software](https://www.realtek.com/Download/List?cate_id=584)                                       	|
| Nvidia Display 	| [GeForce Game Ready Driver/Nvidia Studio Driver](https://www.nvidia.com/en-us/drivers/)                                                                          	|
| AMD Display    	| [AMD Software: Adrenalin Edition](https://www.amd.com/en/support/download/drivers.html)                                                                          	|
| Intel Display  	| [Intel® Arc™ & Iris® Xe Graphics/7th-10th Gen Processor Graphics](https://www.intel.com/content/www/us/en/support/articles/000090440/graphics.html)              	|
| Intel WiFi     	| [Intel® Wireless Wi-Fi Drivers](https://www.intel.com/content/www/us/en/download/19351/intel-wireless-wi-fi-drivers-for-windows-10-and-windows-11.html)          	|
| Intel BT       	| [Intel® Wireless Bluetooth® Drivers](https://www.intel.com/content/www/us/en/download/18649/intel-wireless-bluetooth-drivers-for-windows-10-and-windows-11.html) 	|
| Intel Chipset  	| [Chipset INF Utility](https://www.intel.com/content/www/us/en/support/products/1145/software/chipset-software/intel-chipset-software-installation-utility.html)  	|
| AMD Chipset    	| [AMD Chipset Drivers](https://www.amd.com/en/support/download/drivers.html)                                                                                      	|

不在預設集上的軀動可嘗試在網上以 `軀動名稱` + `silent`／`unattended`／`command line install` 搜尋，或利用 [Silent Install Builder](https://www.silentinstall.org/) 等類似的軟件自行製作。

#### 不能同時安裝
勺選後，在使用「同步安裝」模式時，有關的軀動程式將不會在同一時間執行。

### 安裝

<img src="https://github.com/user-attachments/assets/f028262a-b39f-41d4-9969-1638ae6f6ca5">

*在所有工作執行完成前，執行狀態視窗不能夠被關閉。*

#### 關機設定
關機設定只會在所有工作執行成功及軀動安裝成功後才會執行。

#### 取消執行
只有處於「等待中」或「執行中」的工作才能取消執行。<br>
按下相關工作的「取消」按鈕即可。但注意，程式並不保證相關工作能夠被終止執行。

<p align="right">(<a href="#readme-top">回到最頂</a>)</p>



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
