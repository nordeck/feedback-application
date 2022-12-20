/*
 *  Copyright 2022 Nordeck IT + Consulting GmbH
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and  limitations
 *  under the License.
 *
 */

(function() {
    console.log("FEEDBACK loaded custom feedback");

    const LOG = 'FEEDBACK_HANDLER';
    
    class Feedback {
        constructor(options) {
            this.sendEvent = this.sendEvent.bind(this);
            this.setUserProperties = this.setUserProperties.bind(this);
        }

        // called from lib-jitsi-meet
        sendEvent(event) {
            // console.log(`${LOG} handler`, event);

            if (event.action === 'connection.stage.reached' && event.actionSubject === 'conference_muc.joined') {
                this.handleJoin();
                return;
            }

            if (event.action === 'feedback' && event.actionSubject === 'feedback') {
                this.handleFeedback(event.attributes);
                return;
            }
        }

        // called from lib-jitsi-meet
        setUserProperties(permanentProperties) {
            // nothing
        }

        handleFeedback(data) {
            const {rating, comment} = data;
            console.log(`${LOG} Feedback`, rating, comment);

            const jwt = window.APP.conference.feedbackToken;
            console.log(`${LOG} gather metrics for  JWT: ${jwt}`);
        
            const postFeedback = async (jwt, payload) => {
                const baseUrl = APP.store.getState()['features/base/config'].feedbackBackend;
                const url = `${baseUrl}/feedback`;
                
                const headers = {
                    'authorization': `Bearer ${jwt}`
                };
    
                const res = await fetch(url, {
                    method: 'POST',
                    headers,
                    body: JSON.stringify(payload)
                });
    
                if (!res.ok) {
                    throw `${LOG}  Status error: ${res.status}`;
                }
    
                return res.text();
            };

            // METRICS    
            const getMetricsSafe = () => {
                try {
                    const metrics = this.gatherMetrics();
                    console.log(`${LOG} metrics:`, metrics);
                    return metrics;
                } catch (e) {
                    console.error(`${LOG}  Metrics error:`, e);
                }
                return {};
            }

            const payload = {
                rating: rating,
                rating_comment: comment,
                metadata: {
                    ...getMetricsSafe()
                }
            }

            postFeedback(jwt, payload)
                .then(result => {
                    console.log(`${LOG} feedback result: `, result);
                })
                .catch(e => console.error(`${LOG} failed to feedback`, e));
        }

        enableFeedbackOnLeave() {
            const conf = window.APP.store.getState()['features/base/conference'].conference;
            conf.isCallstatsEnabled = () => true;
        }

        _getMatrixContext() {
            const token = window.APP.store.getState()['features/base/jwt'].jwt;
            const payload = token.split('.')[1];
            const content = JSON.parse(atob(payload));
            return content.context;
        }

        handleJoin() {
            this.enableFeedbackOnLeave();

            console.log(`${LOG} jitsimeet`, JitsiMeetJS);

            // Extract matrix openId token from the Jitsi JWT token
            const oidToken = this._getMatrixContext().matrix.token;

            console.log(`${LOG} Extracted matrix token: ${oidToken}`);

            const getToken = async (oidToken) => {
                const baseUrl = APP.store.getState()['features/base/config'].feedbackBackend;
                const url = `${baseUrl}/token`;
                
                const headers = {
                    'authorization': `Bearer ${oidToken}`
                };
    
                const res = await fetch(url, {
                    method: 'GET',
                    headers
                });
    
                if (!res.ok) {
                    throw `${LOG}  Status error: ${res.status}`;
                }
    
                return res.text();
            };

            // get the feedback JWT token as soon as possible
            getToken(oidToken)
                .then(feedbackToken => {
                    console.log(`${LOG} feedback JWT: ${feedbackToken}`);
                    window.APP.conference.feedbackToken = feedbackToken;
                })
                .catch(e => console.error(`${LOG} failed to fetch JWT`, e));
                
            return;
        }

        gatherMetrics() {
            const metrics = {};

            const config = APP.store.getState()['features/base/config'];
            const conference = window.APP.store.getState()['features/base/conference'].conference;
            const localParticipant = window.APP.store.getState()['features/base/participants'].local;

            // meetingId
            metrics.meetingUrl = window.location.href;
            metrics.meetingId  = JitsiMeetJS.analytics.permanentProperties.conference_name;

            // participantID
            metrics.participantId = conference.myUserId();

            // matrix user Id
            metrics.matrixUserId = localParticipant.email;
            metrics.displayName  = localParticipant.name;

            // 
            metrics.userRegion = JitsiMeetJS.analytics.permanentProperties.userRegion;

            // app data
            // const {appName,environment,releaseNumber,envType,backendRelease} = JitsiMeetJS.analytics.permanentProperties;
            metrics.appLibVersion = JitsiMeetJS.version;
            metrics.appFocusVersion = conference.componentsVersions.versions.focus;
            metrics.appName = JitsiMeetJS.analytics.permanentProperties.appName;
            metrics.appMeetingRegion = config.deploymentInfo.region;
            metrics.appShard = config.deploymentInfo.shard;
            metrics.appRegion = config.deploymentInfo.region;
            metrics.appEnvironment = config.deploymentInfo.environment;
            metrics.appEnvType = config.deploymentInfo.envType;

            // browser, os
            metrics.userAgent = JitsiMeetJS.analytics.permanentProperties.user_agent;
            metrics.browserName = JitsiMeetJS.util.browser.getName();
            metrics.browserVersion = JitsiMeetJS.util.browser.getVersion();
            metrics.osName = JitsiMeetJS.util.browser._bowser.parseOS().name;
            metrics.osVersion = JitsiMeetJS.util.browser._bowser.parseOS().version;
            metrics.osVersionName = JitsiMeetJS.util.browser._bowser.parseOS().versionName;

            // standalone or embedded
            metrics.externalApi = JitsiMeetJS.analytics.permanentProperties.externalApi;
            metrics.inIframe    = JitsiMeetJS.analytics.permanentProperties.inIframe;

            return metrics;
        }
    }

    const getJitsiMeetGlobalNS = () => {
        return window.JitsiMeetJS.app;
    }

    getJitsiMeetGlobalNS().analyticsHandlers.push(Feedback);

})();