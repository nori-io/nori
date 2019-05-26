# Nori Engine

Nori is plugin engine for Golang.

All you need:
1. Implement Plugin interface and compile you project as golang plugin
2. Put your file into plugin directory
3. Start Nori

## Run Nori

```bash
./nori server --config="/etc/nori/config.json"
```

--config - path to your config file, by default Nori looking for config in ~/.config/nori

## Configuration

- plugins.dir []string - paths to plugins dirs

Example config file:
```json
{
  "nori": {
    "storage": {
      "type": "none"
    }
  },
  "plugins": {
    "dir": [
        "/home/nori/.config/nori/plugins"
    ]
  },
  "http": {
    "addr": "localhost:8089"
  }
}
```


#### MySQL Storage

```json
{
  "nori": {
    "storage": {
      "type": "mysql",
      "source": "nori:nori@/noridb"
    }
  }
}
```

## Semantic Versioning

This repo uses Semantic versioning (http://semver.org/), so

MAJOR version when you make incompatible API changes,
MINOR version when you add functionality in a backwards-compatible manner, and
PATCH version when you make backwards-compatible bug fixes.

## Contributors

- Sergei Che
- Stan Shulga 
- Anita Nabieva
