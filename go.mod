module gitlab.com/codifierr/cowin_bot/v1

go 1.16

replace (
	cowin => ./cowin
	telegram => ./telegram
)

require (
	cowin v0.0.0-00010101000000-000000000000 // indirect
	telegram v0.0.0-00010101000000-000000000000
)
