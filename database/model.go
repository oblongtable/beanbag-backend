package database

type Quiz struct {
	ID        int    `gorm:"column:id; primary_key; not null" json:"id"`
	UUID      string `Unique quiz code`
	title     string
	desc      string
	questions []Question `gorm:"column:role" json:"role"`
	isPriv    bool
	timer     int
}

type Question struct {
	ID           int
	desc         string   `Description`
	answers      []Answer `Allow multiple answers`
	sp_timer_opt bool     `Specific question timer option`
	sp_timer     int      `Only used if sp_timer_opt enabled`
}

type Answer struct {
	ID        int
	desc      string `Description`
	isCorrect bool
}
