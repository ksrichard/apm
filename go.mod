module github.com/ksrichard/apm

go 1.15

require (
	github.com/arduino/arduino-cli v0.0.0-20210413144851-088d4276190d
	github.com/manifoldco/promptui v0.8.0
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1 // indirect
	google.golang.org/grpc v1.27.0
)

replace go.bug.st/downloader/v2 => ./go-downloader/
