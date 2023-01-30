# Jitsi Feedback Plugin

## Installation

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
    // ,'USER_AGENT'
    // ,'USER_REGION'
];


// Optional.
config.deploymentInfo = {
         shard: "shard1",
         region: "europe",
         userRegion: "asia",
         envType: "envType",
         backendRelease: "backendRelease",
         environment: "production"
};

```

## Sample Data

| Metadata            | Value                                 |
| ------------------- | ------------------------------------- |
| APP_BACKEND_RELEASE | "backendRelease"                      |
| APP_ENV_TYPE        | "envType"                             |
| APP_ENVIRONMENT     | "production"                          |
| APP_FOCUS_VERSION   | "1.0.954"                             |
| APP_LIB_VERSION     | "{#COMMIT_HASH#}"                     |
| APP_MEETING_REGION  | "europe"                              |
| APP_NAME            | "Jitsi Meet"                          |
| APP_REGION          | "europe"                              |
| APP_SHARD           | "shard1"                              |
| BROWSER_NAME        | "firefox"                             |
| BROWSER_VERSION     | "109.0"                               |
| DISPLAY_NAME        | "ElementUser"                         |
| EXTERNAL_API        | true                                  |
| IN_IFRAME           | true                                  |
| MATRIX_USER_ID      | "@ElementUser:localhost"              |
| MEETING_ID          | "EFKVM...XXG5A"                       |
| MEETING_URL         | "https://localhost:8443/EFKVM...",    |
| OS_NAME             | "Windows"                             |
| OS_VERSION_NAME     | "10"                                  |
| OS_VERSION          | "NT 10.0"                             |
| USER_AGENT          | "Mozilla/5.0 AppleWebKit/537.36 ...." |
| USER_REGION         | "asia"                                |

## Compatibility

| Metadata            | Firefox | Safari | Chrome | Chromium | Edge |
| ------------------- | ------- | ------ | ------ | -------- | ---- |
| APP_BACKEND_RELEASE | :o:     |        | :o:    | :o:      | :o:  |
| APP_ENV_TYPE        | :o:     |        | :o:    | :o:      | :o:  |
| APP_ENVIRONMENT     | :o:     |        | :o:    | :o:      | :o:  |
| APP_FOCUS_VERSION   | :o:     |        | :o:    | :o:      | :o:  |
| APP_LIB_VERSION     | :o:     |        | :o:    | :o:      | :o:  |
| APP_MEETING_REGION  | :o:     |        | :o:    | :o:      | :o:  |
| APP_NAME            | :o:     |        | :o:    | :o:      | :o:  |
| APP_REGION          | :o:     |        | :o:    | :o:      | :o:  |
| APP_SHARD           | :o:     |        | :o:    | :o:      | :o:  |
| BROWSER_NAME        | :o:     |        | :o:    | :o:      | :o:  |
| BROWSER_VERSION     | :o:     |        | :o:    | :o:      | :o:  |
| DISPLAY_NAME        | :o:     |        | :o:    | :o:      | :o:  |
| EXTERNAL_API        | :o:     |        | :o:    | :o:      | :o:  |
| IN_IFRAME           | :o:     |        | :o:    | :o:      | :o:  |
| MATRIX_USER_ID      | :o:     |        | :o:    | :o:      | :o:  |
| MEETING_ID          | :o:     |        | :o:    | :o:      | :o:  |
| MEETING_URL         | :o:     |        | :o:    | :o:      | :o:  |
| OS_NAME             | :o:     |        | :o:    | :o:      | :o:  |
| OS_VERSION_NAME     | :o:     |        | :o:    | :o:      | :o:  |
| OS_VERSION          | :o:     |        | :o:    | :o:      | :o:  |
| USER_AGENT          | :o:     |        | :o:    | :o:      | :o:  |
| USER_REGION         | :o:     |        | :o:    | :o:      | :o:  |

* :o: supported
* :x: problems reported

Tested on 
Jitsi stable-8044-1
Firefox 109.0 (windows 11)
Chrome 108 (Linux)

#### Notes

`APP_FOCUS_VERSION`, `PARTICIPANT_ID`  not available when the feedback is submitted before leaving the call, but are present if the feedback is sent using the "Leave Feedback" button in the toolbar.

`OS_VERSION_NAME` and `OS_VERSION` may be empty on Linux hosts.

`PARTICIPANT_ID` is  the local user's ID in the jitsi room. Lookup `JitsiConference.prototype.myUserId` function  from `lib-jitsi-meet` for reference.```
