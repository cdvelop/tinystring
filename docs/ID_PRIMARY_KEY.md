# ID and Primary Key Detection

The `IDorPrimaryKey` function determines if a field is an ID field and/or a primary key field based on naming conventions.

#### Parameters
- `tableName`: The name of the table or entity that the field belongs to
- `fieldName`: The name of the field to analyze

#### Returns
- `isID`: `true` if the field is an ID field (starts with "id")
- `isPK`: `true` if the field is a primary key (matches specific patterns)

#### Examples
- `IDorPrimaryKey("user", "id")` returns `(true, true)`
- `IDorPrimaryKey("user", "iduser")` returns `(true, true)`
- `IDorPrimaryKey("user", "userId")` returns `(true, true)`
- `IDorPrimaryKey("user", "id_user")` returns `(true, true)`
- `IDorPrimaryKey("user", "user_ID")` returns `(true, true)`
- `IDorPrimaryKey("user", "idaddress")` returns `(true, false)`
- `IDorPrimaryKey("user", "name")` returns `(false, false)`