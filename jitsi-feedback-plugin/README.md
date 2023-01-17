### How to install the custom analytics feedback handler


We suggest to mount the `feedback.js` to the root of the jitsi-meet application, matching the path set in `config.analytics.scriptURLs`.
For example, the location inside the jitsi-meet web container is `/usr/share/jitsi-meet/feedback.js`, and it will be reachable from clients at `/feedback.js`.

If required by your special use case, it is technically possible to host `feedback.js` anywhere you want as long as clients can reach it. Simply adjust the path given in the `scriptURLs` config setting needs to match.

### Required Jitsi configuration

These are the relevant settings that need to be set in jitsi-meet.
The intended way to configure this is by leveraging the `custom-config.js` method.

```javascript
// address of the feedback backend REST API, reachable from the end user device
config.feedbackBackend = 'https://example.org:8080'

// percentage of users to automatically request feedback from when leaving the call
// it's 100 by default if undefined, i.e. always shown
config.feedbackPercentage = 100;

// enables the feedback button in the toolbar, can be any string
config.callStatsID = 'id';

// disables CallStats even if callStatsID is set
config.callStatsSecret = null;

// this will enable analytics, both need to be false for our handler to work
config.disableThirdPartyRequests = false;
config.analytics.disabled = false

// custom analytics handler loads our feedback.js plugin
config.analytics.scriptURLs = ['/feedback.js'];

// Array<string> of enabled metrics and metadata items - uncomment to enable
// some of the metadata may not be available depending on the user's browser and device as well as the configuration of the jitsi backend
config.metadata = [ ''
    // ,'APP_BACKEND_RELEASE'
    // ,'APP_ENV_TYPE'
    // ,'APP_ENVIRONMENT'
    // ,'APP_FOCUS_VERSION'
    // ,'APP_LIB_VERSION'
    // ,'APP_MEETING_REGION'
    // ,'APP_NAME'
    // ,'APP_REGION'
    // ,'APP_SHARD'
    // ,'BROWSER_NAME'
    // ,'BROWSER_VERSION'
    // ,'DISPLAY_NAME'
    // ,'EXTERNAL_API'
    // ,'IN_IFRAME'
    // ,'MATRIX_USER_ID'
    // ,'MEETING_ID'
    // ,'MEETING_URL'
    // ,'OS_NAME'
    // ,'OS_VERSION_NAME'
    // ,'OS_VERSION'
    // ,'PARTICIPANT_ID'
    // ,'USER_AGENT'
    // ,'USER_REGION'
];

```
