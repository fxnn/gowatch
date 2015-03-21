# gowatch

gowatch will provide configurable logfile analysis for your server. It will be able to parse your logfiles and create
summaries in formats ready for delivery via E-Mail or Web.

However, this is still under development and _not_ ready for use yet.

## Related work

* **[logwatch](http://logwatch.sourceforge.net)** is widely used by Linux server administrators round the world, and so did
  I use it for many years. However, I find it to be not flexible enough in its configuration, and as soon as I want to
  change something, I always felt it was hard to extend and hard to change. Gowatch aims to be flexible, configurable
  and extendable.
* **[logstash](http://logstash.net)** is a log processor, that became very popular in combination with the search serer
  [elasticsearch](http://www.elasticsearch.org). Those are really great tools, especially for usage in large server
  parks. However, they need several Gigabytes of RAM and that's just far too heavy for my small tiny server. Gowatch
  aims to be a small and easy-to-be-used tool with low requirements, just as logwatch always was.

## License

Licensed under MIT, see for [LICENSE](LICENSE) file.
