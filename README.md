# Nori CMS

## Configuration

- plugin manager
  - pm.source = sql
  - pm.sql.driver = mysql
  - pm.sql.connection = 
- plugins
  - plugins.dir : string - path to directory with plugins
- routing
- database

## Plugins

Plugins uploaded into plugin directory. Plugin directory defined in configuration.

#### Build Plugin

```bash
go build -buildmode=plugin -o plugin.so
```