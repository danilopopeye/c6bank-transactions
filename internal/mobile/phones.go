package mobile

type Phone struct {
	Width, Height  int
	Header, Footer int
	Month          int
}

const MonthSize = 150

var (
	IPhone13       = Phone{1170, 2532, 755, 245, 640}
	IPhone13ProMax = Phone{1284, 2778, 800, 250, 0}
	IPhone16Pro    = Phone{1206, 2622, 800, 250, 660}

	Phones = []Phone{IPhone13, IPhone13ProMax, IPhone16Pro}
)
