package mcache

// ReqResp - response of request
type ReqResp struct {
	Success bool    `json:"success"`
	Error   *string `json:"error,omitempty"`
	Result  *string `json:"result,omitempty"`
}

// Command -
type Command struct {
	// set, get, remove, keys
	Name string `json:"name"`

	// str, slice, map
	Type *string `json:"type,omitempty"`

	//
	Key *string `json:"key,omitempty"`
	TTL *int64  `json:"ttl,omitempty"`

	//
	Str *string `json:"str,omitempty"`

	//
	Slice []string `json:"slice,omitempty"`

	//
	Map map[string]string `json:"map,omitempty"`

	// for get command, slice index or map key
	IdxKey *string `json:"idx_key,omitempty"`
}

// PStr -
func PStr(s string) *string {
	return &s
}

// PInt64 -
func PInt64(i int64) *int64 {
	return &i
}
