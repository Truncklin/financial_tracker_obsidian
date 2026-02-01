---
view: monthly
categories:
  - "[[Образование]]"
  - "[[Рестораны]]"
  - "[[Бонусы]]"
  - "[[Проценты]]"
  - "[[Связь]]"
  - "[[Наличные]]"
  - "[[Спорттовары]]"
  - "[[Красота]]"
  - "[[Сервис]]"
  - "[[Местный транспорт]]"
  - "[[Тренировки]]"
  - "[[Различные товары]]"
  - "[[Супермаркеты]]"
  - "[[Переводы]]"
  - "[[Экосистема Яндекс]]"
  - "[[Развлечения]]"
  - "[[Фастфуд]]"
  - "[[Фото и копицентры]]"
  - "[[Транспорт]]"
year: 2025
---

```dataviewjs
const year = dv.current().year ?? 2025

dv.header(2, "DEBUG 1 — текущий год")
dv.paragraph(year)

// ----------------------------
// выбранные категории (ИСТИНА)
// ----------------------------
const selectedCategories = (dv.current().categories ?? [])
  .map(c => c.path.split("/").pop().replace(".md", ""))

dv.header(2, "DEBUG 2 — выбранные категории")
dv.list(selectedCategories)

// ----------------------------
// все транзакции
// ----------------------------
const allTx = dv.pages('"finance-automation/vault/finance/transactions"')

dv.header(2, "DEBUG 3 — всего транзакций найдено")
dv.paragraph(allTx.length)

// ----------------------------
// helper: извлечь имена категорий ИЗ ЛЮБОГО ФОРМАТА
// ----------------------------
function extractCategoryNames(field) {
  return dv.array(field)
    .map(v => {
      if (typeof v === "string") return v
      if (v?.path) return v.path
      return null
    })
    .filter(Boolean)
    .map(p => p.split("/").pop().replace(".md", ""))
}

// ----------------------------
// DEBUG — что реально в категориях
// ----------------------------
dv.header(2, "DEBUG 4 — категории транзакций")
dv.table(
  ["Файл", "RAW", "catNames"],
  allTx.slice(0, 10).map(p => [
    p.file.name,
    JSON.stringify(p.category),
    extractCategoryNames(p.category).join(", ")
  ])
)

// ----------------------------
// фильтрация
// ----------------------------
const tx = allTx.where(p => {
  if (typeof p.amount !== "number" || p.amount >= 0) return false

  const catNames = extractCategoryNames(p.category)
  if (catNames.length === 0) return false

  return catNames.some(c => selectedCategories.includes(c))
})

dv.header(2, "DEBUG 5 — прошло фильтр")
dv.paragraph(tx.length)

// ----------------------------
// транзакции после фильтра
// ----------------------------
dv.header(2, "DEBUG 6 — транзакции после фильтра")
dv.table(
  ["Файл", "Дата", "Сумма", "Категории"],
  tx.slice(0, 10).map(p => [
    p.file.name,
    p.date?.toISODate?.() ?? p.date,
    p.amount,
    extractCategoryNames(p.category).join(", ")
  ])
)

// ----------------------------
// агрегация по месяцам
// ----------------------------
const byMonth = {}

for (let t of tx) {
  if (!t.date) continue
  const m = t.date.month
  byMonth[m] = (byMonth[m] || 0) + Math.abs(t.amount)
}

dv.header(2, "DEBUG 7 — агрегированные данные")
dv.paragraph(JSON.stringify(byMonth, null, 2))

// ----------------------------
// финальная таблица
// ----------------------------
dv.header(2, "RESULT — траты по месяцам")
dv.table(
  ["Месяц", "Сумма"],
  Object.entries(byMonth)
    .sort((a, b) => a[0] - b[0])
    .map(([m, sum]) => [
      `${year}-${String(m).padStart(2, "0")}`,
      Math.round(sum)
    ])
)
```
