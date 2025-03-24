package appointment

import "sort"

// Helper function to sort queue items
func sortQueueItems(items []QueueItem) {
	// Define priority order
	priorityOrder := map[string]int{
		"Normal": 1,
		"High":   2,
	}
	// Sort by priority (high first) and then by waiting time
	sort.Slice(items, func(i, j int) bool {
		// If priorities are different, sort by priority (high first)
		if items[i].Priority != items[j].Priority {
			return priorityOrder[items[i].Priority] > priorityOrder[items[j].Priority]
		}
		// If priorities are the same, sort by waiting time (longer wait first)
		return items[i].WaitingSince < items[j].WaitingSince
	})
}
