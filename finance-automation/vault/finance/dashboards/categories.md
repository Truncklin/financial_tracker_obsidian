```TABLE
  enabled AS "ON",
  budget AS "Бюджет"
FROM "finance/categories"
WHERE type = "category"
SORT enabled DESC, file.name ASC
```
