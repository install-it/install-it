export namespace execute {
	
	export class CommandResult {
	    lapse: number;
	    exitCode: number;
	    stdout: string;
	    stderr: string;
	    error: string;
	    aborted: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CommandResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lapse = source["lapse"];
	        this.exitCode = source["exitCode"];
	        this.stdout = source["stdout"];
	        this.stderr = source["stderr"];
	        this.error = source["error"];
	        this.aborted = source["aborted"];
	    }
	}

}

export namespace porter {
	
	export class Progress {
	    name: string;
	    status: status.Status;
	    total: number;
	    current: number;
	    // Go type: time
	    startAt: any;
	    error: any;
	
	    static createFrom(source: any = {}) {
	        return new Progress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.status = source["status"];
	        this.total = source["total"];
	        this.current = source["current"];
	        this.startAt = this.convertValues(source["startAt"], null);
	        this.error = source["error"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Progresses {
	    tasks: Progress[];
	    messages: string[];
	    status: status.Status;
	
	    static createFrom(source: any = {}) {
	        return new Progresses(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tasks = this.convertValues(source["tasks"], Progress);
	        this.messages = source["messages"];
	        this.status = source["status"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace status {
	
	export enum Status {
	    PENDING = "pending",
	    RUNNING = "running",
	    COMPLETED = "completed",
	    FAILED = "failed",
	    ABORTING = "aborting",
	    ABORTED = "aborted",
	    SKIPED = "skiped",
	    SPEEDED = "speeded",
	    ERRORED = "errored",
	}

}

export namespace storage {
	
	export enum DriverType {
	    NETWORK = "network",
	    DISPLAY = "display",
	    MISCELLANEOUS = "miscellaneous",
	}
	export enum SuccessAction {
	    NOTHING = "nothing",
	    REBOOT = "reboot",
	    SHUTDOWN = "shutdown",
	    FIRMWARE = "firmware",
	}
	export class AppSetting {
	    create_partition: boolean;
	    set_password: boolean;
	    password: string;
	    parallel_install: boolean;
	    success_action: SuccessAction;
	    success_action_delay: number;
	    filter_miniport_nic: boolean;
	    filter_microsoft_nic: boolean;
	    language: string;
	    driver_download_url: string;
	    auto_check_update: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AppSetting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.create_partition = source["create_partition"];
	        this.set_password = source["set_password"];
	        this.password = source["password"];
	        this.parallel_install = source["parallel_install"];
	        this.success_action = source["success_action"];
	        this.success_action_delay = source["success_action_delay"];
	        this.filter_miniport_nic = source["filter_miniport_nic"];
	        this.filter_microsoft_nic = source["filter_microsoft_nic"];
	        this.language = source["language"];
	        this.driver_download_url = source["driver_download_url"];
	        this.auto_check_update = source["auto_check_update"];
	    }
	}
	export class Driver {
	    id: string;
	    name: string;
	    type: DriverType;
	    path: string;
	    flags: string[];
	    minExeTime: number;
	    allowRtCodes: number[];
	    incompatibles: string[];
	
	    static createFrom(source: any = {}) {
	        return new Driver(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.path = source["path"];
	        this.flags = source["flags"];
	        this.minExeTime = source["minExeTime"];
	        this.allowRtCodes = source["allowRtCodes"];
	        this.incompatibles = source["incompatibles"];
	    }
	}
	export class DriverGroup {
	    id: string;
	    name: string;
	    type: DriverType;
	    drivers: Driver[];
	
	    static createFrom(source: any = {}) {
	        return new DriverGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.drivers = this.convertValues(source["drivers"], Driver);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace sysinfo {
	
	export class Win32_BaseBoard {
	    Caption: string;
	    ConfigOptions: string[];
	    CreationClassName: string;
	    Depth: number;
	    Description: string;
	    Height: number;
	    HostingBoard: boolean;
	    HotSwappable: boolean;
	    InstallDate: string;
	    Manufacturer: string;
	    Model: string;
	    Name: string;
	    OtherIdentifyingInfo: string;
	    PartNumber: string;
	    PoweredOn: boolean;
	    Product: string;
	    Removable: boolean;
	    Replaceable: boolean;
	    RequirementsDescription: string;
	    RequiresDaughterBoard: boolean;
	    SerialNumber: string;
	    SKU: string;
	    SlotLayout: string;
	    SpecialRequirements: boolean;
	    Status: string;
	    Tag: string;
	    Version: string;
	    Weight: number;
	    Width: number;
	
	    static createFrom(source: any = {}) {
	        return new Win32_BaseBoard(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Caption = source["Caption"];
	        this.ConfigOptions = source["ConfigOptions"];
	        this.CreationClassName = source["CreationClassName"];
	        this.Depth = source["Depth"];
	        this.Description = source["Description"];
	        this.Height = source["Height"];
	        this.HostingBoard = source["HostingBoard"];
	        this.HotSwappable = source["HotSwappable"];
	        this.InstallDate = source["InstallDate"];
	        this.Manufacturer = source["Manufacturer"];
	        this.Model = source["Model"];
	        this.Name = source["Name"];
	        this.OtherIdentifyingInfo = source["OtherIdentifyingInfo"];
	        this.PartNumber = source["PartNumber"];
	        this.PoweredOn = source["PoweredOn"];
	        this.Product = source["Product"];
	        this.Removable = source["Removable"];
	        this.Replaceable = source["Replaceable"];
	        this.RequirementsDescription = source["RequirementsDescription"];
	        this.RequiresDaughterBoard = source["RequiresDaughterBoard"];
	        this.SerialNumber = source["SerialNumber"];
	        this.SKU = source["SKU"];
	        this.SlotLayout = source["SlotLayout"];
	        this.SpecialRequirements = source["SpecialRequirements"];
	        this.Status = source["Status"];
	        this.Tag = source["Tag"];
	        this.Version = source["Version"];
	        this.Weight = source["Weight"];
	        this.Width = source["Width"];
	    }
	}
	export class Win32_DiskDrive {
	    Availability: number;
	    BytesPerSector: number;
	    Capabilities: number[];
	    CapabilityDescriptions: string[];
	    Caption: string;
	    CompressionMethod: string;
	    ConfigManagerErrorCode: number;
	    ConfigManagerUserConfig: boolean;
	    CreationClassName: string;
	    DefaultBlockSize: number;
	    Description: string;
	    DeviceID: string;
	    ErrorCleared: boolean;
	    ErrorDescription: string;
	    ErrorMethodology: string;
	    FirmwareRevision: string;
	    Index: number;
	    InstallDate: string;
	    InterfaceType: string;
	    LastErrorCode: number;
	    Manufacturer: string;
	    MaxBlockSize: number;
	    MaxMediaSize: number;
	    MediaLoaded: boolean;
	    MediaType: string;
	    MinBlockSize: number;
	    Model: string;
	    Name: string;
	    NeedsCleaning: boolean;
	    NumberOfMediaSupported: number;
	    Partitions: number;
	    PNPDeviceID: string;
	    PowerManagementCapabilities: number[];
	    PowerManagementSupported: boolean;
	    SCSIBus: number;
	    SCSILogicalUnit: number;
	    SCSIPort: number;
	    SCSITargetId: number;
	    SectorsPerTrack: number;
	    SerialNumber: string;
	    Signature: number;
	    Size: number;
	    Status: string;
	    StatusInfo: number;
	    SystemCreationClassName: string;
	    SystemName: string;
	    TotalCylinders: number;
	    TotalHeads: number;
	    TotalSectors: number;
	    TotalTracks: number;
	    TracksPerCylinder: number;
	
	    static createFrom(source: any = {}) {
	        return new Win32_DiskDrive(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Availability = source["Availability"];
	        this.BytesPerSector = source["BytesPerSector"];
	        this.Capabilities = source["Capabilities"];
	        this.CapabilityDescriptions = source["CapabilityDescriptions"];
	        this.Caption = source["Caption"];
	        this.CompressionMethod = source["CompressionMethod"];
	        this.ConfigManagerErrorCode = source["ConfigManagerErrorCode"];
	        this.ConfigManagerUserConfig = source["ConfigManagerUserConfig"];
	        this.CreationClassName = source["CreationClassName"];
	        this.DefaultBlockSize = source["DefaultBlockSize"];
	        this.Description = source["Description"];
	        this.DeviceID = source["DeviceID"];
	        this.ErrorCleared = source["ErrorCleared"];
	        this.ErrorDescription = source["ErrorDescription"];
	        this.ErrorMethodology = source["ErrorMethodology"];
	        this.FirmwareRevision = source["FirmwareRevision"];
	        this.Index = source["Index"];
	        this.InstallDate = source["InstallDate"];
	        this.InterfaceType = source["InterfaceType"];
	        this.LastErrorCode = source["LastErrorCode"];
	        this.Manufacturer = source["Manufacturer"];
	        this.MaxBlockSize = source["MaxBlockSize"];
	        this.MaxMediaSize = source["MaxMediaSize"];
	        this.MediaLoaded = source["MediaLoaded"];
	        this.MediaType = source["MediaType"];
	        this.MinBlockSize = source["MinBlockSize"];
	        this.Model = source["Model"];
	        this.Name = source["Name"];
	        this.NeedsCleaning = source["NeedsCleaning"];
	        this.NumberOfMediaSupported = source["NumberOfMediaSupported"];
	        this.Partitions = source["Partitions"];
	        this.PNPDeviceID = source["PNPDeviceID"];
	        this.PowerManagementCapabilities = source["PowerManagementCapabilities"];
	        this.PowerManagementSupported = source["PowerManagementSupported"];
	        this.SCSIBus = source["SCSIBus"];
	        this.SCSILogicalUnit = source["SCSILogicalUnit"];
	        this.SCSIPort = source["SCSIPort"];
	        this.SCSITargetId = source["SCSITargetId"];
	        this.SectorsPerTrack = source["SectorsPerTrack"];
	        this.SerialNumber = source["SerialNumber"];
	        this.Signature = source["Signature"];
	        this.Size = source["Size"];
	        this.Status = source["Status"];
	        this.StatusInfo = source["StatusInfo"];
	        this.SystemCreationClassName = source["SystemCreationClassName"];
	        this.SystemName = source["SystemName"];
	        this.TotalCylinders = source["TotalCylinders"];
	        this.TotalHeads = source["TotalHeads"];
	        this.TotalSectors = source["TotalSectors"];
	        this.TotalTracks = source["TotalTracks"];
	        this.TracksPerCylinder = source["TracksPerCylinder"];
	    }
	}
	export class Win32_DiskPartition {
	    AdditionalAvailability: number;
	    Availability: number;
	    PowerManagementCapabilities: number[];
	    IdentifyingDescriptions: string[];
	    MaxQuiesceTime: number;
	    OtherIdentifyingInfo: number;
	    StatusInfo: number;
	    PowerOnHours: number;
	    TotalPowerOnHours: number;
	    Access: number;
	    BlockSize: number;
	    Bootable: boolean;
	    BootPartition: boolean;
	    Caption: string;
	    ConfigManagerErrorCode: number;
	    ConfigManagerUserConfig: boolean;
	    CreationClassName: string;
	    Description: string;
	    DeviceID: string;
	    DiskIndex: number;
	    ErrorCleared: boolean;
	    ErrorDescription: string;
	    ErrorMethodology: string;
	    HiddenSectors: number;
	    Index: number;
	    InstallDate: string;
	    LastErrorCode: number;
	    Name: string;
	    PNPDeviceID: string;
	    PowerManagementSupported: boolean;
	    PrimaryPartition: boolean;
	    Purpose: string;
	    RewritePartition: boolean;
	    Size: number;
	    StartingOffset: number;
	    Status: string;
	    SystemCreationClassName: string;
	    SystemName: string;
	    Type: string;
	
	    static createFrom(source: any = {}) {
	        return new Win32_DiskPartition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AdditionalAvailability = source["AdditionalAvailability"];
	        this.Availability = source["Availability"];
	        this.PowerManagementCapabilities = source["PowerManagementCapabilities"];
	        this.IdentifyingDescriptions = source["IdentifyingDescriptions"];
	        this.MaxQuiesceTime = source["MaxQuiesceTime"];
	        this.OtherIdentifyingInfo = source["OtherIdentifyingInfo"];
	        this.StatusInfo = source["StatusInfo"];
	        this.PowerOnHours = source["PowerOnHours"];
	        this.TotalPowerOnHours = source["TotalPowerOnHours"];
	        this.Access = source["Access"];
	        this.BlockSize = source["BlockSize"];
	        this.Bootable = source["Bootable"];
	        this.BootPartition = source["BootPartition"];
	        this.Caption = source["Caption"];
	        this.ConfigManagerErrorCode = source["ConfigManagerErrorCode"];
	        this.ConfigManagerUserConfig = source["ConfigManagerUserConfig"];
	        this.CreationClassName = source["CreationClassName"];
	        this.Description = source["Description"];
	        this.DeviceID = source["DeviceID"];
	        this.DiskIndex = source["DiskIndex"];
	        this.ErrorCleared = source["ErrorCleared"];
	        this.ErrorDescription = source["ErrorDescription"];
	        this.ErrorMethodology = source["ErrorMethodology"];
	        this.HiddenSectors = source["HiddenSectors"];
	        this.Index = source["Index"];
	        this.InstallDate = source["InstallDate"];
	        this.LastErrorCode = source["LastErrorCode"];
	        this.Name = source["Name"];
	        this.PNPDeviceID = source["PNPDeviceID"];
	        this.PowerManagementSupported = source["PowerManagementSupported"];
	        this.PrimaryPartition = source["PrimaryPartition"];
	        this.Purpose = source["Purpose"];
	        this.RewritePartition = source["RewritePartition"];
	        this.Size = source["Size"];
	        this.StartingOffset = source["StartingOffset"];
	        this.Status = source["Status"];
	        this.SystemCreationClassName = source["SystemCreationClassName"];
	        this.SystemName = source["SystemName"];
	        this.Type = source["Type"];
	    }
	}
	export class Win32_NetworkAdapter {
	    AdapterType: string;
	    AdapterTypeID: number;
	    AutoSense: boolean;
	    Availability: number;
	    Caption: string;
	    ConfigManagerErrorCode: number;
	    ConfigManagerUserConfig: boolean;
	    CreationClassName: string;
	    Description: string;
	    DeviceID: string;
	    ErrorCleared: boolean;
	    ErrorDescription: string;
	    GUID: string;
	    Index: number;
	    InstallDate: string;
	    Installed: boolean;
	    InterfaceIndex: number;
	    LastErrorCode: number;
	    MACAddress: string;
	    Manufacturer: string;
	    MaxNumberControlled: number;
	    MaxSpeed: number;
	    Name: string;
	    NetConnectionID: string;
	    NetConnectionStatus: number;
	    NetEnabled: boolean;
	    NetworkAddresses: string[];
	    PermanentAddress: string;
	    PhysicalAdapter: boolean;
	    PNPDeviceID: string;
	    PowerManagementCapabilities: number[];
	    PowerManagementSupported: boolean;
	    ProductName: string;
	    ServiceName: string;
	    Speed: number;
	    Status: string;
	    StatusInfo: number;
	    SystemCreationClassName: string;
	    SystemName: string;
	    TimeOfLastReset: string;
	
	    static createFrom(source: any = {}) {
	        return new Win32_NetworkAdapter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AdapterType = source["AdapterType"];
	        this.AdapterTypeID = source["AdapterTypeID"];
	        this.AutoSense = source["AutoSense"];
	        this.Availability = source["Availability"];
	        this.Caption = source["Caption"];
	        this.ConfigManagerErrorCode = source["ConfigManagerErrorCode"];
	        this.ConfigManagerUserConfig = source["ConfigManagerUserConfig"];
	        this.CreationClassName = source["CreationClassName"];
	        this.Description = source["Description"];
	        this.DeviceID = source["DeviceID"];
	        this.ErrorCleared = source["ErrorCleared"];
	        this.ErrorDescription = source["ErrorDescription"];
	        this.GUID = source["GUID"];
	        this.Index = source["Index"];
	        this.InstallDate = source["InstallDate"];
	        this.Installed = source["Installed"];
	        this.InterfaceIndex = source["InterfaceIndex"];
	        this.LastErrorCode = source["LastErrorCode"];
	        this.MACAddress = source["MACAddress"];
	        this.Manufacturer = source["Manufacturer"];
	        this.MaxNumberControlled = source["MaxNumberControlled"];
	        this.MaxSpeed = source["MaxSpeed"];
	        this.Name = source["Name"];
	        this.NetConnectionID = source["NetConnectionID"];
	        this.NetConnectionStatus = source["NetConnectionStatus"];
	        this.NetEnabled = source["NetEnabled"];
	        this.NetworkAddresses = source["NetworkAddresses"];
	        this.PermanentAddress = source["PermanentAddress"];
	        this.PhysicalAdapter = source["PhysicalAdapter"];
	        this.PNPDeviceID = source["PNPDeviceID"];
	        this.PowerManagementCapabilities = source["PowerManagementCapabilities"];
	        this.PowerManagementSupported = source["PowerManagementSupported"];
	        this.ProductName = source["ProductName"];
	        this.ServiceName = source["ServiceName"];
	        this.Speed = source["Speed"];
	        this.Status = source["Status"];
	        this.StatusInfo = source["StatusInfo"];
	        this.SystemCreationClassName = source["SystemCreationClassName"];
	        this.SystemName = source["SystemName"];
	        this.TimeOfLastReset = source["TimeOfLastReset"];
	    }
	}
	export class Win32_PhysicalMemory {
	    Attributes: number;
	    BankLabel: string;
	    Capacity: number;
	    Caption: string;
	    ConfiguredClockSpeed: number;
	    ConfiguredVoltage: number;
	    CreationClassName: string;
	    DataWidth: number;
	    Description: string;
	    DeviceLocator: string;
	    FormFactor: number;
	    HotSwappable: boolean;
	    InstallDate: string;
	    InterleaveDataDepth: number;
	    InterleavePosition: number;
	    Manufacturer: string;
	    MaxVoltage: number;
	    MemoryType: number;
	    MinVoltage: number;
	    Model: string;
	    Name: string;
	    OtherIdentifyingInfo: string;
	    PartNumber: string;
	    PositionInRow: number;
	    PoweredOn: boolean;
	    Removable: boolean;
	    Replaceable: boolean;
	    SerialNumber: string;
	    SKU: string;
	    SMBIOSMemoryType: number;
	    Speed: number;
	    Status: string;
	    Tag: string;
	    TotalWidth: number;
	    TypeDetail: number;
	    Version: string;
	
	    static createFrom(source: any = {}) {
	        return new Win32_PhysicalMemory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Attributes = source["Attributes"];
	        this.BankLabel = source["BankLabel"];
	        this.Capacity = source["Capacity"];
	        this.Caption = source["Caption"];
	        this.ConfiguredClockSpeed = source["ConfiguredClockSpeed"];
	        this.ConfiguredVoltage = source["ConfiguredVoltage"];
	        this.CreationClassName = source["CreationClassName"];
	        this.DataWidth = source["DataWidth"];
	        this.Description = source["Description"];
	        this.DeviceLocator = source["DeviceLocator"];
	        this.FormFactor = source["FormFactor"];
	        this.HotSwappable = source["HotSwappable"];
	        this.InstallDate = source["InstallDate"];
	        this.InterleaveDataDepth = source["InterleaveDataDepth"];
	        this.InterleavePosition = source["InterleavePosition"];
	        this.Manufacturer = source["Manufacturer"];
	        this.MaxVoltage = source["MaxVoltage"];
	        this.MemoryType = source["MemoryType"];
	        this.MinVoltage = source["MinVoltage"];
	        this.Model = source["Model"];
	        this.Name = source["Name"];
	        this.OtherIdentifyingInfo = source["OtherIdentifyingInfo"];
	        this.PartNumber = source["PartNumber"];
	        this.PositionInRow = source["PositionInRow"];
	        this.PoweredOn = source["PoweredOn"];
	        this.Removable = source["Removable"];
	        this.Replaceable = source["Replaceable"];
	        this.SerialNumber = source["SerialNumber"];
	        this.SKU = source["SKU"];
	        this.SMBIOSMemoryType = source["SMBIOSMemoryType"];
	        this.Speed = source["Speed"];
	        this.Status = source["Status"];
	        this.Tag = source["Tag"];
	        this.TotalWidth = source["TotalWidth"];
	        this.TypeDetail = source["TypeDetail"];
	        this.Version = source["Version"];
	    }
	}
	export class Win32_Processor {
	    AddressWidth: number;
	    Architecture: number;
	    AssetTag: string;
	    Availability: number;
	    Caption: string;
	    Characteristics: number;
	    ConfigManagerErrorCode: number;
	    ConfigManagerUserConfig: boolean;
	    CpuStatus: number;
	    CreationClassName: string;
	    CurrentClockSpeed: number;
	    CurrentVoltage: number;
	    DataWidth: number;
	    Description: string;
	    DeviceID: string;
	    ErrorCleared: boolean;
	    ErrorDescription: string;
	    ExtClock: number;
	    Family: number;
	    InstallDate: string;
	    L2CacheSize: number;
	    L2CacheSpeed: number;
	    L3CacheSize: number;
	    L3CacheSpeed: number;
	    LastErrorCode: number;
	    Level: number;
	    LoadPercentage: number;
	    Manufacturer: string;
	    MaxClockSpeed: number;
	    Name: string;
	    NumberOfCores: number;
	    NumberOfEnabledCore: number;
	    NumberOfLogicalProcessors: number;
	    OtherFamilyDescription: string;
	    PartNumber: string;
	    PNPDeviceID: string;
	    PowerManagementCapabilities: number[];
	    PowerManagementSupported: boolean;
	    ProcessorId: string;
	    ProcessorType: number;
	    Revision: number;
	    Role: string;
	    SecondLevelAddressTranslationExtensions: boolean;
	    SerialNumber: string;
	    SocketDesignation: string;
	    Status: string;
	    StatusInfo: number;
	    Stepping: string;
	    SystemCreationClassName: string;
	    SystemName: string;
	    ThreadCount: number;
	    UniqueId: string;
	    UpgradeMethod: number;
	    Version: string;
	    VirtualizationFirmwareEnabled: boolean;
	    VMMonitorModeExtensions: boolean;
	    VoltageCaps: number;
	
	    static createFrom(source: any = {}) {
	        return new Win32_Processor(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AddressWidth = source["AddressWidth"];
	        this.Architecture = source["Architecture"];
	        this.AssetTag = source["AssetTag"];
	        this.Availability = source["Availability"];
	        this.Caption = source["Caption"];
	        this.Characteristics = source["Characteristics"];
	        this.ConfigManagerErrorCode = source["ConfigManagerErrorCode"];
	        this.ConfigManagerUserConfig = source["ConfigManagerUserConfig"];
	        this.CpuStatus = source["CpuStatus"];
	        this.CreationClassName = source["CreationClassName"];
	        this.CurrentClockSpeed = source["CurrentClockSpeed"];
	        this.CurrentVoltage = source["CurrentVoltage"];
	        this.DataWidth = source["DataWidth"];
	        this.Description = source["Description"];
	        this.DeviceID = source["DeviceID"];
	        this.ErrorCleared = source["ErrorCleared"];
	        this.ErrorDescription = source["ErrorDescription"];
	        this.ExtClock = source["ExtClock"];
	        this.Family = source["Family"];
	        this.InstallDate = source["InstallDate"];
	        this.L2CacheSize = source["L2CacheSize"];
	        this.L2CacheSpeed = source["L2CacheSpeed"];
	        this.L3CacheSize = source["L3CacheSize"];
	        this.L3CacheSpeed = source["L3CacheSpeed"];
	        this.LastErrorCode = source["LastErrorCode"];
	        this.Level = source["Level"];
	        this.LoadPercentage = source["LoadPercentage"];
	        this.Manufacturer = source["Manufacturer"];
	        this.MaxClockSpeed = source["MaxClockSpeed"];
	        this.Name = source["Name"];
	        this.NumberOfCores = source["NumberOfCores"];
	        this.NumberOfEnabledCore = source["NumberOfEnabledCore"];
	        this.NumberOfLogicalProcessors = source["NumberOfLogicalProcessors"];
	        this.OtherFamilyDescription = source["OtherFamilyDescription"];
	        this.PartNumber = source["PartNumber"];
	        this.PNPDeviceID = source["PNPDeviceID"];
	        this.PowerManagementCapabilities = source["PowerManagementCapabilities"];
	        this.PowerManagementSupported = source["PowerManagementSupported"];
	        this.ProcessorId = source["ProcessorId"];
	        this.ProcessorType = source["ProcessorType"];
	        this.Revision = source["Revision"];
	        this.Role = source["Role"];
	        this.SecondLevelAddressTranslationExtensions = source["SecondLevelAddressTranslationExtensions"];
	        this.SerialNumber = source["SerialNumber"];
	        this.SocketDesignation = source["SocketDesignation"];
	        this.Status = source["Status"];
	        this.StatusInfo = source["StatusInfo"];
	        this.Stepping = source["Stepping"];
	        this.SystemCreationClassName = source["SystemCreationClassName"];
	        this.SystemName = source["SystemName"];
	        this.ThreadCount = source["ThreadCount"];
	        this.UniqueId = source["UniqueId"];
	        this.UpgradeMethod = source["UpgradeMethod"];
	        this.Version = source["Version"];
	        this.VirtualizationFirmwareEnabled = source["VirtualizationFirmwareEnabled"];
	        this.VMMonitorModeExtensions = source["VMMonitorModeExtensions"];
	        this.VoltageCaps = source["VoltageCaps"];
	    }
	}
	export class Win32_UserAccount {
	    AccountType: number;
	    Caption: string;
	    Description: string;
	    Disabled: boolean;
	    Domain: string;
	    FullName: string;
	    InstallDate: string;
	    LocalAccount: boolean;
	    Lockout: boolean;
	    Name: string;
	    PasswordChangeable: boolean;
	    PasswordExpires: boolean;
	    PasswordRequired: boolean;
	    SID: string;
	    SIDType: number;
	    Status: string;
	
	    static createFrom(source: any = {}) {
	        return new Win32_UserAccount(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AccountType = source["AccountType"];
	        this.Caption = source["Caption"];
	        this.Description = source["Description"];
	        this.Disabled = source["Disabled"];
	        this.Domain = source["Domain"];
	        this.FullName = source["FullName"];
	        this.InstallDate = source["InstallDate"];
	        this.LocalAccount = source["LocalAccount"];
	        this.Lockout = source["Lockout"];
	        this.Name = source["Name"];
	        this.PasswordChangeable = source["PasswordChangeable"];
	        this.PasswordExpires = source["PasswordExpires"];
	        this.PasswordRequired = source["PasswordRequired"];
	        this.SID = source["SID"];
	        this.SIDType = source["SIDType"];
	        this.Status = source["Status"];
	    }
	}
	export class Win32_VideoController {
	    AcceleratorCapabilities: number[];
	    AdapterCompatibility: string;
	    AdapterDACType: string;
	    AdapterRAM: number;
	    Availability: number;
	    CapabilityDescriptions: string[];
	    Caption: string;
	    ColorTableEntries: number;
	    ConfigManagerErrorCode: number;
	    ConfigManagerUserConfig: boolean;
	    CreationClassName: string;
	    CurrentBitsPerPixel: number;
	    CurrentHorizontalResolution: number;
	    CurrentNumberOfColors: number;
	    CurrentNumberOfColumns: number;
	    CurrentNumberOfRows: number;
	    CurrentRefreshRate: number;
	    CurrentScanMode: number;
	    CurrentVerticalResolution: number;
	    Description: string;
	    DeviceID: string;
	    DeviceSpecificPens: number;
	    DitherType: number;
	    DriverDate: string;
	    DriverVersion: string;
	    ErrorCleared: boolean;
	    ErrorDescription: string;
	    ICMIntent: number;
	    ICMMethod: number;
	    InfFilename: string;
	    InfSection: string;
	    InstallDate: string;
	    InstalledDisplayDrivers: string;
	    LastErrorCode: number;
	    MaxMemorySupported: number;
	    MaxNumberControlled: number;
	    MaxRefreshRate: number;
	    MinRefreshRate: number;
	    Monochrome: boolean;
	    Name: string;
	    NumberOfColorPlanes: number;
	    NumberOfVideoPages: number;
	    PNPDeviceID: string;
	    PowerManagementCapabilities: number[];
	    PowerManagementSupported: boolean;
	    ProtocolSupported: number;
	    ReservedSystemPaletteEntries: number;
	    SpecificationVersion: number;
	    Status: string;
	    StatusInfo: number;
	    SystemCreationClassName: string;
	    SystemName: string;
	    SystemPaletteEntries: number;
	    TimeOfLastReset: string;
	    VideoArchitecture: number;
	    VideoMemoryType: number;
	    VideoMode: number;
	    VideoModeDescription: string;
	    VideoProcessor: string;
	
	    static createFrom(source: any = {}) {
	        return new Win32_VideoController(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AcceleratorCapabilities = source["AcceleratorCapabilities"];
	        this.AdapterCompatibility = source["AdapterCompatibility"];
	        this.AdapterDACType = source["AdapterDACType"];
	        this.AdapterRAM = source["AdapterRAM"];
	        this.Availability = source["Availability"];
	        this.CapabilityDescriptions = source["CapabilityDescriptions"];
	        this.Caption = source["Caption"];
	        this.ColorTableEntries = source["ColorTableEntries"];
	        this.ConfigManagerErrorCode = source["ConfigManagerErrorCode"];
	        this.ConfigManagerUserConfig = source["ConfigManagerUserConfig"];
	        this.CreationClassName = source["CreationClassName"];
	        this.CurrentBitsPerPixel = source["CurrentBitsPerPixel"];
	        this.CurrentHorizontalResolution = source["CurrentHorizontalResolution"];
	        this.CurrentNumberOfColors = source["CurrentNumberOfColors"];
	        this.CurrentNumberOfColumns = source["CurrentNumberOfColumns"];
	        this.CurrentNumberOfRows = source["CurrentNumberOfRows"];
	        this.CurrentRefreshRate = source["CurrentRefreshRate"];
	        this.CurrentScanMode = source["CurrentScanMode"];
	        this.CurrentVerticalResolution = source["CurrentVerticalResolution"];
	        this.Description = source["Description"];
	        this.DeviceID = source["DeviceID"];
	        this.DeviceSpecificPens = source["DeviceSpecificPens"];
	        this.DitherType = source["DitherType"];
	        this.DriverDate = source["DriverDate"];
	        this.DriverVersion = source["DriverVersion"];
	        this.ErrorCleared = source["ErrorCleared"];
	        this.ErrorDescription = source["ErrorDescription"];
	        this.ICMIntent = source["ICMIntent"];
	        this.ICMMethod = source["ICMMethod"];
	        this.InfFilename = source["InfFilename"];
	        this.InfSection = source["InfSection"];
	        this.InstallDate = source["InstallDate"];
	        this.InstalledDisplayDrivers = source["InstalledDisplayDrivers"];
	        this.LastErrorCode = source["LastErrorCode"];
	        this.MaxMemorySupported = source["MaxMemorySupported"];
	        this.MaxNumberControlled = source["MaxNumberControlled"];
	        this.MaxRefreshRate = source["MaxRefreshRate"];
	        this.MinRefreshRate = source["MinRefreshRate"];
	        this.Monochrome = source["Monochrome"];
	        this.Name = source["Name"];
	        this.NumberOfColorPlanes = source["NumberOfColorPlanes"];
	        this.NumberOfVideoPages = source["NumberOfVideoPages"];
	        this.PNPDeviceID = source["PNPDeviceID"];
	        this.PowerManagementCapabilities = source["PowerManagementCapabilities"];
	        this.PowerManagementSupported = source["PowerManagementSupported"];
	        this.ProtocolSupported = source["ProtocolSupported"];
	        this.ReservedSystemPaletteEntries = source["ReservedSystemPaletteEntries"];
	        this.SpecificationVersion = source["SpecificationVersion"];
	        this.Status = source["Status"];
	        this.StatusInfo = source["StatusInfo"];
	        this.SystemCreationClassName = source["SystemCreationClassName"];
	        this.SystemName = source["SystemName"];
	        this.SystemPaletteEntries = source["SystemPaletteEntries"];
	        this.TimeOfLastReset = source["TimeOfLastReset"];
	        this.VideoArchitecture = source["VideoArchitecture"];
	        this.VideoMemoryType = source["VideoMemoryType"];
	        this.VideoMode = source["VideoMode"];
	        this.VideoModeDescription = source["VideoModeDescription"];
	        this.VideoProcessor = source["VideoProcessor"];
	    }
	}

}

