package database

type Question struct {
	ID      int
	desc    string   `Description`
	answers []Answer `Allow multiple answers`
	timeout int      `Timer`
}

type Answer struct {
	ID        int
	desc      string `Description`
	isCorrect bool
}

type Quiz struct {
	ID        int      `gorm:"column:id; primary_key; not null" json:"id"`
	UUID      string   `Unique quiz code`
	questions Question `gorm:"column:role" json:"role"`
	private   bool
}
