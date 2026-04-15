package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type Skill struct {
	ID          string
	Name        string
	Description string
	Trigger     string
	FilePath    string
	Enabled     bool
}

type Registry struct {
	Skills []Skill
}

func LoadRegistry(skillsDir string) (*Registry, error) {
	// TODO: Load skills.yaml
	// For now: return empty registry
	return &Registry{Skills: []Skill{}}, nil
}

func (r *Registry) Match(input string) *Skill {
	// TODO: Match input against skill triggers (regex or keyword)
	for i := range r.Skills {
		if !r.Skills[i].Enabled {
			continue
		}
		if r.matchesTrigger(&r.Skills[i], input) {
			return &r.Skills[i]
		}
	}
	return nil
}

func (r *Registry) matchesTrigger(skill *Skill, input string) bool {
	if skill.Trigger == "" {
		return false
	}
	matched, _ := regexp.MatchString(skill.Trigger, input)
	return matched
}

func LoadSkillFile(path string) (string, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", fmt.Errorf("reading skill file: %w", err)
	}
	return string(data), nil
}
