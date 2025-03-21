package aybox

import (
	"github.com/dromara/carbon/v2"
	"testing"
)

func TestSetTimezone(t *testing.T) {
	c := carbon.Now()
	t.Log(c.Format("Y-m-d H:i:s T"))
	c = c.SetTimezone("Asia/Tokyo")
	c = c.SetTimezone("Asia/Tokyo")
	t.Log(c.Format("Y-m-d H:i:s T"))
}
