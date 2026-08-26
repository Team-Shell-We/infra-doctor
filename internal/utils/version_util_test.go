package utils

import (
	"testing"

	"github.com/Team-Shell-We/infra-doctor/internal/i18n"
)

func TestLocalizeVersionTranslatesUnknownFallback(t *testing.T) {

	got := LocalizeVersion(i18n.Korean, "unknown")

	if got != i18n.Get(i18n.Korean, "version.unknown") {
		t.Errorf("LocalizeVersion(ko, unknown) = %q, want %q", got, i18n.Get(i18n.Korean, "version.unknown"))
	}
}

func TestLocalizeVersionTranslatesDevelopmentFallback(t *testing.T) {

	got := LocalizeVersion(i18n.English, "development")

	if got != i18n.Get(i18n.English, "version.development") {
		t.Errorf("LocalizeVersion(en, development) = %q, want %q", got, i18n.Get(i18n.English, "version.development"))
	}
}

func TestLocalizeVersionLeavesRealVersionUnchanged(t *testing.T) {

	got := LocalizeVersion(i18n.Korean, "v0.2.0")

	if got != "v0.2.0" {
		t.Errorf("LocalizeVersion(ko, v0.2.0) = %q, want %q (real version numbers aren't translated)", got, "v0.2.0")
	}
}

// 회귀 테스트: go test 바이너리에도 빌드 정보가 있으므로 패닉 없이 값을 반환해야 함
func TestVersionAndGoVersionReturnNonEmpty(t *testing.T) {

	if v := Version(); v == "" {
		t.Error("Version() returned empty string")
	}

	if v := GoVersion(); v == "" {
		t.Error("GoVersion() returned empty string")
	}
}
