# Eloi's Server

This is the server of [Eloi](https://github.com/mbaraa/eloi), where it scrapes over Zugania's website and provides a list of overlays to be used by the CLI client.

This server updates the Overlays and Ebuilds metadata everyday after midnight, so that the client can get the latest version of the packages as soon as possible.

Also you can use the endpoints if you need a proper JSON of the Ebuilds or Overlays.

- `/overlays/all` returns all of the overlays' data
- `/overlays/single?name=someOverlay` returns the requested overlay's data
- `/overlays/ebuilds` returns all ebuilds' data

Overlays and Ebuilds JSON schema can be found in the server's or client's [models](https://github.com/mbaraa/eloi-server/tree/main/models).
