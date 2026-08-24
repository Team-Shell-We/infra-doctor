package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeGradleDetectsSpringBootGroovyDSL(t *testing.T) {

	root := t.TempDir()
	buildFile := filepath.Join(root, "build.gradle")

	content := `plugins {
	id 'org.springframework.boot' version '3.2.0'
	id 'java'
}
`
	if err := os.WriteFile(buildFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	framework, _, _, err := AnalyzeGradle(buildFile)
	if err != nil {
		t.Fatalf("AnalyzeGradle failed: %v", err)
	}

	if !framework.SpringBoot.Enabled || framework.SpringBoot.Version != "3.2.0" {
		t.Errorf("SpringBoot = %+v, want Enabled=true Version=3.2.0", framework.SpringBoot)
	}
}

func TestAnalyzeGradleDetectsSpringBootKotlinDSL(t *testing.T) {

	root := t.TempDir()
	buildFile := filepath.Join(root, "build.gradle.kts")

	content := `plugins {
	id("org.springframework.boot") version "3.2.0"
	id("java")
}
`
	if err := os.WriteFile(buildFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	framework, _, _, err := AnalyzeGradle(buildFile)
	if err != nil {
		t.Fatalf("AnalyzeGradle failed: %v", err)
	}

	if !framework.SpringBoot.Enabled || framework.SpringBoot.Version != "3.2.0" {
		t.Errorf("SpringBoot = %+v, want Enabled=true Version=3.2.0", framework.SpringBoot)
	}
}

func TestAnalyzeGradleJavaVersionFromToolchain(t *testing.T) {

	root := t.TempDir()
	buildFile := filepath.Join(root, "build.gradle")

	content := `java {
	toolchain {
		languageVersion = JavaLanguageVersion.of(21)
	}
}
`
	if err := os.WriteFile(buildFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	framework, _, _, err := AnalyzeGradle(buildFile)
	if err != nil {
		t.Fatalf("AnalyzeGradle failed: %v", err)
	}

	if framework.Java.Version != "21" {
		t.Errorf("Java.Version = %q, want 21", framework.Java.Version)
	}
}

func TestAnalyzeGradleJavaVersionFromSourceCompatibility(t *testing.T) {

	root := t.TempDir()
	buildFile := filepath.Join(root, "build.gradle")

	content := "sourceCompatibility = '17'\n"
	if err := os.WriteFile(buildFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	framework, _, _, err := AnalyzeGradle(buildFile)
	if err != nil {
		t.Fatalf("AnalyzeGradle failed: %v", err)
	}

	if framework.Java.Version != "17" {
		t.Errorf("Java.Version = %q, want 17", framework.Java.Version)
	}
}

func TestAnalyzeGradleDependenciesAndDatabase(t *testing.T) {

	root := t.TempDir()
	buildFile := filepath.Join(root, "build.gradle.kts")

	content := `dependencies {
	implementation("org.springframework.boot:spring-boot-starter-security")
	implementation("org.springframework.boot:spring-boot-starter-data-jpa")
	implementation("org.springframework.boot:spring-boot-starter-data-redis")
	runtimeOnly("org.postgresql:postgresql")
}
`
	if err := os.WriteFile(buildFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, dependency, database, err := AnalyzeGradle(buildFile)
	if err != nil {
		t.Fatalf("AnalyzeGradle failed: %v", err)
	}

	if !dependency.Security.Enabled {
		t.Error("expected Security dependency to be detected in Kotlin DSL")
	}
	if !dependency.JPA.Enabled {
		t.Error("expected JPA dependency to be detected in Kotlin DSL")
	}
	if database.Primary.Type != "PostgreSQL" {
		t.Errorf("database.Primary.Type = %q, want PostgreSQL", database.Primary.Type)
	}
	if database.Redis == nil || !database.Redis.Enabled {
		t.Error("expected Redis to be detected in Kotlin DSL")
	}
}
