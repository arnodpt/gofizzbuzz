package http

// FizzbuzzQueryParams are the accepted query parameters for a Fizzbuzz query
// There are pointers in the struct to detect empty parameters on query parsing
type FizzbuzzQueryParams struct {
	Str1  *string `query:str1`
	Str2  *string `query:str2`
	Int1  *int    `query:int1`
	Int2  *int    `query:int2`
	Limit *int    `query:limit`
}

// FizzbuzzResult is returned by the Fizzbuzz request on success
type FizzbuzzResult struct {
	Result []string `json:result`
}

// PopularResult is returned by the PopularRequest request on success
type PopularResult struct {
	Str1  string `json:str1`
	Str2  string `json:str2`
	Int1  int    `json:int1`
	Int2  int    `json:int2`
	Limit int    `json:limit`
	Count int64  `json:count`
}