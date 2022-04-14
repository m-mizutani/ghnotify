package github.notify

notify[msg] {
    input.name == "deploy_key"
    input.event.action == "created"
    msg := {
        "channel": "#alert",
        "text": "A new deploy key created",
        "fields": [
            {
                "name": "title",
                "value": input.event.key.title,
            },
            {
                "name": "read_only",
                "value": input.event.key.read_only,
            },
        ],
    }
}
