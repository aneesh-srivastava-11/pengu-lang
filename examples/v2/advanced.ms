version 2
service advanced_v2

// Service-level configs
middleware logging
middleware auth jwt
health enable
metrics enable

// Basic route from V1
route GET "/profile"
    log "fetching profile..."
    respond 200 "admin profile data"

// New JSON parsing route from V2
route POST "/user"
    parse json UserRequest
    auth jwt
    log "creating user..."
    respond 201 "created"
