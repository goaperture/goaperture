package main

import (
	"fmt"
	"io/fs"

	"github.com/goaperture/goaperture/v2/skills"
)

func main() {
	fmt.Println("=== Testing embed ===")

	// List all files in the embedded FS
	err := fs.WalkDir(skills.SkillFolder, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path, d.IsDir())
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println()
	fmt.Println("=== ReadDir tests ===")
	e, err := skills.SkillFolder.ReadDir("a2")
	fmt.Printf("ReadDir(\"a2\"):   %d entries, err=%v\n", len(e), err)
	e, err = skills.SkillFolder.ReadDir("a2/")
	fmt.Printf("ReadDir(\"a2/\"):  %d entries, err=%v\n", len(e), err)
	e, err = skills.SkillFolder.ReadDir(".")
	fmt.Printf("ReadDir(\".\"):    %d entries, err=%v\n", len(e), err)
	_, err = skills.SkillFolder.ReadFile("a2/SKILL.md")
	fmt.Printf("ReadFile(\"a2/SKILL.md\"): err=%v\n", err)
}
