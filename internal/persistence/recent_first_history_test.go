package persistence

import (
	"reflect"
	"testing"
)

func TestRecentFirstHistory_Append(t *testing.T) {
	// Append tests have unique expectations, so their cases are defined locally.
	testCases := []struct {
		name          string
		capacity      int
		initialItems  [][]byte
		itemToAppend  []byte
		expectedOrder [][]byte // Newest to oldest
		expectedCount int
	}{
		{"append to empty", 3, [][]byte{}, []byte("a"), [][]byte{[]byte("a")}, 1},
		{"append to non-full", 3, [][]byte{[]byte("a")}, []byte("b"), [][]byte{[]byte("b"), []byte("a")}, 2},
		{"append to full, causing eviction", 2, [][]byte{[]byte("a"), []byte("b")}, []byte("c"), [][]byte{[]byte("c"), []byte("b")}, 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestRecentFirstHistory(tc.capacity, tc.initialItems)
			h.Append(tc.itemToAppend)

			if h.count != tc.expectedCount {
				t.Errorf("Expected count %d, got %d", tc.expectedCount, h.count)
			}

			actualOrder := make([][]byte, 0, h.count)
			for i := 0; i < h.count; i++ {
				item, _ := h.GetPos(i)
				actualOrder = append(actualOrder, item.Data)
			}

			if !reflect.DeepEqual(actualOrder, tc.expectedOrder) {
				t.Errorf("Expected order %#v, got %#v", tc.expectedOrder, actualOrder)
			}
		})
	}
}

func TestRecentFirstHistory_Newest(t *testing.T) {
	for _, tc := range newestTestCases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestRecentFirstHistory(3, tc.initialItems)
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

func TestRecentFirstHistory_Contains(t *testing.T) {
	h := newTestRecentFirstHistory(3, [][]byte{[]byte("apple"), []byte("banana")})

	for _, tc := range containsTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got := h.Contains(tc.itemToFind)
			if got != tc.expectToFind {
				t.Errorf("Contains(%q) = %v, want %v", tc.itemToFind, got, tc.expectToFind)
			}
		})
	}
}

func TestRecentFirstHistory_LogicalToPhysical(t *testing.T) {
	t.Run("no wrap around", func(t *testing.T) {
		h := newTestRecentFirstHistory(5, [][]byte{[]byte("a"), []byte("b"), []byte("c")})

		if phys := h.logicalToPhysical(0); phys != 2 {
			t.Errorf("Expected logical 0 to be physical 2, got %d", phys)
		}
		if phys := h.logicalToPhysical(2); phys != 0 {
			t.Errorf("Expected logical 2 to be physical 0, got %d", phys)
		}
	})

	t.Run("with wrap around", func(t *testing.T) {
		h := newTestRecentFirstHistory(3, [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d")})

		if phys := h.logicalToPhysical(0); phys != 0 {
			t.Errorf("Expected logical 0 to be physical 0, got %d", phys)
		}
		if phys := h.logicalToPhysical(2); phys != 1 {
			t.Errorf("Expected logical 2 to be physical 1, got %d", phys)
		}
	})
}
