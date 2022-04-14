package github.notify

notify[msg] {
    input.name == "issue_comment"
    contains(input.event.comment.body, "mizutani")
    msg := {
        "channel": "#notify-mizutani",
        "text": "Hello, mizutani",
        "body": input.event.comment.body,
    }
}
