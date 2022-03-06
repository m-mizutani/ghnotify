package github.notify

notify[msg] {
    input.name == "check_suite"
    input.event.check_suite.conclusion == "failure"

    msg := {
        "text": "Check suite failed",
    }
}
