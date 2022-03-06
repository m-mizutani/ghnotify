package github.notify

notify[msg] {
    input.name == "issues"
    input.event.action == "opened"
    input.event.repository.full_name == "mizutani-sandbox/test-repo"
    labels = { name | name := input.event.issue.labels[_].name }
    msg := {
        "channel": "#alert",
        "color": "#123456",
        "text": "issue opened",
        "body": input.event.issue.body,
        "fields": [
            {
                "name": "labels",
                "value": concat(", ", labels),
            },
        ],
    }
}
