# Buy a PS5 Digital Edition

1. Go into [pkg/util/cookies.go](https://github.com/drpaneas/ps5amazon/blob/main/pkg/util/cookies.go) and populate the variables with your session cookie, taken from Developer Tools (e.g. using Chrome)

```go
var (
	EuronicsCookie  string = "euronics=blahblahblah"
	AmazonCookie    string = "session-id=blahclahblah"
	AlternateCookie string = "secure_session=blahclahblah"
)
```

2. Build: `go build`
3. Run it and make sure your speakers are on.

### Test speakers

To make sure your speakers are working fine with this software, change one of the if statements at `main.go`. Build. Run.
For example:

```diff
- if amazon.IsReadyToBuy() {
+ if !amazon.IsReadyToBuy() {
```

You should hear the siren