# Jitsi Feedback Plugin

This plugin enables participants of a [Jitsi Meet](https://jitsi.org/jitsi-meet/) call to leave feedback to their service provider.
Upon leaving the call, the user is presented a form where they can give a rating on a 5-star scale and additionally may leave a comment.

This implementation is tailored to work when Jitsi is used as a video/VoIP provider as a widget within the [Element](https://element.io/) chat application.
In particular, it implements an authorization mechanism in conjunction with the backend, which ensures that feedback can only be submitted by authorized users.

The plugin may be configured through the normal Jitsi configuration file by the service provider to also include some technical metadata in the call.

## Key Features

- This application is a Jitsi plugin implementing a custom way to collect user feedback and call metrics.
    - The normal Jitsi feedback form is reused.
    - The address of the custom [backend](../backend) is configurable.
    - The collected metadata is configurable.
- The feedback form appears automatically at the end of a call (configurable).
    - Feedback can also be given during a call, if that option is enabled.
- It uses Jitsi's ability to load custom code by [including it in the index.html](https://github.com/jitsi/jitsi-meet/blob/3081b41d0d6f5b13ffddadbda1460f3548ceefbf/index.html#L191). This approach alleviates the maintainance workload as we don't need to merge code into upstream.

## How To Use

- Copy, mount, or put otherwise the `plugin.head.html` file into the right place: it should be placed next to `index.html`. For example, using the official docker image or .deb package, it should be available at `/usr/share/jitsi-meet/plugin.head.html`.
- Configure settings: use the regular Jitsi configuration file and modify/add the plugin-specific settings ([see below](#configuration)). The `custom-config.js` file provided by this repository may be used as a starting point.
- Run/restart Jitsi Meet and the feedback [backend](../backend).

## Configuration

You will need to set the following variables in the Jitsi configuration.
The Jitsi Docker container provides a way to adjust config using a [custom-config.js](https://jitsi.github.io/handbook/docs/devops-guide/devops-guide-docker#jitsi-meet-configuration) file.
<div style="margin-left: auto;
            margin-right: auto;
            width: 70%">

| Setting name                | Description                                                                      | Example     |
|-----------------------------|----------------------------------------------------------------------------------|-------------|
| `config.feedbackPercentage` | percentage of users to automatically request feedback from when leaving the call | `100`       |
| `config.callStatsID`        | optional: enable button in toolbar                                               | `'enabled'` |
| `config.metadata.???`       | optional: configure what metadata should be gathered                             | ``          |

</div>
