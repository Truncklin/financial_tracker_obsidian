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

// выбранные категории (пути и имена)
const selectedCatPathSet = new Set(
  dv.array(dv.current().categories).map(c => c.path)
)
const selectedCatNames = [...selectedCatPathSet].map(p => p.split("/").pop())

// все транзакции
const allTx = dv.pages('"finance-automation/vault/finance/transactions"')

// финальная таблица — массив строк
const result = []

for (let catName of selectedCatNames) {
  // фильтруем транзакции по категории и отрицательной сумме
  const tx = allTx.where(p => 
    typeof p.amount === "number" &&
    p.amount < 0 &&
    typeof p.categories === "string" &&
    p.categories.split("/").pop() === catName
  )

  // агрегируем по месяцам
  const byMonth = {}
  for (let t of tx) {
    if (!t.date) continue
    const m = t.date.month
    byMonth[m] = (byMonth[m] || 0) + Math.abs(t.amount)
  }

  // добавляем строки в итоговый массив
  for (let [m, sum] of Object.entries(byMonth).sort((a, b) => a[0] - b[0])) {
    result.push([
      catName.replace(".md", ""),
      `${year}-${String(m).padStart(2, "0")}`,
      Math.round(sum)
    ])
  }
}

// выводим таблицу
dv.header(2, "Расходы по категориям и месяцам")
dv.table(["Категория", "Месяц", "Сумма"], result)

```
