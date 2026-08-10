package clerk

// SCIMDirectory is the former name of [Directory].
//
// Deprecated: use [Directory]. The two types are identical; this alias exists
// so that code written against the SCIM-era names keeps compiling.
type SCIMDirectory = Directory

// SCIMDirectoryList is the former name of [DirectoryList].
//
// Deprecated: use [DirectoryList]. This is a distinct type rather than an
// alias because [DirectoryList] renamed the data field from SCIMDirectories to
// Directories. Both decode the same `data` payload.
type SCIMDirectoryList struct {
	APIResource
	SCIMDirectories []*SCIMDirectory `json:"data"`
	TotalCount      int64            `json:"total_count"`
}
