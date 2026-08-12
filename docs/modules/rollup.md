# Rollup Module

`internal/rollup` validates interval names and produces aligned UTC buckets.

- Minutes: 60 minute buckets.
- Hours: 24 hour buckets.
- Days: 30 day buckets.
- Months: 12 month buckets.

Grid SQL averages successful values at each level before feeding the next level. A bucket is failed only when no successful value exists.
