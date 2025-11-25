package storage

import (
	"testing"
)

func TestAppSettingStorage_AllWithoutExistingFile(t *testing.T) {
	storage := &AppSettingStorage{Store: &MemoryStore{}}

	result, err := storage.All()

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Verify default values
	if result.AutoCheckUpdate != true {
		t.Errorf("expected AutoCheckUpdate to be true, got %v", result.AutoCheckUpdate)
	}
	if result.FilterMiniportNic != true {
		t.Errorf("expected FilterMiniportNic to be true, got %v", result.FilterMiniportNic)
	}
	if result.FilterMicrosoftNic != true {
		t.Errorf("expected FilterMicrosoftNic to be true, got %v", result.FilterMicrosoftNic)
	}
	if result.Language != "en" {
		t.Errorf("expected Language to be 'en', got %s", result.Language)
	}
	if result.ParallelInstall != true {
		t.Errorf("expected ParallelInstall to be true, got %v", result.ParallelInstall)
	}
	if result.SuccessAction != Nothing {
		t.Errorf("expected SuccessAction to be 'nothing', got %s", result.SuccessAction)
	}
	if result.SuccessActionDelay != 5 {
		t.Errorf("expected SuccessActionDelay to be 5, got %d", result.SuccessActionDelay)
	}
}

func TestAppSettingStorage_AllWithExistingFile(t *testing.T) {
	existingSetting := AppSetting{
		CreatePartition:    true,
		SetPassword:        true,
		Password:           "test123",
		ParallelInstall:    false,
		SuccessAction:      Reboot,
		SuccessActionDelay: 10,
		FilterMiniportNic:  false,
		FilterMicrosoftNic: false,
		Language:           "zh_Hant_HK",
		DriverDownloadUrl:  "https://example.com",
		AutoCheckUpdate:    false,
		HideNotFound:       true,
	}

	memStore := &MemoryStore{}
	memStore.Write(existingSetting)
	storage := &AppSettingStorage{Store: memStore}

	result, err := storage.All()

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Verify loaded values
	if result.CreatePartition != true {
		t.Errorf("expected CreatePartition to be true, got %v", result.CreatePartition)
	}
	if result.SetPassword != true {
		t.Errorf("expected SetPassword to be true, got %v", result.SetPassword)
	}
	if result.Password != "test123" {
		t.Errorf("expected Password to be 'test123', got %s", result.Password)
	}
	if result.ParallelInstall != false {
		t.Errorf("expected ParallelInstall to be false, got %v", result.ParallelInstall)
	}
	if result.SuccessAction != Reboot {
		t.Errorf("expected SuccessAction to be 'reboot', got %s", result.SuccessAction)
	}
	if result.SuccessActionDelay != 10 {
		t.Errorf("expected SuccessActionDelay to be 10, got %d", result.SuccessActionDelay)
	}
	if result.Language != "zh_Hant_HK" {
		t.Errorf("expected Language to be 'zh_Hant_HK', got %s", result.Language)
	}
	if result.DriverDownloadUrl != "https://example.com" {
		t.Errorf("expected DriverDownloadUrl to be 'https://example.com', got %s", result.DriverDownloadUrl)
	}
	if result.AutoCheckUpdate != false {
		t.Errorf("expected AutoCheckUpdate to be false, got %v", result.AutoCheckUpdate)
	}
	if result.HideNotFound != true {
		t.Errorf("expected HideNotFound to be true, got %v", result.HideNotFound)
	}
}

func TestAppSettingStorage_Update(t *testing.T) {
	storage := &AppSettingStorage{Store: &MemoryStore{}}

	newSetting := AppSetting{
		CreatePartition:    true,
		SetPassword:        true,
		Password:           "newpass",
		ParallelInstall:    false,
		SuccessAction:      Shutdown,
		SuccessActionDelay: 30,
		FilterMiniportNic:  true,
		FilterMicrosoftNic: false,
		Language:           "zh_Hant_HK",
		DriverDownloadUrl:  "https://drivers.example.com",
		AutoCheckUpdate:    false,
		HideNotFound:       true,
	}

	result, err := storage.Update(newSetting)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if result.Password != "newpass" {
		t.Errorf("expected Password to be 'newpass', got %s", result.Password)
	}
	if result.SuccessAction != Shutdown {
		t.Errorf("expected SuccessAction to be 'shutdown', got %s", result.SuccessAction)
	}
	if result.SuccessActionDelay != 30 {
		t.Errorf("expected SuccessActionDelay to be 30, got %d", result.SuccessActionDelay)
	}

}

func TestAppSettingStorage_UpdateMultipleTimes(t *testing.T) {
	storage := &AppSettingStorage{Store: &MemoryStore{}}

	setting1 := AppSetting{
		Language:        "en",
		AutoCheckUpdate: true,
		SuccessAction:   Nothing,
	}
	storage.Update(setting1)

	setting2 := AppSetting{
		Language:        "zh_Hant_HK",
		AutoCheckUpdate: false,
		SuccessAction:   Reboot,
	}
	result, err := storage.Update(setting2)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if result.Language != "zh_Hant_HK" {
		t.Errorf("expected Language to be 'zh_Hant_HK', got %s", result.Language)
	}
	if result.AutoCheckUpdate != false {
		t.Errorf("expected AutoCheckUpdate to be false, got %v", result.AutoCheckUpdate)
	}
	if result.SuccessAction != Reboot {
		t.Errorf("expected SuccessAction to be 'reboot', got %s", result.SuccessAction)
	}
}

func TestAppSettingStorage_UpdateCachesResult(t *testing.T) {
	storage := &AppSettingStorage{Store: &MemoryStore{}}

	setting := AppSetting{
		Language:        "en",
		AutoCheckUpdate: true,
		SuccessAction:   Nothing,
	}
	storage.Update(setting)

	result, _ := storage.All()

	if result.Language != "en" {
		t.Errorf("expected cached Language to be 'en', got %s", result.Language)
	}
}
