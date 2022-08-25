package cut

// todo: optimization func

type searchRange struct {
	data map[[2]int]struct{}
}

func newSearchRange() *searchRange {
	return &searchRange{
		data: map[[2]int]struct{}{},
	}
}

// Adds range from min to max inclusive
func (r *searchRange) addRange(min, max int) {
	r.data[[2]int{min, max}] = struct{}{}
}

func (r *searchRange) addSingle(a int) {
	r.addRange(a, a)
}

func (r *searchRange) hasRange(min, max int) bool {
	_, has := r.data[[2]int{min, max}]
	return has
}

func (r *searchRange) isInRange(n int) bool {
	for k := range r.data {
		if n >= k[0] && n <= k[1] {
			return true
		}
	}
	return false
}
