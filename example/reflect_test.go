package example

import "testing"

//  author: LaoGou
//  since: 2022/10/31
//  desc: example

func TestJSONArray(t *testing.T) {
	j := `[{"name":"的sdd","age":0,"activity":true,"fighting":33.55,"inDate":"2022-08-31T09:08:29.837820+00:00"},{"name":"的sdd","age":0,"activity":true,"fighting":33.55,"inDate":"2022-08-31T09:08:29.837820+00:00"}]`
	t.Log(j)
}
