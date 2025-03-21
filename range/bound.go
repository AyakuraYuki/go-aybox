package ay_range

// Bound indicates whether an endpoint of some range is contained in the range
// itself ("closed") or not ("open"). If a range is unbounded on a side, it is
// neither open nor closed on that side; the bound simply does not exist.
type Bound struct {
	inclusive bool
}

var (
	BoundOpen   Bound = Bound{false}
	BoundClosed Bound = Bound{true}
)

// BoundForBool returns the bound type corresponding to a boolean value for
// inclusivity.
func BoundForBool(inclusive bool) Bound {
	if inclusive {
		return BoundClosed
	}
	return BoundOpen
}
