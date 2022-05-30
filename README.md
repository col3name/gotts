
### Install
```
go get "github.com/col3name/gotts"
```

### Update
```
go get -u "github.com/col3name/gotts"
```

### Remove
```
go clean -i "github.com/col3name/gotts"
```

### Import
```go
import "github.com/col3name/gotts"
import "github.com/col3name/gotts/voices"
```

### Use
```go
speech := gotts.Speech{Folder: "audio", Language: voices.English, Volume: 0, Speed: 1}
speech.Speak("Your sentence.")
```

### Use with Handlers
```go
import (
    "github.com/col3name/gotts"
    handlers "github.com/col3name/gotts/handlers"
    voices "github.com/col3name/gotts/voices"
)

speech := gotts.Speech{Folder: "audio", Language: voices.English, Volume: 0, Speed: 1}
speech.Speak("Your sentence.")
```

Have Fun!
