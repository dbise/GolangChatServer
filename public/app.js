new Vue({
    el: '#app',

    data: {
        ws: null, // Connecting WebSocket
        newMsg: '', // New message to be sent to the server
        chatContent: '', // The list of chat content on the server
        email: null, // Email
        username: null, // Username
        joined: false // Return true if username and password filled out
    },

    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            self.chatContent += '<div class="chip">'
                + msg.username
                + ': '
                + '</div>'
                + msg.message + '<br/>';
            var element = document.getElementById('chatMessages');
            element.scrollTop = element.scrollHeight; // Automatically scroll to bottom
        });
    },

    methods: {
        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        email: this.email,
                        username: this.username,
                        message: $('<p>').html(this.newMsg).text() // get message to send
                    }
                ));
                this.newMsg = ''; // Reset next message to blank
            }
        },

        join: function () {
            if (!this.email) {
                Materialize.toast('An email must be entered', 2000);
                return
            }
            if (!this.username) {
                Materialize.toast('A Username must be entered', 2000);
                return
            }
            this.email = $('<p>').html(this.email).text();
            this.username = $('<p>').html(this.username).text();
            this.joined = true;
        }
    }
});
