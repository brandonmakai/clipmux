package persistence

import (
	"reflect"
	"testing"
)

func TestChronologicalHistory_Append(t *testing.T) {
	// Append tests have unique expectations, so their cases are defined locally.
	testCases := []struct {
		name          string
		capacity      int
		initialItems  [][]byte
		itemToAppend  []byte
		expectedOrder [][]byte
		expectedCount int
	}{
		{"append to empty", 3, [][]byte{}, []byte("a"), [][]byte{[]byte("a")}, 1},
		{"append to non-full", 3, [][]byte{[]byte("a")}, []byte("b"), [][]byte{[]byte("a"), []byte("b")}, 2},
		{"append to full, causing eviction", 2, [][]byte{[]byte("a"), []byte("b")}, []byte("c"), [][]byte{[]byte("b"), []byte("c")}, 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestChronologicalHistory(tc.capacity, tc.initialItems)
			h.Append(tc.itemToAppend)

			if h.count != tc.expectedCount {
				t.Errorf("Expected count %d, got %d", tc.expectedCount, h.count)
			}

			actualOrder := make([][]byte, 0, h.count)
			for i := 0; i < h.count; i++ {
				idx := (h.head + i) % h.capacity
				actualOrder = append(actualOrder, h.buf[idx].Data)
			}

			if !reflect.DeepEqual(actualOrder, tc.expectedOrder) {
				t.Errorf("Expected order %#v, got %#v", tc.expectedOrder, actualOrder)
			}
		})
	}
}

func TestChronologicalHistory_Newest(t *testing.T) {
	for _, tc := range newestTestCases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestChronologicalHistory(3, tc.initialItems)
			item, found := h.Newest()

			if found != tc.expectedFound {
				t.Fatalf("Expected found to be %v, got %v", tc.expectedFound, found)
			}
			if !reflect.DeepEqual(item.Data, tc.expectedItem) {
				t.Errorf("Expected item %#v, got %#v", tc.expectedItem, item.Data)
			}
		})
	}
}

func TestChronologicalHistory_Contains(t *testing.T) {
	h := newTestChronologicalHistory(3, [][]byte{[]byte("apple"), []byte("banana")})

	for _, tc := range containsTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got := h.Contains(tc.itemToFind)
			if got != tc.expectToFind {
				t.Errorf("Contains(%q) = %v, want %v", tc.itemToFind, got, tc.expectToFind)
			}
		})
	}
}
