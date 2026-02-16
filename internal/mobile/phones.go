package mobile

// Phone represents an iPhone model with screen dimensions and crop regions for transaction processing.
// The struct stores dimensions and margins for smart image cropping to extract transaction data.
type Phone struct {
	// Width, Height are the image dimensions in pixels
	Width, Height int
	// Header is the top margin to exclude (in pixels) when cropping transaction area
	Header int
	// Footer is the bottom margin to exclude (in pixels) when cropping transaction area
	Footer int
	// Month is the Y-position (vertical offset) where the month reference region starts
	Month int
	// MonthSize is the height of the month region to extract (in pixels)
	// - 100px for IPhoneMirror (smaller header area)
	// - 150px for other iPhone models
	// - 0 means use fallback default of 150px
	MonthSize int
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
