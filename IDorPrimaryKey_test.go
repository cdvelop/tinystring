package tinystring

import "testing"

// IDorPrimaryKeyTestCases returns test cases for IDorPrimaryKey function
func IDorPrimaryKeyTestCases() []struct {
	tableName  string
	fieldName  string
	expectedID bool
	expectedPK bool
} {
	return []struct {
		tableName  string
		fieldName  string
		expectedID bool
		expectedPK bool
	}{
		// Test "ID" field (case insensitive)
		{"user", "ID", true, true},
		{"user", "id", true, true},
		{"user", "Id", true, true},
		{"user", "iD", true, true},

		// Test "id_TABLE_NAME" pattern (case insensitive)
		{"user", "id_user", true, true},
		{"user", "ID_USER", true, true},
		{"user", "Id_User", true, true},
		{"product", "id_product", true, true},
		{"product", "ID_PRODUCT", true, true},

		// Test "idTABLE_NAME" pattern (case insensitive)
		{"user", "iduser", true, true},
		{"user", "IDUSER", true, true},
		{"user", "Iduser", true, true},
		{"product", "idproduct", true, true},
		{"product", "IDPRODUCT", true, true},

		// Test "TABLE_NAMEid" pattern (case insensitive)
		{"user", "userid", false, true},
		{"user", "USERID", false, true},
		{"user", "Userid", false, true},
		{"user", "user_id", false, true},
		{"user", "USER_ID", false, true},
		{"product", "productid", false, true},
		{"product", "PRODUCTID", false, true},

		// Test non-ID fields
		{"user", "name", false, false},
		{"user", "email", false, false},
		{"user", "id_other", true, false}, // Starts with id but not PK for user
		{"user", "", false, false},        // Empty field name
		{"", "ID", true, false},
		{"i", "i", false, false},
	}
}

func TestIDorPrimaryKey(t *testing.T) {
	testCases := IDorPrimaryKeyTestCases()

	for _, tc := range testCases {
		t.Run(tc.tableName+"_"+tc.fieldName, func(t *testing.T) {
			isID, isPK := IDorPrimaryKey(tc.tableName, tc.fieldName)
			if isID != tc.expectedID || isPK != tc.expectedPK {
				t.Errorf("IDorPrimaryKey(%q, %q) = (%v, %v); want (%v, %v)",
					tc.tableName, tc.fieldName, isID, isPK, tc.expectedID, tc.expectedPK)
			}
		})
	}
}

func BenchmarkIDorPrimaryKey(b *testing.B) {
	testCases := IDorPrimaryKeyTestCases()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			IDorPrimaryKey(tc.tableName, tc.fieldName)
		}
	}
}
