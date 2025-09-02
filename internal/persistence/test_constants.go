package persistence

// Shared test case structure for testing the Newest() method.
var newestTestCases = []struct {
	name          string
	initialItems  [][]byte
	expectedItem  []byte
	expectedFound bool
}{
	{"empty history", [][]byte{}, nil, false},
	{"single item", [][]byte{[]byte("a")}, []byte("a"), true},
	{"multiple items", [][]byte{[]byte("a"), []byte("b"), []byte("c")}, []byte("c"), true},
}

// Shared test case structure for testing the Contains() method.
var containsTestCases = []struct {
	name         string
	itemToFind   string
	expectToFind bool
}{
	{"item exists", "apple", true},
	{"item exists 2", "banana", true},
	{"item does not exist", "cherry", false},
	{"empty string search", "", false},
}

// newTestChronologicalHistory is a helper to create and populate a ChronologicalHistory instance for testing.
func newTestChronologicalHistory(capacity int, items [][]byte) *ChronologicalHistory {
	h := newChronologicalHistory(capacity, 1000, 50)
	for _, item := range items {
		h.Append(item)
	}
	return h
}

// newTestRecentFirstHistory is a helper to create and populate a RecentFirstHistory instance for testing.
func newTestRecentFirstHistory(capacity int, items [][]byte) *RecentFirstHistory {
	h := newRecentFirstHistory(capacity, 1000, 50)
	for _, item := range items {
		h.Append(item)
	}
	return h
}
