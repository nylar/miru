![miru](http://i.imgur.com/0ssoHoU.png)

[![Build Status](https://semaphoreci.com/api/v1/projects/3ebec8f7-d164-4823-b23c-665f92d8a7da/374112/badge.png)](https://semaphoreci.com/nylar/miru)
[![Coverage Status](https://coveralls.io/repos/nylar/miru/badge.svg?branch=HEAD)](https://coveralls.io/r/nylar/miru?branch=HEAD)
[![license](http://img.shields.io/badge/license-unlicense-blue.svg "license")](https://raw.githubusercontent.com/nylar/miru/master/UNLICENSE)

## API

### Queues

```
/api/queues/
```

Returns a list of queues.

### Queue

```
/api/queue/bbc.co.uk
```

Return an individual queue.

### Sites

```
/api/sites
```

Returns a list of sites.

### Crawl

```
/api/crawl?url=http%3A%2F%2Fbbc.co.uk%2F
```

Crawls a given URL, will then recursively crawl each found link until the queue list is exhausted.

### Search

```
/api/search?q=news
```

Searches the datastore for any pages with an index matching the keywords.
