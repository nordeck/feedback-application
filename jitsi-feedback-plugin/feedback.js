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

            const metrics = this._gatherMetrics();

            const payload = {
                rating: rating,
                rating_comment: comment,
                metadata: {
                    ...metrics
                }
            }

            postFeedback(jwt, payload)
                .then(result => {
                    console.log(`${LOG} feedback result: `, payload);
                })
                .catch(e => console.error(`${LOG} failed to feedback`, e));
        }

        enableFeedbackOnLeave() {
            const conf = window.APP.store.getState()['features/base/conference'].conference;
            conf.isCallstatsEnabled = () => true;
        }

        handleJoin() {
            this.enableFeedbackOnLeave();

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

        _gatherMetrics() {
            const metrics = {};
            const config = APP.store.getState()['features/base/config'];
            const flags = config.metadata || [];

            flags.forEach(flag => {
                try {
                    this._addMetric(flag, metrics)
                } catch (e) {
                    console.error(`${LOG} Metrics error ${flag}:`, e);
                }
            });
            
            return metrics;
        }

        _getMatrixContext() {
            const token = window.APP.store.getState()['features/base/jwt'].jwt;
            const payload = token.split('.')[1];
            const content = JSON.parse(atob(payload));
            return content.context;
        }

        _addMetric(flag, metrics)  {
            const config = APP.store.getState()['features/base/config'];
            const conference = window.APP.store.getState()['features/base/conference'].conference;
            const localParticipant = window.APP.store.getState()['features/base/participants'].local;

            switch (flag) {
                // meetingId
                case 'MEETING_URL':
                    metrics.meetingUrl = window.location.href;
                    break;
                case 'MEETING_ID':
                    metrics.meetingId  = JitsiMeetJS.analytics.permanentProperties.conference_name;
                    break;

                // participantID    
                case 'PARTICIPANT_ID':
                    metrics.participantId = conference.myUserId();
                    break;

                // matrix user Id    
                case 'MATRIX_USER_ID':
                    metrics.matrixUserId = localParticipant.email;
                    break;
                case 'DISPLAY_NAME':
                    metrics.displayName  = localParticipant.name;
                    break;

                case 'USER_REGION':
                    metrics.userRegion = JitsiMeetJS.analytics.permanentProperties.userRegion;
                    break;
                
                // app data
                case 'APP_LIB_VERSION':
                    metrics.appLibVersion = JitsiMeetJS.version;
                    break;   
                case 'APP_FOCUS_VERSION':
                    metrics.appFocusVersion = conference.componentsVersions.versions.focus;
                    break; 
                case 'APP_NAME':
                    metrics.appName = JitsiMeetJS.analytics.permanentProperties.appName;
                    break;
                case 'APP_MEETING_REGION':
                    metrics.appMeetingRegion = config.deploymentInfo.region;
                    break;
                case 'APP_SHARD':
                    metrics.appShard = config.deploymentInfo.shard;
                    break;
                case 'APP_REGION':
                    metrics.appRegion = config.deploymentInfo.region;
                    break;
                case 'APP_ENVIRONMENT':
                    metrics.appEnvironment = config.deploymentInfo.environment;
                    break; 
                case 'APP_ENV_TYPE':
                    metrics.appEnvType = config.deploymentInfo.envType;
                    break;
                case 'APP_BACKEND_RELEASE':
                    metrics.appBackendRelease = config.deploymentInfo.backendRelease;
                    break;

                // browser, os
                case 'USER_AGENT':
                    metrics.userAgent = JitsiMeetJS.analytics.permanentProperties.user_agent;
                    break;
                case 'BROWSER_NAME':
                    metrics.browserName = JitsiMeetJS.util.browser.getName();
                    break;
                case 'BROWSER_VERSION':
                    metrics.browserVersion = JitsiMeetJS.util.browser.getVersion();
                    break;
                case 'OS_NAME':
                    metrics.osName = JitsiMeetJS.util.browser._bowser.parseOS().name;
                    break;
                case 'OS_VERSION':
                    metrics.osVersion = JitsiMeetJS.util.browser._bowser.parseOS().version;
                    break;
                case 'OS_VERSION_NAME':
                    metrics.osVersionName = JitsiMeetJS.util.browser._bowser.parseOS().versionName;
                    break;
                
                // standalone or embedded
                case 'EXTERNAL_API':
                    metrics.externalApi = JitsiMeetJS.analytics.permanentProperties.externalApi;
                    break;
                case 'IN_IFRAME':
                    metrics.inIframe = JitsiMeetJS.analytics.permanentProperties.inIframe;
                    break;
            }
        } 
    }

    const getJitsiMeetGlobalNS = () => {
        return window.JitsiMeetJS.app;
    }

    getJitsiMeetGlobalNS().analyticsHandlers.push(Feedback);

})();