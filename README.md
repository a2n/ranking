# Search Engine Ranking
Query (🔎) ranking of the site with the keywords in specified search engine.

## Usage 🔧

### Google
```Go
NewGoogle().Get("gov.tw", "Taiwan", 5)
```
It shows the keyword **Taiwan** results, and find the site **gov.tw** ranking, up to **5 * 10** items result, in **google.com.tw**.


### Bing, aka Yahoo
```Go
NewBing("KEY").Get(zh_TW, "gov.tw", "Taiwan", 5)
```

It shows the keyword **Taiwan** results, and find the site **gov.tw** ranking, up to **5 * 10** item result, in **zh_TW** market.


## License
[MIT License](http://choosealicense.com/licenses/mit/)
