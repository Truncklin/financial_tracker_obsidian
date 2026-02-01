package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	csvFile      = "Pay.csv"
	vaultPath    = "vault"
	outputDir    = "finance/transactions"
	outputDirCat = "finance/categories"
	source       = "tinkoff"
	account      = "tinkoff_black"
)

func slugify(s string) string {
	s = strings.ToLower(s)
	re := regexp.MustCompile(`[^a-zа-я0-9]+`)
	s = re.ReplaceAllString(s, "_")
	return strings.Trim(s, "_")
}

func main() {

	f, err := os.Open(csvFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = ';'
	reader.FieldsPerRecord = 0

	rows, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	headers := rows[0]
	index := map[string]int{}
	for i, h := range headers {
		h = strings.TrimSpace(h)
		h = strings.TrimPrefix(h, "\ufeff")
		index[h] = i
	}

	categorySet := make(map[string]struct{})

	for _, row := range rows[1:] {
		dateStr := row[index["Дата операции"]]
		amountStr := strings.Replace(row[index["Сумма операции"]], ",", ".", 1)
		currency := row[index["Валюта операции"]]
		merchant := row[index["Описание"]]
		rawCategory := row[index["Категория"]]

		amount, _ := strconv.ParseFloat(amountStr, 64)
		dateTime, err := time.Parse("02.01.2006 15:04:05", dateStr)
		if err != nil {
			panic(err)
		}

		date := dateTime.Format("2006-01-02")

		categoryPath := fmt.Sprintf("finance-automation/vault/finance/categories/%s.md", rawCategory)
		categorySet[rawCategory] = struct{}{}

		filename := fmt.Sprintf(
			"%s__%s__%d.md",
			date,
			slugify(merchant),
			int(amount),
		)

		fullPath := filepath.Join(
			vaultPath,
			outputDir,
			filename,
		)

		if _, err := os.Stat(fullPath); err == nil {
			continue // уже есть
		}

		content := fmt.Sprintf(
			`---
type: transaction
source: %s
date: %s
amount: %.2f
currency: %s
category: [[%s]]
merchant: %s
account: %s
tags: [%s]
---
			
# %s
			
- 💸 Сумма: %.2f %s
- 📅 Дата: %s
- 🏷 Категория: [[%s]]
`,
			source,
			date,
			amount,
			currency,
			categoryPath,
			merchant,
			account,
			map[bool]string{true: "expense", false: "income"}[amount < 0],
			merchant,
			amount,
			currency,
			date,
			categoryPath)

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			panic(err)
		}
	}
	for cat := range categorySet {
		content := fmt.Sprintf(
			`---
type: %s
enabled: true
color: "%s"
---
`,
			cat,
			randomColor())

		filename := fmt.Sprintf("%s.md", cat)

		fullPath := filepath.Join(
			vaultPath,
			outputDirCat,
			filename,
		)

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			panic(err)
		}
	}
}

func randomColor() string {
	return fmt.Sprintf("#%06X", rand.Intn(0xFFFFFF))
}
