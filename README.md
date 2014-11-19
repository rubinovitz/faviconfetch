# FaviconFetch


FaviconFetch is a go package for fetching a favicon given a url.

## Examples

FaviconFetch can be run as a gorountine and have the favicon returned via a channel.

```
import(
    "github.com/rubinovitz/faviconfetch"
)

// make a channel
faviconChannel := make(chan []byte)
// launch a gorountine to get the favicon
go func() {
        faviconChannel <- Fetch(uri)
}()
//do whatever for a while
// get the favicon
favicon := <-faviconChannel

```

## Testing

Run
```
go test
```
