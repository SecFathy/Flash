package markdown

import (
	"fmt"
	"os"
)

func SaveMarkdown(fileName string, content string) error {
	filePath := fmt.Sprintf("%s.md", fileName)
	return os.WriteFile(filePath, []byte(content), 0644)
}
