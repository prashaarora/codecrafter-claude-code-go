package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Skill struct {
	Name string `yaml:"name"`
	Description string `yaml:"description"`
}

func loadSkillsSystemPrompt() string {
	skills, err := loadSkills()
	if err != nil || len(skills) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("You have access to the following skills:\n")
	for _, skill := range skills {
		b.WriteString(fmt.Sprintf("- %s: %s\n", skill.Name, skill.Description))
	}
	return b.String()
}

func loadSkills() ([]Skill, error) {
	skillsDir := ".claude/skills"
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var skills []Skill
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		folderName := entry.Name()
		skillFilePath := filepath.Join(skillsDir, folderName, "SKILL.md")

		content, err := os.ReadFile(skillFilePath)
		if err != nil {
			continue
		}

		skill, err := parseSkillFomatter(string(content), folderName)
		if err != nil {
			continue
		}
		skills = append(skills, skill)
	}
	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Name < skills[j].Name
	})
	return skills, nil
}

func parseSkillFomatter(content string, folderName string) (Skill, error) {
	content = strings.ReplaceAll(content, "\r\n", "\n")

	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return Skill{}, fmt.Errorf("missing frontmatter")
	}

	frontmatter := strings.TrimSpace(parts[1])
	var skill Skill
	if err := yaml.Unmarshal([]byte(frontmatter), &skill); err != nil {
		return Skill{}, err
	}
	if skill.Name == "" {
		skill.Name = folderName
	}
	return skill, nil
}


