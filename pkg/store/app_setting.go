package store

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

type AppSettingStore struct {
	Store   Store
	setting AppSetting
}

func (s *AppSettingStore) All() (AppSetting, error) {
	if s.Store.Modified() {
		s.setting = AppSetting{}
		s.Store.Read(&s.setting)
	} else if !s.Store.Exist() {
		s.setting = AppSetting{SuccessAction: Nothing, SuccessActionDelay: 5, Language: "en"}
		s.Store.Write(s.setting)
	}
	return s.setting, nil
}

func (s *AppSettingStore) Update(v AppSetting) {
	s.setting = v
	s.Store.Write(v)
}
