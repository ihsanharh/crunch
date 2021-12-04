package emojis

var (
	PingPong = "🏓"
	music = "🎵"
	xmark     = "✖️"
	loading = "🥏"
	check      = "✅"
)

func GetNumberAsEmoji(num int) string {
	switch num {
	case 0:
		return "0⃣"
	case 1:
		return "1⃣"
	case 2:
		return "2⃣"
	case 3:
		return "3⃣"
	case 4:
		return "4⃣"
	case 5:
		return "5⃣"
	case 6:
		return "6⃣"
	case 7:
		return "7⃣"
	case 8:
		return "8⃣"
	case 9:
		return "9⃣"
	case 10:
		return "🔟"
	default:
		return ""
	}
}
