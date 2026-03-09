version 1

service user

route GET "/profile"
    log "profile request"
    respond 200 "ok"

route POST "/create"
    log "creating user"
    respond 201 "created user successfully"
