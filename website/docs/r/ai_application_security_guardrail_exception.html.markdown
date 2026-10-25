---
subcategory: "AI Application Security"
layout: "incapsula"
page_title: "incapsula_ai_application_security_guardrail_exception"
description: |-
  Provides an Imperva AI Application Security guardrail exception resource.
---

# incapsula_ai_application_security_guardrail_exception

Provides an Imperva AI Application Security guardrail exception resource.

A guardrail exception is a text pattern that one guardrail of a policy lets through instead of blocking or
alerting on it - for example, test prompts that the `PROMPT_INJECTION` guardrail would otherwise flag. The guardrail
is chosen by its `guardrail_type` and `phase`; the provider looks it up in the policy.

## Example Usage

```hcl
resource "incapsula_ai_application_security_guardrail_exception" "qa_prompts" {
  application_id    = incapsula_ai_application_security_application.app.id
  policy_id         = incapsula_ai_application_security_policy.policy.id
  guardrail_type    = "PROMPT_INJECTION"
  phase             = "PROMPT"
  exception_pattern = "ignore previous instructions in test prompts"
  exception_reason  = "FALSE_POSITIVE"
  note              = "QA regression prompts"
}
```

## Argument Reference

* `account_id` - (Optional) Numeric ID of the account that owns the application. Defaults to the account of the API credentials. Changing it creates a new guardrail exception.
* `application_id` - (Required) UUID of the application the policy belongs to. Changing it creates a new guardrail exception.
* `policy_id` - (Required) UUID of the policy whose guardrail this exception applies to. Changing it creates a new guardrail exception.
* `guardrail_type` - (Required) Type of the policy's guardrail: `PROMPT_INJECTION`, `PII_STATIC`, `MODERATION`, `ZERO_SHOT_CLASSIFICATION` or `SYSTEM_PROMPT_LEAK`. `RATE_LIMIT` is not supported: it has no text to except.
* `phase` - (Required) Phase of that guardrail: `PROMPT` or `RESPONSE`. The policy must have exactly one guardrail of this type in this phase; otherwise the apply fails with `policy <id> has no <type> guardrail in phase <phase>` (or, if it has several, `... can't tell which one to use`).
* `exception_pattern` - (Required) The text the guardrail lets through, by default 3 to 200 words and at most 1500 characters. The service removes leading and trailing whitespace when it saves it, so such whitespace (e.g. a heredoc's trailing newline) is not a change.
* `exception_reason` - (Required) Why the exception exists: `BUSINESS_REQUIREMENT`, `FALSE_POSITIVE` or `OTHER` (exact, case-sensitive).
* `note` - (Optional) Note for the people managing the exception, up to 1000 characters. Not sent to the gateway. Leading and trailing whitespace is removed when saved.
* `active` - (Optional) Whether the exception is applied. Default: `true`.

Changing `guardrail_type`, `phase`, `exception_pattern`, `exception_reason`, `note` or `active` updates the guardrail
exception in place - it keeps its ID, also when it moves to another guardrail of the policy.

A policy can have at most 50 guardrail exceptions, enabled and disabled together. Creating one more fails with
`This policy already has 50 guardian exceptions, which is the maximum. Delete an exception to add a new one.`
Setting `active` on an existing exception is always allowed.

## Attributes Reference

* `id` - UUID of the guardrail exception.
* `guardrail_id` - ID of the guardrail the exception is attached to. If the policy's guardrails are recreated outside Terraform, the next plan shows a change and the apply attaches the exception to the guardrail of its type and phase again.
* `author` - The user who created the guardrail exception.
* `updated_by` - The user who last changed it; empty if it was never changed.

## Import

Import with `<application_id>/<policy_id>/<id>`, optionally prefixed with `<account_id>/`:

```
$ terraform import incapsula_ai_application_security_guardrail_exception.example 6f1c2d3e-0a4b-4c5d-8e9f-1a2b3c4d5e6f/7c9e6679-7425-40de-944b-e07fc1f90ae7/9b1d6e3a-4c2f-4b8e-a1d7-5f0e2c3b4a69
$ terraform import incapsula_ai_application_security_guardrail_exception.example 1234567/6f1c2d3e-0a4b-4c5d-8e9f-1a2b3c4d5e6f/7c9e6679-7425-40de-944b-e07fc1f90ae7/9b1d6e3a-4c2f-4b8e-a1d7-5f0e2c3b4a69
```

In the Imperva UI, the guardrail exception window, the application page and the policy page each have a button to copy their ID.
