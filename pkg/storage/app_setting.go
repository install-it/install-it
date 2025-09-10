package storage

type AppSetting struct {
	CreatePartition    bool          `json:"create_partition"`
	SetPassword        bool          `json:"set_password"`
	Password           string        `json:"password"`
	ParallelInstall    bool          `json:"parallel_install"`
	SuccessAction      SuccessAction `json:"success_action"`
	SuccessActionDelay int           `json:"success_action_delay"`
	FilterMiniportNic  bool          `json:"filter_miniport_nic"`
	FilterMicrosoftNic bool          `json:"filter_microsoft_nic"`
	Language           string        `json:"language"`
	DriverDownloadUrl  string        `json:"driver_download_url"`
	AutoCheckUpdate    bool          `json:"auto_check_update"`
}

type SuccessAction string

const (
	Nothing  SuccessAction = "nothing"
	Shutdown SuccessAction = "shutdown"
	Reboot   SuccessAction = "reboot"
	Firmware SuccessAction = "firmware"
)

type AppSettingStorage struct {
	Store   Store
	setting AppSetting
}

func (s *AppSettingStorage) All() (AppSetting, error) {
	if !s.Store.Exist() {
		s.setting = AppSetting{
			AutoCheckUpdate: true, FilterMiniportNic: true, FilterMicrosoftNic: true,
			Language: "en", SuccessAction: Nothing, SuccessActionDelay: 5,
		}
		s.Store.Write(s.setting)
	} else {
		s.Store.Read(&s.setting)
	}
	return s.setting, nil
}

func (s *AppSettingStorage) Update(v AppSetting) (AppSetting, error) {
	s.setting = v
	return s.setting, s.Store.Write(v)
}
