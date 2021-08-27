package gateway_test

//type apinumtest struct {
//	num int
//	max interface{}
//	err error
//}
//
//var apinumdata = []apinumtest{
//	// valid request
//	{
//		num: 2,
//		max: 4,
//		err: nil,
//	},
//	// valid request(boundary)
//	{
//		num: 2,
//		max: 3,
//		err: nil,
//	},
//	// valid request(permit unlimited call)
//	{
//		num: 2,
//		max: "-",
//		err: nil,
//	},
//	// invalid request
//	{
//		num: 2,
//		max: 1,
//		err: errors.New("limit exceeded"),
//	},
//	// invalid request(boundary)
//	{
//		num: 2,
//		max: 2,
//		err: errors.New("limit exceeded"),
//	},
//	// unexpected value in data
//	{
//		num: 2,
//		max: "unlimited",
//		err: errors.New("unexpected limit value"),
//	},
//	// unexpected value in data
//	{
//		num: 2,
//		max: false,
//		err: errors.New("unexpected limit value"),
//	},
//}
//
//func TestAPILimitChecker(t *testing.T) {
//	for i, tt := range apinumdata {
//		gateway.APIData["key"] = []gateway.Field{
//			{
//				Template: *gateway.NewURITemplate("/path"),
//				Path:     *gateway.NewURITemplate("/path"),
//				Num:      tt.num,
//				Max:      tt.max,
//			},
//		}
//
//		switch err := gateway.APILimitChecker("key", "path"); err {
//		case nil:
//			if tt.err != nil {
//				t.Fatalf("case %d: expected %v, get %v", i, tt.err, err)
//			}
//		default:
//			if tt.err == nil {
//				t.Fatalf("case %d: expected %v, get %v", i, tt.err, err)
//			}
//			if err.Error() != tt.err.Error() {
//				t.Fatalf("case %d: expected %v, get %v", i, tt.err, err)
//			}
//		}
//	}
//}
