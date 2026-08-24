package analyzer

import (
	"encoding/xml"
	"strings"
)

type mavenProject struct {
	XMLName      xml.Name          `xml:"project"`
	ModelVersion string            `xml:"modelVersion"`
	GroupID      string            `xml:"groupId"`
	ArtifactID   string            `xml:"artifactId"`
	Version      string            `xml:"version"`
	Packaging    string            `xml:"packaging"`
	Parent       mavenParent       `xml:"parent"`
	Properties   mavenProperties   `xml:"properties"`
	Dependencies []mavenDependency `xml:"dependencies>dependency"`
	Plugins      []mavenPlugin     `xml:"build>plugins>plugin"`
	Modules      []string          `xml:"modules>module"`
}

type mavenParent struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
}

type mavenDependency struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
	Scope      string `xml:"scope"`
	Optional   bool   `xml:"optional"`
}

type mavenPlugin struct {
	GroupID       string `xml:"groupId"`
	ArtifactID    string `xml:"artifactId"`
	Version       string `xml:"version"`
	Configuration struct {
		Release string `xml:"release"`
		Source  string `xml:"source"`
		Target  string `xml:"target"`
	} `xml:"configuration"`
}

type mavenProperties map[string]string

func (properties *mavenProperties) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	values := make(map[string]string)
	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch element := token.(type) {
		case xml.StartElement:
			var value string
			if err := decoder.DecodeElement(&value, &element); err != nil {
				return err
			}
			values[element.Name.Local] = strings.TrimSpace(value)
		case xml.EndElement:
			if element.Name == start.Name {
				*properties = values
				return nil
			}
		}
	}
}
