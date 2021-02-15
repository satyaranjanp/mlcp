# Multi Level Car Parking

## Design Flow
1. MLCP implements a cache to keep track of available slots in each level (floor)
2. MLCP receives request to get nearest parking slot (via api or MQ)
3. It checks the cache and returns the nearest available slot (lowes probable floor, lowest probable slot number)
4. Cache is maintained to keep application running in case DB is not reachable
4. It receives request to assign a slot to a vehicle
5. It updates cache (removing the free slot form cache) and then make request to write to DB
6. Writing to DB is buffered, so that any DB related issue can be gracefully handled
8. When any vehicle leaves a slot, cache is updated and a write request is made to put the data in DB, like slot details and vehicle details
## Features
1. Can be run in single intance mode (with internal cache) and multiple instance mode (with external cache, like redis, memcache)
1. Can be extended to use different DB
2. Can be extended to use different MQ
