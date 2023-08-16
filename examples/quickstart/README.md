# Quickstart

For a quick verification of `server` subcommand of `ghnotify` in a local environment, this document may be helpful.

Note: in this quickstart, you have to POST event payload manually, as opposed to using actual GitHub event.

## Prerequisits

- [Docker](https://docs.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Preparations

### Setting Up Slack Channel and Policy File

1. Create Slack channel dedicated to receiving notification from `ghnotify`
1. Draft a policy `.rego` file under `examples/policy/`
    1. You can achive this easily by duplicating any existing `.rego` files in `examples/policy/`
1. Ensure the name (not the ID) of the Slack channel created in step `1.` is configured as the `"channel"` parameter in the `.rego` file
1. Generate the bundled policy file, `bundle.tar.gz`, using the command `make bundle`

### Configuring Sensitive Values

Start by duplicating `.env.template` to `.env` to store the sensitive values generated in this section.

#### Slack

1. Create Slack App capable of posting messages to the designated channel
1. Populate the `GHNOTIFY_SLACK_API_TOKEN` variable in `.env` with the bot user's OAuth token
1. Install the Slack App to the channel dedicatd to receiving notification

#### GitHub

1. Prepare your GitHub webhook secret
    1. Generate a high-entropy random string to be used as GitHub Webhook secret
        - cf: https://docs.github.com/webhooks-and-events/webhooks/securing-your-webhooks#setting-your-secret-token
    1. Populate the `GHNOTIFY_WEBHOOK_SECRET` variable in `.env` with the generated secret
1. Setup GitHub App
    1. Instantiate a new GitHub App, determining its affiliation (be it personal or organizational) based on your specific requirements
    1. Assign permission to the app based on the events you're interested in monitoring
    1. In the GitHub App's settings, enter the webhook secret from step `1.i.` into the "Webhook secret (optional)" field
    1. Finalize the app setup and install it to your user profile or organization as necessary

### Crafting the Webhook Event Payload

`ghnotify` authenticates the webhook event payload POSTed to it using HMAC-sha256 hash method, consistent with GitHub's procedure.
Thus for this quickstart, when POSTing the event payload, you must also send the hased value.

You can achive the step explained below by using Postman collection [GitHub Webhook Events](https://web.postman.co/workspace/201d9100-8e44-4ea1-8028-db63b9593ad1/collection/27582439-e25dc84e-e2a3-4a8d-bf6b-d1ed02e91a46 )

#### JSON Event Payload

Draw from [GitHub's official documentation](https://docs.github.com/en/webhooks-and-events/webhooks/webhook-events-and-payloads ) to craft your event payload manually.

For example, the structure of a payload of an event triggered by labelling a pull request as `breaking-change` is `{"action": "labeled", "label": {"name": "breaking-change"}}` with header `X-GitHub-Event: pull_request`.

#### HMAC-sha256 Hashe

Use the webhook secret from earlier to hash the event payload with HMAC-sha256.
Given the secret `c15262968fc607f44b4009c7df4dac623395c6e0`, the hashing process would look like:
```bash
$ echo -n '{action": "labeled", "label": {"name": "breaking-change"}}' | openssl dgst -sha256 -hmac "c15262968fc607f44b4009c7df4dac623395c6e0"
```

This yields the hash `558fda5efec10bafa52d14a6c12ccad6ed512826928a1493333f488f5787a124`.


## Deploy and POST event

1. Deploy `ghnotify` notification system using `docker-compose up`
1. POST your event to the `ghnotify` endpoint

For example:
```bash
$ curl -X POST \
  -d '{"action": "labeled", "label": {"name": "breaking-change"}}' \
  -H 'Content-Type: application/json' \
  -H 'X-GitHub-Event: pull_request' \
  -H 'x-hub-signature-256: sha256=558fda5efec10bafa52d14a6c12ccad6ed512826928a1493333f488f5787a124' \
  localhost:4080/webhook/github
```

Happy notifying! :rocket:
