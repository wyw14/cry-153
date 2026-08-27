package gate

import "time"

func RetrySchedule(retries int, max time.Duration) []time.Duration {
	if retries < 0 {
		retries = 0
	}
	out := make([]time.Duration, retries)
	for i := range out {
		out[i] = BoundedBackoff(i, max)
	}
	return out
}
func TotalBackoff(retries int, max time.Duration) time.Duration {
	total := time.Duration(0)
	for _, d := range RetrySchedule(retries, max) {
		if max-total < d {
			return max
		}
		total += d
	}
	return total
}
func IsSaturated(attempt int, max time.Duration) bool { return BoundedBackoff(attempt, max) == max }
