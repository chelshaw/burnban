# Burn Ban

## How to run this

1. Pull down the repo

2. Install deps `go mod vendor`

3. Run it `go run .`

## Deploy

1. Build binary `go build -o bin/burnban -v`

2. `git push heroku main`

### NEXT STEPS / NOTES

- Figure out local vendoring bullshit

- go mod edit -replace example.com/greetings=../greetings
