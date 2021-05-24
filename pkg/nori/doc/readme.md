### Plugin states

1. Nil - plugin was added in registry
2. Inited 
   a) plugin.Init func successfully executed, or
   b) plugin.Stop func successfully executed
3. Running - plugin.Start func successfully executed

Possible transitions:
- Nil -> Inited 
- Inited -> Running 
- Running -> Inited

#### None

None is a basic plugin state. It represents state of new imported plugin.

#### Inited

Inited represents the state of plugin ready to be started.
a) state of initialised plugin,  
b) state of stopped plugin.

#### Running

Running represents the state of started plugin i.e. plugin.Start have been successfully called.

### Plugin function calls

There three main plugin function that may change the state of plugin.
1. Init()
2. Start()
3. Stop()

 