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

            const payload = {
                rating: rating,
                rating_comment: comment,
                metadata: {
                    field1: "field1"
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

        handleJoin() {
            this.enableFeedbackOnLeave();

            // Extract matrix openId token from the Jitsi JWT token
            const token = window.APP.store.getState()['features/base/jwt'].jwt;
            const payload = token.split('.')[1];
            const content = JSON.parse(atob(payload));
            const oidToken = content.context.matrix.token;

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
    }

    const getJitsiMeetGlobalNS = () => {
        return window.JitsiMeetJS.app;
    }

    getJitsiMeetGlobalNS().analyticsHandlers.push(Feedback);

})();