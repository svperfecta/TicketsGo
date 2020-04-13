package command

type Category string

const (
	General    Category = "ℹ️ General"
	Tickets    Category = "📩 Tickets"
	Settings   Category = "🔧 Settings"
	Tags       Category = "✍️ Tags"
	Statistics Category = "📈 Statistics"
)

var categories = []Category{
	General,
	Tickets,
	Settings,
	Tags,
	Statistics,
}
