# Autoclone

## Run autoclone
To download and run autoclone server at your machine you need to add `config.json`,
download the executable file `rm -f Autoclone && wget https://github.com/GestureToSpeech/Autoclone/raw/master/bin/Autoclone -O Autoclone && sudo chmod -R 0777 Autoclone`
and run it `rm -f nohup.out && nohup ./Autoclone &`. This will start the server in the background. Logs will be in
`nohup.out`. To stop it, run `kill $(pgrep Autoclone)`.

## Development
To install all dependencies for development: `sudo bash install.sh`.

Add bin path: `go env -w GOBIN=/path/to/folder/Autoclone/bin`

Compile to executable: `go install Autoclone`
