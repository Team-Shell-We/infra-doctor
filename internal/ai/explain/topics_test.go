package explain

import "testing"

func TestEveryTopicHasADisplayName(t *testing.T) {

	for _, topic := range Topics {
		if _, ok := displayNames[topic]; !ok {
			t.Errorf("topic %q has no entry in displayNames", topic)
		}
	}
}

func TestDisplayNameFallback(t *testing.T) {

	if got := DisplayName("not-a-real-topic"); got != "not-a-real-topic" {
		t.Errorf("DisplayName should fall back to the raw topic, got %q", got)
	}
}

func TestDisplayNameKnownTopic(t *testing.T) {

	if got := DisplayName("k8s"); got != "Kubernetes" {
		t.Errorf(`DisplayName("k8s") = %q, want "Kubernetes"`, got)
	}
}
