package storage

import "testing"

func TestGetConfig(t *testing.T) {
	if *GetConfig() != testConfig {
		t.Errorf("Initial configuration is not the test configuration")
	}
}

func TestPushPopConfig(t *testing.T) {
	popTest := func() {
		if err := PopConfig(); err != nil {
			t.Errorf("Failed to pop the staging configuration")
		}
		if *GetConfig() != testConfig {
			t.Errorf("Config after pop is not the test configuration")
		}
	}
	if err := PushConfig("../.env.staging"); err != nil {
		t.Errorf("Failed to push config and load staging configuration")
	}
	defer popTest()
	if *GetConfig() == testConfig {
		t.Errorf("Config after push is still the test configuration")
	}
}

func TestPushPopFailedConfig(t *testing.T) {
	popTest := func() {
		if err := PopConfig(); err == nil {
			t.Errorf("Was able to pop after failed push")
		}
	}
	if err := PushConfig(".no-such-environment-file"); err == nil {
		t.Errorf("Was able to push a non-existent environment")
	}
	defer popTest()
	if *GetConfig() != testConfig {
		t.Errorf("Config after failed push is not the test configuration")
	}
}