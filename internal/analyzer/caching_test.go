package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeCachingCountsCacheAnnotations(t *testing.T) {

	root := t.TempDir()

	src := `package com.example.demo;

import org.springframework.cache.annotation.*;

@CacheConfig(cacheNames = "facilities")
public class FacilityService {

    @Cacheable(key = "#id")
    public Object find(Long id) { return null; }

    @CacheEvict(allEntries = true)
    public void clear() {}
}
`
	if err := os.WriteFile(filepath.Join(root, "FacilityService.java"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	info, err := AnalyzeCaching(root)
	if err != nil {
		t.Fatalf("AnalyzeCaching failed: %v", err)
	}

	if info.CacheableCount != 2 {
		t.Errorf("CacheableCount = %d, want 2 (Cacheable + CacheEvict)", info.CacheableCount)
	}
	if !info.HasCacheConfig {
		t.Error("HasCacheConfig = false, want true")
	}
}

func TestAnalyzeCachingDetectsTokenCandidatesByFilename(t *testing.T) {

	root := t.TempDir()

	for _, name := range []string{"RefreshTokenRepository.java", "UserSession.java", "FacilityController.java"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("package com.example.demo;\npublic interface X {}\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	info, err := AnalyzeCaching(root)
	if err != nil {
		t.Fatalf("AnalyzeCaching failed: %v", err)
	}

	if len(info.TokenCandidates) != 2 {
		t.Fatalf("len(TokenCandidates) = %d, want 2, got %v", len(info.TokenCandidates), info.TokenCandidates)
	}
}

func TestAnalyzeCachingNoFalsePositivesOnPlainClass(t *testing.T) {

	root := t.TempDir()

	src := "package com.example.demo;\npublic class FacilityController {}\n"
	if err := os.WriteFile(filepath.Join(root, "FacilityController.java"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	info, err := AnalyzeCaching(root)
	if err != nil {
		t.Fatalf("AnalyzeCaching failed: %v", err)
	}

	if info.CacheableCount != 0 || info.HasCacheConfig || len(info.TokenCandidates) != 0 {
		t.Errorf("expected all-zero CachingInfo, got %+v", info)
	}
}
