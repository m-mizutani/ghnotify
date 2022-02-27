# ghnotify

General GitHub event notification tool to Slack with [Open Policy Agent](https://github.com/open-policy-agent/opa) and [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/). There are a lot of notification tools from GitHub to Slack. However, in most case, notification rules are deeply integrated with source code implementation and customization by end-user is limited. `ghnotify` uses generic policy language Rego and OPA as runtime to separate implementation and policy completely.

## Usage

### 1) Retrieve incoming webhook URL of Slack

Setup incoming webhook according to https://api.slack.com/messaging/webhooks and note *incoming webhook URL*.

### 2) Creating Rego policy

Policy rules are following.

**Input**: What data will be provided
- `input.name`: Event name. It comes from `X-GitHub-Event` header.
- `input.event`: Webhook events. See [docs](https://docs.github.com/en/developers/webhooks-and-events/webhooks) for more detail and schema.

**Result**: What data should be returned
- `notify`: Set of notification messages
    - `notify[_].channel`: Destination channel of Slack. If empty, notification will be sent to default channel of the incoming webhook.
    - `notify[_].text`: Custom message of slack notification
    - `notify[_].header`: Custom message header (like title)
    - `notify[_].color`: Message bar color
    - `notify[_].fields`: Set of custom message fields.
        - `notify[_].fields[_].name`: Field name
        - `notify[_].fields[_].value`: Field value
        - `notify[_].fields[_].url`: Link assigned to the field

#### Example 1) Notification of "call me" in issue comment

```rego
package github.notify

notify[msg] {
    input.name == "issue_comment"
    contains(input.event.comment.body, "mizutani")
    msg := {
        "channel": "notify-mizutani",
        "text": "Hello, mizutani",
    }
}
```

#### Example 2) Notification of workflow (actions) failed

```rego
package github.notify

notify[msg] {
    input.name == "workflow_run"
    input.event.action == "completed"
    input.event.conclusion == "failure"

    msg := {
        "channel": "notify-failure",
        "text": "workflow failed",
        "color": "#E01E5A", # red
        "fields": [
            "name": "Repository",
            "value": input.event.full_name,
        ],
    }
}
```

#### Example 3) Assigned "breaking-change" label to PR

```rego
package github.notify

notify[msg] {
    input.name == "pull_request"
    input.event.action == "labeled"
    input.event.label.name == "breaking-change"

    msg := {
        "channel": "notify-mizutani",
        "text": "breaking change assigned",
        "fields": [
            "name": "Repository",
            "value": input.event.full_name,
        ],
        "fields": [
            "name": "PR",
            "value": input.event.pull_request.title,
            "url": input.event.pull_request.html_url,
        ],
    }
}
```

### 3) Deploy `ghnotify`

Deploy `ghnotify` to your environment and prepare URL that can be accessed from public internet. I recommend [Cloud Run](https://cloud.google.com/run) of Google Cloud in the use case.

When deploying `ghnotify`, I recommend to generate and use *Webhook secret* value. Please prepare random token and provide it to `--webhook-secret`.

### 4) Create a new GitHub App

1. Go to https://github.com/settings/apps and click `New GitHub App`
2. Grant permissions and check events you want to subscribe in `Subscribe to events`.
3. Check `Active` in `Webhook` section
4. Set URL of deployed `ghnotify` to `Webhook URL`
5. Set *Webhook secret* to `Webhook secret` if you configured
6. Then click `Create GitHub App`

## Options

- Server
    - `--addr`: Server address and port to listen webhook. e.g. `0.0.0.0:8080`
    - `--webhook-secret`: Webhook secret
- Policy (either one of `--local-policy` and `--remote-url` is required)
    - `--local-policy`: Policy files or directory.
    - `--local-package`: Package name of policy file
    - `--remote-url`: URL of OPA server
    - `--remote-header`: HTTP header to query OPA server
- Notification
    - `--slack-webhook`: Incoming webhook URL of Slack

## License

Apache License 2.0
