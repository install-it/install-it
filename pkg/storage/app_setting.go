package storage

import (
	"encoding/json"
	"os"
)

type AppSettingManager struct {
	Path    string
	setting AppSetting
	fstat   os.FileInfo
}

func (s *AppSettingManager) Read() (AppSetting, error) {
	if s.fstat == nil || s.modified() {
		var setting AppSetting

		if _, err := os.Stat(s.Path); err != nil {
			s.Update(AppSetting{SuccessAction: Nothing, SuccessActionDelay: 5, Language: "en"})
		}

		bytes, err := os.ReadFile(s.Path)
		if err != nil {
			return AppSetting{}, err
		}

		if err := json.Unmarshal(bytes, &setting); err != nil {
			return AppSetting{}, err
		}

		s.setting = setting
	}

	return s.setting, nil
}

func (s *AppSettingManager) Update(setting AppSetting) error {
	s.setting = setting

	bytes, err := json.Marshal(s.setting)
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.Path, bytes, os.ModePerm); err == nil {
		s.fstat, _ = os.Stat(s.Path)
		return nil
	} else {
		return err
	}

}

func (s AppSettingManager) modified() bool {
	if s.fstat == nil {
		return false
	}

	if stat, err := os.Stat(s.Path); err != nil {
		return false
	} else {
		return stat.ModTime().After(s.fstat.ModTime())
	}
}

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
