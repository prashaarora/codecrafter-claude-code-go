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
	Name        string `yaml:"name"`
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

		skill, err := parseSkillFrontmatter(string(content), folderName)
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

func parseSkillFrontmatter(content string, folderName string) (Skill, error) {
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

func resolveUserContent(prompt string) (string, error) {
	if !strings.HasPrefix(prompt, "/") {
		return prompt, nil
	}

	fields := strings.Fields(prompt)
	if len(fields) == 0 {
		return prompt, nil
	}

	skillName := strings.TrimPrefix(fields[0], "/")
	if skillName == "" {
		return "", fmt.Errorf("skill name is empty")
	}

	body, err := loadSkillBody(skillName)
	if err != nil {
		return "", err
	}
	return body, nil
}

func loadSkillBody(skillName string) (string, error) {
	skillFilePath := filepath.Join(".claude/skills", skillName, "SKILL.md")
	contentBytes, err := os.ReadFile(skillFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read skill %s: %w", skillName, err)
	}

	content := strings.ReplaceAll(string(contentBytes), "\r\n", "\n")
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return "", fmt.Errorf("missing frontmatter in skill %s", skillName)
	}

	body := strings.TrimSpace(parts[2])
	return body, nil
}
