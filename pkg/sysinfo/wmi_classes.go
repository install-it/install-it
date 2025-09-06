package sysinfo

/*
The Win32_Processor WMI class represents a device that can interpret a sequence of instructions on a computer running on a Windows operating system.

See: https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/win32-processor
*/
type Win32_Processor struct {
	AddressWidth                            uint16
	Architecture                            uint16
	AssetTag                                string
	Availability                            uint16
	Caption                                 string
	Characteristics                         uint32
	ConfigManagerErrorCode                  uint32
	ConfigManagerUserConfig                 bool
	CpuStatus                               uint16
	CreationClassName                       string
	CurrentClockSpeed                       uint32
	CurrentVoltage                          uint16
	DataWidth                               uint16
	Description                             string
	DeviceID                                string
	ErrorCleared                            bool
	ErrorDescription                        string
	ExtClock                                uint32
	Family                                  uint16
	InstallDate                             string
	L2CacheSize                             uint32
	L2CacheSpeed                            uint32
	L3CacheSize                             uint32
	L3CacheSpeed                            uint32
	LastErrorCode                           uint32
	Level                                   uint16
	LoadPercentage                          uint16
	Manufacturer                            string
	MaxClockSpeed                           uint32
	Name                                    string
	NumberOfCores                           uint32
	NumberOfEnabledCore                     uint32
	NumberOfLogicalProcessors               uint32
	OtherFamilyDescription                  string
	PartNumber                              string
	PNPDeviceID                             string
	PowerManagementCapabilities             []uint16
	PowerManagementSupported                bool
	ProcessorId                             string
	ProcessorType                           uint16
	Revision                                uint16
	Role                                    string
	SecondLevelAddressTranslationExtensions bool
	SerialNumber                            string
	SocketDesignation                       string
	Status                                  string
	StatusInfo                              uint16
	Stepping                                string
	SystemCreationClassName                 string
	SystemName                              string
	ThreadCount                             uint32
	UniqueId                                string
	UpgradeMethod                           uint16
	Version                                 string
	VirtualizationFirmwareEnabled           bool
	VMMonitorModeExtensions                 bool
	VoltageCaps                             uint32
}

/*
The Win32_MotherboardDevice WMI class represents a device that contains the central components of the Windows computer system.

See: https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/win32-motherboarddevice
*/
type Win32_MotherboardDevice struct {
	Availability                uint16
	Caption                     string
	ConfigManagerErrorCode      uint32
	ConfigManagerUserConfig     bool
	CreationClassName           string
	Description                 string
	DeviceID                    string
	ErrorCleared                bool
	ErrorDescription            string
	InstallDate                 string
	LastErrorCode               uint32
	Name                        string
	PNPDeviceID                 string
	PowerManagementCapabilities []uint16
	PowerManagementSupported    bool
	PrimaryBusType              string
	RevisionNumber              string
	SecondaryBusType            string
	Status                      string
	StatusInfo                  uint16
	SystemCreationClassName     string
	SystemName                  string
}

/*
The Win32_BaseBoard WMI class represents a baseboard, which is also known as a motherboard or system board.

See: https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/win32-baseboard
*/
type Win32_BaseBoard struct {
	Caption                 string
	ConfigOptions           []string
	CreationClassName       string
	Depth                   float32
	Description             string
	Height                  float32
	HostingBoard            bool
	HotSwappable            bool
	InstallDate             string
	Manufacturer            string
	Model                   string
	Name                    string
	OtherIdentifyingInfo    string
	PartNumber              string
	PoweredOn               bool
	Product                 string
	Removable               bool
	Replaceable             bool
	RequirementsDescription string
	RequiresDaughterBoard   bool
	SerialNumber            string
	SKU                     string
	SlotLayout              string
	SpecialRequirements     bool
	Status                  string
	Tag                     string
	Version                 string
	Weight                  float32
	Width                   float32
}

/*
The Win32_PhysicalMemory WMI class represents a physical memory device located on a computer system and available to the operating system.

See: https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/win32-physicalmemory
*/
type Win32_PhysicalMemory struct {
	Attributes           uint32
	BankLabel            string
	Capacity             uint64
	Caption              string
	ConfiguredClockSpeed uint32
	ConfiguredVoltage    uint32
	CreationClassName    string
	DataWidth            uint16
	Description          string
	DeviceLocator        string
	FormFactor           uint16
	HotSwappable         bool
	InstallDate          string
	InterleaveDataDepth  uint16
	InterleavePosition   uint32
	Manufacturer         string
	MaxVoltage           uint32
	MemoryType           uint16
	MinVoltage           uint32
	Model                string
	Name                 string
	OtherIdentifyingInfo string
	PartNumber           string
	PositionInRow        uint32
	PoweredOn            bool
	Removable            bool
	Replaceable          bool
	SerialNumber         string
	SKU                  string
	SMBIOSMemoryType     uint32
	Speed                uint32
	Status               string
	Tag                  string
	TotalWidth           uint16
	TypeDetail           uint16
	Version              string
}

/*
The Win32_VideoController WMI class represents the capabilities and management capacity of the video controller on a computer system running Windows.

See: https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/win32-videocontroller
*/
type Win32_VideoController struct {
	AcceleratorCapabilities      []uint16
	AdapterCompatibility         string
	AdapterDACType               string
	AdapterRAM                   uint32
	Availability                 uint16
	CapabilityDescriptions       []string
	Caption                      string
	ColorTableEntries            uint32
	ConfigManagerErrorCode       uint32
	ConfigManagerUserConfig      bool
	CreationClassName            string
	CurrentBitsPerPixel          uint32
	CurrentHorizontalResolution  uint32
	CurrentNumberOfColors        uint64
	CurrentNumberOfColumns       uint32
	CurrentNumberOfRows          uint32
	CurrentRefreshRate           uint32
	CurrentScanMode              uint16
	CurrentVerticalResolution    uint32
	Description                  string
	DeviceID                     string
	DeviceSpecificPens           uint32
	DitherType                   uint32
	DriverDate                   string
	DriverVersion                string
	ErrorCleared                 bool
	ErrorDescription             string
	ICMIntent                    uint32
	ICMMethod                    uint32
	InfFilename                  string
	InfSection                   string
	InstallDate                  string
	InstalledDisplayDrivers      string
	LastErrorCode                uint32
	MaxMemorySupported           uint32
	MaxNumberControlled          uint32
	MaxRefreshRate               uint32
	MinRefreshRate               uint32
	Monochrome                   bool
	Name                         string
	NumberOfColorPlanes          uint16
	NumberOfVideoPages           uint32
	PNPDeviceID                  string
	PowerManagementCapabilities  []uint16
	PowerManagementSupported     bool
	ProtocolSupported            uint16
	ReservedSystemPaletteEntries uint32
	SpecificationVersion         uint32
	Status                       string
	StatusInfo                   uint16
	SystemCreationClassName      string
	SystemName                   string
	SystemPaletteEntries         uint32
	TimeOfLastReset              string
	VideoArchitecture            uint16
	VideoMemoryType              uint16
	VideoMode                    uint16
	VideoModeDescription         string
	VideoProcessor               string
}

/*
The Win32_NetworkAdapter class is deprecated. Use the MSFT_NetAdapter class instead. The Win32_NetworkAdapterWMI class represents a network adapter of a computer running a Windows operating system.

Win32_NetworkAdapter only supplies IPv4 data. For more information, see IPv6 and IPv4 Support in WMI.

See: https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/win32-networkadapter
*/
type Win32_NetworkAdapter struct {
	AdapterType                 string
	AdapterTypeID               uint16
	AutoSense                   bool
	Availability                uint16
	Caption                     string
	ConfigManagerErrorCode      uint32
	ConfigManagerUserConfig     bool
	CreationClassName           string
	Description                 string
	DeviceID                    string
	ErrorCleared                bool
	ErrorDescription            string
	GUID                        string
	Index                       uint32
	InstallDate                 string
	Installed                   bool
	InterfaceIndex              uint32
	LastErrorCode               uint32
	MACAddress                  string
	Manufacturer                string
	MaxNumberControlled         uint32
	MaxSpeed                    uint64
	Name                        string
	NetConnectionID             string
	NetConnectionStatus         uint16
	NetEnabled                  bool
	NetworkAddresses            []string
	PermanentAddress            string
	PhysicalAdapter             bool
	PNPDeviceID                 string
	PowerManagementCapabilities []uint16
	PowerManagementSupported    bool
	ProductName                 string
	ServiceName                 string
	Speed                       uint64
	Status                      string
	StatusInfo                  uint16
	SystemCreationClassName     string
	SystemName                  string
	TimeOfLastReset             string
}

/*
The Win32_DiskDrive WMI class represents a physical disk drive as seen by a computer running the Windows operating system.

See: https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/win32-diskdrive
*/
type Win32_DiskDrive struct {
	Availability                uint16
	BytesPerSector              uint32
	Capabilities                []int32 // []uint16 causes panic (call of reflect.Value on int32 Value)
	CapabilityDescriptions      []string
	Caption                     string
	CompressionMethod           string
	ConfigManagerErrorCode      uint32
	ConfigManagerUserConfig     bool
	CreationClassName           string
	DefaultBlockSize            uint64
	Description                 string
	DeviceID                    string
	ErrorCleared                bool
	ErrorDescription            string
	ErrorMethodology            string
	FirmwareRevision            string
	Index                       uint32
	InstallDate                 string
	InterfaceType               string
	LastErrorCode               uint32
	Manufacturer                string
	MaxBlockSize                uint64
	MaxMediaSize                uint64
	MediaLoaded                 bool
	MediaType                   string
	MinBlockSize                uint64
	Model                       string
	Name                        string
	NeedsCleaning               bool
	NumberOfMediaSupported      uint32
	Partitions                  uint32
	PNPDeviceID                 string
	PowerManagementCapabilities []uint16
	PowerManagementSupported    bool
	SCSIBus                     uint32
	SCSILogicalUnit             uint16
	SCSIPort                    uint16
	SCSITargetId                uint16
	SectorsPerTrack             uint32
	SerialNumber                string
	Signature                   uint32
	Size                        uint64
	Status                      string
	StatusInfo                  uint16
	SystemCreationClassName     string
	SystemName                  string
	TotalCylinders              uint64
	TotalHeads                  uint32
	TotalSectors                uint64
	TotalTracks                 uint64
	TracksPerCylinder           uint32
}

/*
The Win32_DiskPartition WMI class represents the capabilities and management capacity of a partitioned area of a physical disk on a computer system running Windows.
Example: Disk #0, Partition #1.

See: https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/win32-diskpartition
*/
type Win32_DiskPartition struct {
	AdditionalAvailability      uint16
	Availability                uint16
	PowerManagementCapabilities []uint16
	IdentifyingDescriptions     [1]string
	MaxQuiesceTime              uint64
	OtherIdentifyingInfo        uint64
	StatusInfo                  uint16
	PowerOnHours                uint64
	TotalPowerOnHours           uint64
	Access                      uint16
	BlockSize                   uint64
	Bootable                    bool
	BootPartition               bool
	Caption                     string
	ConfigManagerErrorCode      uint32
	ConfigManagerUserConfig     bool
	CreationClassName           string
	Description                 string
	DeviceID                    string
	DiskIndex                   uint32
	ErrorCleared                bool
	ErrorDescription            string
	ErrorMethodology            string
	HiddenSectors               uint32
	Index                       uint32
	InstallDate                 string
	LastErrorCode               uint32
	Name                        string
	PNPDeviceID                 string
	PowerManagementSupported    bool
	PrimaryPartition            bool
	Purpose                     string
	RewritePartition            bool
	Size                        uint64
	StartingOffset              uint64
	Status                      string
	SystemCreationClassName     string
	SystemName                  string
	Type                        string
}

/*
The Win32_UserAccount WMI class contains information about a user account on a computer system running Windows.

See: https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/win32-useraccount
*/
type Win32_UserAccount struct {
	AccountType        uint32
	Caption            string
	Description        string
	Disabled           bool
	Domain             string
	FullName           string
	InstallDate        string
	LocalAccount       bool
	Lockout            bool
	Name               string
	PasswordChangeable bool
	PasswordExpires    bool
	PasswordRequired   bool
	SID                string
	SIDType            uint8
	Status             string
}
