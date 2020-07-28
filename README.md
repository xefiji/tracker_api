# Sigfox's callback endpoint

Just a basic api to expose an endpoint for Sigfox's backend datas.

## Server

Currently running on a raspi 3, monitored with supervisord.

This is a work in progress: it just logs coords and datas for the moment.

## Todo

- save an Mqtt job
- write consumer
- choose a database and save it in
- notify
- expose GET endpoint to populate an OSM map of all geo points

## Device

Datas come from a [Pytrack](https://docs.pycom.io/datasheets/boards/pytrack/) with a [Lopy4](https://docs.pycom.io/gettingstarted/connection/lopy4/) on top of it.

![Pytrack](https://pycom.io/wp-content/uploads/2020/03/Website-Product-Shots-Pytrack-LoPy4.png)

## Further

- write endpoints for LoraWan and the [Ttgo T-Beam](https://www.hackster.io/news/the-ttgo-t-beam-an-esp32-lora-board-d44b08f18628)
