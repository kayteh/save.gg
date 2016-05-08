# Infrastructure

## Backend Data Stores

- **Primary**: Postgres (XC in production)
  * Operates as the single source of truth, period.

- **Secondary/NoSQL**: RethinkDB
  * All unstructured data.

- **Caching**: Redis
  * Write-through LRU + timeout cache.

- **Metrics**: InfluxDB
  * Measure fucking everything.

- **Config**: etcd
  * Store as many config options as possible here, but in dev use flat files.
	
- **Logging**: Graylog
  * Log fucking everything.

- **Cluster Comms**: NATS
  * Ties everything together.

## Queues

sgg-workers will subscribe to all priorities (high, medium, low) prefixed by `work-queue:priority`. They will process jobs in order of priority. Ideally, low priority jobs can be done on a slower server with more cores, but medium jobs will be done after high queues on faster CPU systems. 

- high priority jobs are anything that needs to be in real-time.

- medium priority jobs are default, and are anything that needs to be done behind the scenes. (e.g. crawling of other services, sending emails, tweeting/etc)

- high prority jobs are used for non-realtime things, and should definitely be used for any gamesave manipulation and data mining. 
