# whats-the-weather
Golang cli for telling you the weather

## Installation

I reccomend installing this using homebrew
```
brew tap ciaarraa/by-ciaarraa
brew install whats-the-weather
```

## Configuration
This CLI uses the geocoder api as a location service. This location service needs an api key. In order to use this cli, you will need to make a `GEOCODE_API_KEY` secret availaible to the application.

You can generate an api key at the following website: https://geocode.maps.co/

```
GEOCODE_API_KEY=<secret-value> whats-the-weather location Paris
```