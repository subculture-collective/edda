package srd521

// monsterPrintedXP lists stat blocks whose printed XP disagrees with the
// Experience Points by Challenge Rating table (PDF p. 256). The record keeps the
// printed CR and XP so the transcription stays literal; validation accepts only
// that exact pair for that record, and encounter budgets still reject it.
var monsterPrintedXP = map[string][2]string{
	// PDF p. 305 prints "CR 12 (XP 8,000; PB +4)"; the table gives 8,400.
	"srd-5.2.1:content/monsters/archmage": {"12", "8000"},
}
