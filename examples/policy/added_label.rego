package github.notify

notify[msg] {
    input.name == "pull_request"
    input.event.action == "labeled"
    input.event.label.name == "breaking-change"

    msg := {
        "channel": "#alert",
        "text": "A new breaking change PR",
    }
}
