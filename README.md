# current-ftp-csv-export

[![Build Status](https://travis-ci.org/sesam-community/current-ftp-csv-export.svg?branch=master)](https://travis-ci.org/sesam-community/current-ftp-csv-export)

Export av Current time data (csv file) from FTPS to Sesam

## System setup
```json
{
  "_id": "current-time-ftp-source",
  "type": "system:microservice",
  "docker": {
    "environment": {
      "FTP_PASSWORD": "<password>",
      "FTP_PORT": 21,
      "FTP_SERVER": "ftp.current.no",
      "FTP_USER": "<username>",
      "WS_PORT": 8080
    },
    "image": "ohuenno/sesam-ftps-csv-source",
    "port": 8080
  },
  "verify_ssl": true
}

```
