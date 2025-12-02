package utils

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

type H map[string]interface{}

// GetFirstElementInSlice returns the first element in a slice.
//
// useful when set optional args in a function.
func GetFirstElementInSlice[T interface{}](data []T) (t T, err error) {
	if len(data) == 0 {
		return t, errors.New("data is empty")
	}
	return data[0], nil
}

// NormalizePhone transform phone number to E.164
func NormalizePhone(phone string) (string, error) {
	// remove all non-digit non-'+' characters
	re := regexp.MustCompile(`[^\d+]`)
	cleaned := re.ReplaceAllString(phone, "")

	if !strings.HasPrefix(cleaned, "+") {
		cleaned = "+" + cleaned
	}

	e164 := regexp.MustCompile(`^\+[1-9]\d{7,14}$`)
	if !e164.MatchString(cleaned) {
		return "", errors.New("invalid phone number format")
	}

	return cleaned, nil
}

func PrintBox(title string, content string) {
	// 1. 按 \n 分割内容为行
	contentLines := strings.Split(content, "\n")

	// 2. 计算最大宽度：取 title 和所有 content 行中最长的
	boxWidth := utf8.RuneCountInString(title)
	for _, line := range contentLines {
		lineWidth := utf8.RuneCountInString(line)
		if lineWidth > boxWidth {
			boxWidth = lineWidth
		}
	}

	// 3. 渲染边框和内容
	horizontalBorder := "+" + strings.Repeat("-", boxWidth+2) + "+"
	emptyLine := "|" + strings.Repeat(" ", boxWidth+2) + "|"

	println(horizontalBorder)

	// 4. 渲染标题（居中对齐）
	titleWidth := utf8.RuneCountInString(title)
	leftPad := (boxWidth - titleWidth) / 2
	rightPad := boxWidth - titleWidth - leftPad
	println("| " + strings.Repeat(" ", leftPad) + title + strings.Repeat(" ", rightPad) + " |")

	println(horizontalBorder)

	// 5. 上边距（空行）
	println(emptyLine)

	// 6. 渲染内容（左对齐）
	for _, line := range contentLines {
		lineWidth := utf8.RuneCountInString(line)
		padding := strings.Repeat(" ", boxWidth-lineWidth)
		println("| " + line + padding + " |")
	}

	// 7. 下边距（空行）
	println(emptyLine)

	println(horizontalBorder)
}
