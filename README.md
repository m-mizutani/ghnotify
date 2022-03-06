# ghnotify

`ghnotify` is a general GitHub event notification tool to Slack with [Open Policy Agent](https://github.com/open-policy-agent/opa) and [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/). There are a lot of notification tools from GitHub to Slack. However, in most case, notification rules are deeply integrated with source code implementation and customization of notification rule by end-user is limited.

`ghnotify` uses generic policy language Rego and OPA as runtime to separate implementation and policy completely. Therefore, `ghnotify` can handle and notify **all type of GitHub event**, not only issue/PR comments but also such as following events according to your Rego policy.

- GitHub Actions success/failure
- deploy key creation
- push to specified branch
- add/remove a label
- repository creation, archived, transferred
- team modification
- a new GitHub App installation

## Setup

### 1) Retrieve Bot User OAuth Token of Slack

Create your Slack bot and keep *OAuth Tokens for Your Workspace* in *OAuth & Permissions* page.

### 2) Creating Rego policy

`ghnotify` evaluates received GitHub event one by one. If `notify` variable exists in evaluation results, `ghnotify` notifies a message to Slack according to the results.

Policy rules are following.

**Input**: What data will be provided
- `input.name`: Event name. It comes from `X-GitHub-Event` header.
- `input.event`: Webhook events. See [docs](https://docs.github.com/en/developers/webhooks-and-events/webhooks) for more detail and schema.

**Result**: What data should be returned
- `notify`: Set of notification messages
    - `notify[_].channel`: Destination channel of Slack. It can be used by only API token
    - `notify[_].text`: Custom message of slack notification
    - `notify[_].body`: Custom message body
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
        "channel": "#notify-mizutani",
        "text": "Hello, mizutani",
        "body": input.event.comment.body,
    }
}
```

Then, you shall get a message like following.

![](https://user-images.githubusercontent.com/605953/155864886-c9c8ccbb-809c-44df-8925-fe69a0d820f4.png)


#### Example 2) Notification of workflow (actions) failed

```rego
package github.notify

notify[msg] {
    input.name == "workflow_run"
    input.event.action == "completed"
    input.event.conclusion == "failure"

    msg := {
        "channel": "#notify-failure",
        "text": "workflow failed",
        "color": "#E01E5A", # red
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
    labels := { name | name := input.event.pull_request.labels[_].name }

    msg := {
        "channel": "#notify-mizutani",
        "text": "breaking change assigned",
        "fields": [
            {
                "name": "All labels",
                "value": concat(", ", labels),
            },
        ],
    }
}
```

## Run

### Use Case 1: As GitHub Actions

- Pros: Easy to install
- Cons: GitHub Actions can receive events from only the repository

Create GitHub Actions workflow as following.

```yaml
name: Build and publish container image

on:
  push:
  issue:
  issue_comment:

jobs:
  build:
    runs-on: ubuntu-latest
    env:
      GHNOTIFY_SLACK_API_TOKEN: ${{ secrets.GHNOTIFY_SLACK_API_TOKEN }}
    steps:
      - name: checkout
        uses: actions/checkout@v2
      - name: dump event
        run: echo '${{ toJSON(github.event) }}' > /tmp/event.json
      - uses: docker://ghcr.io/m-mizutani/ghnotify:latest
        with:
          args: "emit -f /tmp/event.json -t ${{ github.event_name }} --local-policy ./policy"
```

### Use Case 2: As GitHub App server

- Pros: Easy to install
- Cons: GitHub Actions can receive event on each repository, can not watch organization wide

#### Deploy `ghnotify`

Deploy `ghnotify` to your environment and prepare URL that can be accessed from public internet. I recommend [Cloud Run](https://cloud.google.com/run) of Google Cloud in the use case.

When deploying `ghnotify`, I recommend to generate and use *Webhook secret* value. Please prepare random token and provide it to `--webhook-secret`.

Callback endpoint will be `http://{hostname}:4080/webhook/github`. You can change port number by `--addr` option.

#### Create a new GitHub App

1. Go to https://github.com/settings/apps and click `New GitHub App`
2. Grant permissions and check events you want to subscribe in `Subscribe to events`.
3. Check `Active` in `Webhook` section
4. Set URL of deployed `ghnotify` to `Webhook URL`
5. Set *Webhook secret* to `Webhook secret` if you configured
6. Then click `Create GitHub App`

#### Example

## Options

- Server
    - `--addr`: Server address and port to listen webhook. e.g. `0.0.0.0:8080`
    - `--webhook-secret`: Webhook secret
- Policy (either one of `--local-policy` and `--remote-url` is required)
    - `--local-policy`: Policy files or directory.
    - `--local-package`: Package name of policy file
    - `--remote-url`: URL of OPA server
    - `--remote-header`: HTTP header to query OPA server
- Notification (either one of following is required)
    - `--slack-api-token`: API token retrieved in Step 1 (Recommended)
    - `--slack-webhook`: Incoming webhook URL of Slack

## License

Apache License 2.0
