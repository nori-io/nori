# Nori docker image

### How to build

To build docker image, run `make docker-build` command.

### Structure
```bash
/app
/app/nori [nori binary]
/app/config/nori.yml [nori config file]
/plugins [mount plugins]
/configs [mount configs]
```

### Env variables

Variable | Description | Default | Examples
---|---|---|---
NORI_CONFIG | Set of paths to config files | | "/configs/*.yaml", "/configs/nori.yaml,/configs/plugins.yaml"
NORI_PLUGINS_DIR | Set of paths to plugin dirs | /plugins/ | "/nori/plugins", "/plugins/web,/plugins/admin"