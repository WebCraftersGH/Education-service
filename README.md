# Education-service
For LearnFlow

# Entities

- Problem (Задача из пулла задач)
  - ID (ID)
  - Name (Название)
  - Slug (Уникальный id)
  - Difficulty (Уровень сложности)
  - Tag (Тег, для сортировки)
  - AuthorID (ID автора)
  - VerifiedAt (Дата прохождения модерации)
  - CreatedAt (Дата создания)
  - UpdatedAt (Дата обновления)

- ProblemContent (Задание задачи)
  - ID (ID)
  - ProblemID (ID задачи)
  - DescriptionMD (Описание в формате .md)
  - InputFormatMD ()
  - OutputFormatMD ()
  - ConstraintsMD ()
  - NotesMD ()
  - CreatedAt ()
  - UpdatedAt ()

