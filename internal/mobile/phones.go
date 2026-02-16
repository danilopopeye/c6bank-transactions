package mobile

type Phone struct {
	Width, Height  int
	Header, Footer int
	Month          int
	MonthSize      int // Height of month region to extract (100 for IPhoneMirror, 150 for others)
}

const MonthSize = 150

var (
	// IPhone13: iPhone 13 (1170×2532)
	// Header=755px, Footer=245px, Month starts at Y=640, MonthSize=150px
	IPhone13 = Phone{1170, 2532, 755, 245, 640, 150}

	// IPhone13ProMax: iPhone 13 Pro Max (1284×2782)
	// Header=800px, Footer=250px, Month disabled (Y=0), MonthSize=150px
	IPhone13ProMax = Phone{1284, 2778, 800, 250, 0, 150}

	// IPhone15Pro: iPhone 15 Pro (1179×2556)
	// Header=776px, Footer=250px, Month starts at Y=660, MonthSize=150px
	IPhone15Pro = Phone{1179, 2556, 776, 250, 660, 150}

	// IPhone16Pro: iPhone 16 Pro (1206×2622)
	// Header=800px, Footer=250px, Month starts at Y=660, MonthSize=150px
	IPhone16Pro = Phone{1206, 2622, 800, 250, 660, 150}

	// IPhoneMirror: iPhone Mirror screenshots from macOS
	// Dimensions: 836×1840 (smaller than physical screens)
	// Regions: Header=600px, Footer=180px, Month starts at Y=0, MonthSize=100px
	// Characteristic: Transparent pixels at top (first row)
	IPhoneMirror = Phone{836, 1840, 600, 180, 0, 100}

	Phones = []Phone{IPhone13, IPhone13ProMax, IPhone15Pro, IPhone16Pro, IPhoneMirror}
)
