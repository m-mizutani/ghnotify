package github.notify

notify[msg] {
    input.name == "repository"
    input.event.action == "created"

    msg := {
        "channel": "#alert",
        "text": "repository created",
    }
}

notify[msg] {
    input.name == "repository"
    input.event.action == "deleted"

    msg := {
        "channel": "#alert",
        "text": "repository created",
    }
}

notify[msg] {
    input.name == "repository"
    input.event.action == "publicized"

    msg := {
        "channel": "#alert",
        "text": "repository publicized",
    }
}
