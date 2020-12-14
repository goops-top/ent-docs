// Code generated by entc, DO NOT EDIT.

package owner

const (
	// Label holds the string label denoting the owner type in the database.
	Label = "owner"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"

	// Table holds the table name of the owner in the database.
	Table = "owners"
)

// Columns holds all SQL columns for owner fields.
var Columns = []string{
	FieldID,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}
