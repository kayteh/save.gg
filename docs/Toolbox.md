# SGG Toolbox

sgg-tools is a collection of tools to manage various things about the app and how it's running. It's an analogue to both `rake` from Ruby and `tinker` from Laravel Artisan.

Code: `cmd/sgg-tools`  
Usage: `sgg-tools COMMAND`  



## Tools


### `migrate`

Migrate controls database schema from the `migrations` folder. 

Usage: `sgg-tools migrate [[COMMAND] <args>]`

#### Commands

- `up`, also default when no command is specified. Rolls database schema forward.  
   In production, this is the only command that is permitted.

- `down`. Rolls back a schema change. Use it if your most recent schema change sucks.

- `redo`. Runs `down` then `up`. Use it when you've made a schema improvement and want to start it over.

- `create <NAME>`. Creates a new migration with the specified name. This command will ensure this tool can use the migrations. 


### `touch`

Invalidates a model in the cache. 

Usage: 

- `sgg-tools touch id <TYPE> <ID>` - By GUID

- `sgg-tools touch slug <TYPE> <SLUG>` - By URL identifier

- `sgg-tools touch url <URL>` - By canonical URL


### `debug-config`

Outputs the current configuration.

Usage: `sgg-tools debug-config`