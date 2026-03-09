version 1

service auth

route POST "/login"
    log "login attempt"
    respond 200 "success"
