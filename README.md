# wcag-scan-go
This app allows for quick scan of a website. 

Features include:
- [x] Concurrency safe
```go
app := App{
    Workers:  8 // Set number of workers
}
```
- [x] Regex pattern matching on URLs
```go
ExcludeRegex: []string{"/path/"} // Exlcude URLs that match to /path/
IncludeRegex: []string{"/path/"} // ONLY include URLs that match /path/
```
- [ ] Content search engine with [TF-IDF](https://en.wikipedia.org/wiki/Tf%E2%80%93idf) based content indexing
- [ ] Browser based GUI built using VueJS
- [ ] [WCAG 2.2 A and AA](https://www.w3.org/TR/WCAG22/) compliance check

# License
[The MIT License (MIT)](https://github.com/joaooliveirapro/wcag-scan-go/blob/main/LICENSE.md)

