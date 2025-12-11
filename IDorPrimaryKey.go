package fmt

// IDorPrimaryKey determines if a field is an ID field and/or a primary key field.
// This function analyzes field names to identify ID fields and primary keys based on naming conventions.
//
// Parameters:
//   - tableName: The name of the table or entity that the field belongs to
//   - fieldName: The name of the field to analyze
//
// Returns:
//   - isID: true if the field is an ID field (starts with "id")
//   - isPK: true if the field is a primary key (matches specific patterns)
//
// Examples:
//   - IDorPrimaryKey("user", "id") returns (true, true)
//   - IDorPrimaryKey("user", "iduser") returns (true, true)
//   - IDorPrimaryKey("user", "userID") returns (true, true)
//   - IDorPrimaryKey("user", "USER_id") returns (true, true)
//   - IDorPrimaryKey("user", "id_user") returns (true, true)
//   - IDorPrimaryKey("user", "idaddress") returns (true, false)
//   - IDorPrimaryKey("user", "name") returns (false, false)
func IDorPrimaryKey(tableName, fieldName string) (isID, isPK bool) {
	c := GetConv()
	defer c.PutConv()

	// Convert tableName to lower in work buffer
	c.ResetBuffer(BuffOut)
	c.WrString(BuffOut, tableName)
	c.ToLower()
	c.swapBuff(BuffOut, BuffWork)
	tableLower := c.GetString(BuffWork)

	// Convert fieldName to lower in out buffer
	c.ResetBuffer(BuffOut)
	c.WrString(BuffOut, fieldName)
	c.ToLower()
	fieldLower := c.GetString(BuffOut)

	// Check if it's an ID field (starts with "id")
	if HasPrefix(fieldLower, "id") {
		isID = true
	}

	// Check for primary key patterns
	if fieldLower == "id" && tableName != "" {
		isPK = true
		return
	}

	if HasPrefix(fieldLower, "id_") && fieldLower[3:] == tableLower && tableName != "" {
		isPK = true
		return
	}

	if HasPrefix(fieldLower, "id") && fieldLower[2:] == tableLower && tableName != "" {
		isPK = true
		return
	}

	if HasSuffix(fieldLower, "id") && fieldLower[:len(fieldLower)-2] == tableLower && tableName != "" {
		isPK = true
		return
	}

	if HasSuffix(fieldLower, "_id") && fieldLower[:len(fieldLower)-3] == tableLower && tableName != "" {
		isPK = true
		return
	}

	return
}
