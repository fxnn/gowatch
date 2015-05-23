# gowatch

gowatch will provide configurable logfile analysis for your server. It will be able to parse your logfiles and create
summaries in formats ready for delivery via E-Mail or Web.

However, this is still under development and _not_ ready for use yet.

[![Build Status](https://travis-ci.org/fxnn/gowatch.svg)](https://travis-ci.org/fxnn/gowatch)

## Installation

Just `go get` the program.
The following assumes that your `GOPATH` is set and your `PATH` contains your `$GOPATH/bin` directory;
if not so, please read the [Code Organization](https://golang.org/doc/code.html#Organization) chapter of the go manual.

```
$ go get github.com/fxnn/gowatch
$ gowatch
2015/04/08 19:10:44 No configuration file given. Specify one using `-c /path/to/config.yml`
```

## Usage

`gowatch` will always run with a configuration file, pass it with the `-c` option:

```
gowatch -c /path/to/config.yml
```

Relative paths will always be resolved based on your current working directory. Note, that this also holds for paths inside the configuration file.

The configuration files itself are separated into three main sections: logfiles, mappers *(not implemented yet)* and summarizers. This reflects the architecture (see below).

An example configuration file would be

```
logfiles:

- filename: /var/log/auth.log
  config: {pattern: '%{SYSLOGBASE} %{GREEDYDATA:Message}'}
  where: {
    timestamp: {"younger than": "24h"}
  }

- filename: /var/log/mail.log
  tags: ['mail.log']
  config: {pattern: '%{SYSLOGBASE} %{GREEDYDATA:Message}'}
  where: {
    timestamp: {"younger than": "24h"}
  }
  
summary:
  
  - do: count
    title: auth.log
    where: {tags: {contains: 'auth.log'}}
    config: {
      '%{pam_caller} [%{user}] %{pam_session_state}': '%{WORD:pam_module}\(%{DATA:pam_caller}(?::%{WORD:pam_facility})\): session %{WORD:pam_session_state} for user %{USERNAME:user}(?: by %{GREEDYDATA:pam_by})?',
      'sudo [%{user}->%{effective_user}] %{command}': '\s*%{USER:user}\s*: TTY=%{DATA:tty} ; PWD=%{PATH:pwd} ; USER=%{USER:effective_user} ; COMMAND=%{PATH:command}(: %{GREEDYDATA:arguments})?'
    }
  
  - do: count
    title: Stored Mails
    where: {tags: {contains: 'mail.log'}}
    config: {
      'Discarded': 'deliver\(%{USER:user}\): sieve: msgid=<%{DATA:msgid}>: marked message to be discarded if not explicitly delivered',
      'Stored [%{mailboxname}]': "deliver\\(%{USER:user}\\): sieve: msgid=<%{DATA:msgid}>: stored mail into mailbox '%{DATA:mailboxname}'",
    }
```


## Architecture

The core of `gowatch` is the following pipeline.

```
  +------------+     +------------+     +------------+
  +   Parser   | --> |   Mapper   | --> | Summarizer |
  +------------+     +------------+     +------------+
```

While each `parser.Parser` creates `logentry.LogEntry` instances (by parsing logfiles) and sends them into the pipeline, the `mapper.Mapper`s will modify these log entries and pass them to the summarizers. Each `summary.Summarizer` produces human readable output, e.g. by counting occurences or listing search results. The concatenation of output might then be given to the user, e.g. by mail.

The names are more specific than what [Logstash](http://logstash.net) uses, and this is by intention. The aim was to build an application specifically for creating reports from logfiles. Further usecases, like network support etc., are out of scope.


## Related work

* **[logwatch](http://logwatch.sourceforge.net)** is widely used by Linux server administrators round the world, and so
  did I use it for many years. However, I find it to be not flexible enough in its configuration, and as soon as I want
  to change something, I always felt it was hard to extend and hard to change. Gowatch aims to be flexible, configurable
  and extendable.
* **[logstash](http://logstash.net)** is a log processor, that became very popular in combination with the search serer
  [elasticsearch](http://www.elasticsearch.org). Those are really great tools, especially for usage in large server
  parks. However, they need several Gigabytes of RAM and that's just far too heavy for my small tiny server. Gowatch
  aims to be a small and easy-to-be-used tool with low requirements, just as logwatch always was.

## 3rd Party Libraries

[Standing on the shoulders of giants](http://en.wikipedia.org/wiki/Standing_on_the_shoulders_of_giants), this wouldn't
be what it is without:

* **[gemsi/grok](http://github.com/gemsi/grok)** is a great Grok implementation in Go, throughoutly tested.
  Grok itself is a simple DRY method for log parsing, known from
  [logstash](http://logstash.net/docs/latest/filters/grok), but there is also a standalone C implementation -- see for
  [jordansissel/grok](https://github.com/jordansissel/grok).
* **[stretchr/testify](http://github.com/stretchr/testify)** brings assertions to Go, just the way they feel right.
  Great for testing!
* **[go-yaml/yaml](https://github.com/go-yaml/yaml)** (un)marshalls YAML files into native Go data structures with few
  more than a single line of code. gowatch wouldn't have the configuration files it has without this library.

...among others. Thanks a lot for your work!

## License

Licensed under MIT, see for [LICENSE](LICENSE) file.
