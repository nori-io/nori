# Nori Engine

Nori is plugin engine for Golang.

All you need:
1. Implement Plugin interface and compile you project as golang plugin
2. Put your file into plugin directory
3. Start Nori

# Run Nori

```bash
./nori server --config="/etc/nori/config.json"
```

--config - path to your config file, by default Nori looking for config in ~/.config/nori

# Configuration

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

# Contributors

- Sergei Che
- Stan Shulga 
