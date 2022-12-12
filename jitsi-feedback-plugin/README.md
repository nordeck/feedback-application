### How to install the custom analytics feedback handler


Mount the feedback.js to the root of the jitsi-meet. The location inside the jitsi-meet web container `/usr/share/jitsi-meet/feedback.js`.


The corresponding settings from the jitsi custom-config.js
```javascript
// address of the feedback backend REST API
config.feedbackBackend = 'http://localhost:8333'

// percentage of users to automatically request feedback from when leaving the call
// it's 100 by default if undefined
config.feedbackPercentage = 100;

// enables the feedback button in the toolbar
config.callStatsID = 'id';

// disables CallStats even if callStatsID is set
config.callStatsSecret = null;

// this will enable analytics, both need to be false for our handler to work
config.disableThirdPartyRequests = false;
config.analytics.disabled = false

// custom analytics handler
config.analytics.scriptURLs = ['/feedback.js'];

// enable to collect metadata
// config.metadata.??? = true;
```