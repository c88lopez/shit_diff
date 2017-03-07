config.json file

```
{
  "Login": {
    "Endpoint": "/login/endpoint",
    "Fields": {
      "Username": ["field.name", "field.value"],
      "Password": ["field.name", "field.value"]
    }
  },

  "Logout": {
    "Endpoint": "/logout/endpoint"
  },

  "Endpoints": "Endpoints.csv",
  "Results": "Results.csv",

  "Domains": [
    "production.domain.com",
    "stage.domain.com",
  ],

  "Timeout": 900
}
```
