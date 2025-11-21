## v0.5.0 (2025-11-21)

### Feat

- websocket now support keep alive and welcome
- refactor bootstrap and log levels

### Fix

- main bootstrap passes server sk to dependencies
- timestamp inconsistency
- syscheck will check golang 1.24 or higher

### Refactor

- removed component registry
- revert config from user space to etc
- keys creation simplified
- keys generation simplification
- remove tls support

## v0.4.1 (2025-08-19)

### Fix

- ci release commit message

## v0.4.0 (2025-08-15)

### Feat

- add sub-apps support
- update css
- add app launcher mounting and actions
- add app table, model and API

### Fix

- typo in constant name

### Refactor

- wasm file loading process
- add fontawesome offline
- css stylesheet.
- overlay data source passing
- template loading process
- token selection

## v0.3.0 (2025-06-06)

### Feat

- component integration
- add middleware support
- add crypto message metadata property
- moved config in user space
- add port flag to start command

### Fix

- auth process

### Refactor

- wasm engine

## v0.2.1 (2025-03-10)

### Feat

- add machine informations
- add gh workflows
- refactored wasm framework
- add finial component interface
- add brotli compression
- add brotli compression
- add basic job runner
- ws support
- move proto in separate package
- add ws support
- add stringer pkg
- improved console display
- add login and dashboard landing page
- add base models and 1st user creation on app init
- add sqlite library
- add docker integration
- add orbital certificate gen.
- add async support and api call
- add wasm and wasm loader
- add dom control and ui events
- add scss working wireframe
- add DOM controller and bootstrap
- add init command
- add base

### Fix

- cz file
- ci branch rules
- add missing cz file
- state reject pointers
- ws connection duplicate
- init command. cleanup unneeded code

### Refactor

- components nesting and work wire framing
- state fix, template fix, dom enhance
- moving css classes to file
- state and event api streamline
- moved away from scss to pure css
- state and events mangers
- localstorage pkg
- dom manipulation pkg
- cleanup
- id to atomic uint
- ws backend
- update init and start commands
- add back web storage ignored by mistake
- add orbital root dir retrieve helper
- **cmd**: allow passing dependencies
- rename dashboard to web
- remove generated files
